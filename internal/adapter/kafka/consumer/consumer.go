package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/segmentio/kafka-go"

	"github.com/Emin-07/TaskManager/internal/adapter/kafka/shared"
	"github.com/Emin-07/TaskManager/internal/core/domain"
	"github.com/Emin-07/TaskManager/internal/core/port"
)

const (
	workerCount    = 4
	jobChannelSize = 64
)

type KafkaConsumer struct {
	cfg          *shared.KafkaConfig
	reader       *kafka.Reader
	userServices port.UserService
	taskServices port.TaskService
}

func NewKafkaConsumer(topic string) *KafkaConsumer {
	cfg := shared.NewKafkaConfig()
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   cfg.Addr,
		Topic:     topic,
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	return &KafkaConsumer{
		cfg:    cfg,
		reader: r,
	}
}

func (kc *KafkaConsumer) Consume(ctx context.Context) error {
	// job channel carries messages from the fetch loop to worker goroutines
	jobs := make(chan kafka.Message, jobChannelSize)
	errCh := make(chan error, workerCount)

	var wg sync.WaitGroup

	// start the fixed worker pool
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range jobs {
				kc.processMessage(ctx, msg, errCh)
			}
		}()
	}

	// fetch loop (single producer)
	fetchErr := kc.fetchLoop(ctx, jobs)

	// close jobs so workers drain remaining messages and exit
	close(jobs)
	wg.Wait()

	// if fetch loop errored due to context, return context error
	if ctx.Err() != nil {
		_ = kc.reader.Close()
		return ctx.Err()
	}

	if err := kc.reader.Close(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return fetchErr
}

func (kc *KafkaConsumer) fetchLoop(ctx context.Context, jobs chan<- kafka.Message) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m, err := kc.reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return err
				}
				log.Printf("error fetching message: %v", err)
				continue
			}
			select {
			case jobs <- m:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func (kc *KafkaConsumer) processMessage(ctx context.Context, m kafka.Message, errCh chan<- error) {
	var err error
	switch kc.reader.Config().Topic {
	case domain.TopicUsers:
		err = kc.handleUserMessage(ctx, m.Value)
	case domain.TopicTasks:
		err = kc.handleTaskMessage(ctx, m.Value)
	}
	if err != nil {
		errCh <- fmt.Errorf("couldn't process message at topic/partition/offset %v/%v/%v: %w", m.Topic, m.Partition, m.Offset, err)
		return
	}

	log.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	if err = kc.reader.CommitMessages(ctx, m); err != nil {
		errCh <- fmt.Errorf("couldn't commit message: %w", err)
	}
}

func (kc *KafkaConsumer) handleUserMessage(ctx context.Context, data []byte) error {
	operation := shared.MsgOperation{}
	if err := json.Unmarshal(data, &operation); err != nil {
		return fmt.Errorf("couldn't extract operation from struct: %s", err)
	}

	switch strings.ToLower(operation.Operation) {
	case domain.CreateOperation:
		postData := shared.UserBasic{}
		if err := json.Unmarshal(data, &postData); err != nil {
			return fmt.Errorf("couldn't extract data from struct: %s", err)
		}
		newUser := &domain.SignupUser{
			Username: postData.Username,
			Role:     postData.Role,
			Email:    postData.Email,
			Password: postData.Password,
		}
		return kc.userServices.Insert(ctx, newUser)
	case domain.PatchOperation:
		patchData := shared.UserPatch{}
		if err := json.Unmarshal(data, &patchData); err != nil {
			return fmt.Errorf("couldn't extract data from struct: %s", err)
		}
		patchedUser := &domain.SignupUser{
			Username: patchData.Username,
			Role:     patchData.Role,
			Email:    patchData.Email,
			Password: patchData.Password,
		}
		return kc.userServices.Patch(ctx, patchedUser, patchData.ID)
	//case shared.ChangeOperation:
	//	kc.userServices.Patch()
	case domain.DeleteOperation:
		deleteData := shared.UserDelete{}
		if err := json.Unmarshal(data, &deleteData); err != nil {
			return fmt.Errorf("couldn't extract data from struct: %s", err)
		}
		return kc.userServices.Delete(ctx, deleteData.ID)
	default:
		return fmt.Errorf("unknown operation passed in : %s", operation.Operation)
	}
}

func (kc *KafkaConsumer) handleTaskMessage(ctx context.Context, data []byte) error {
	operation := shared.MsgOperation{}
	if err := json.Unmarshal(data, &operation); err != nil {
		return fmt.Errorf("couldn't extract operation from struct: %s", err)
	}
	switch strings.ToLower(operation.Operation) {
	case domain.CreateOperation:
		postData := shared.TaskCreate{}
		if err := json.Unmarshal(data, &postData); err != nil {
			return fmt.Errorf("couldn't extract data from struct: %s", err)
		}
		newTask := &domain.CreateTask{
			Title:      postData.Title,
			Text:       postData.Text,
			Priority:   postData.Priority,
			ExpireDays: postData.ExpireDays,
		}
		return kc.taskServices.Post(ctx, newTask, postData.UserId)
	case domain.PatchOperation:
		patchData := shared.TaskPatch{}
		if err := json.Unmarshal(data, &patchData); err != nil {
			return fmt.Errorf("couldn't extract data from struct: %s", err)
		}
		patchedTask := &domain.CreateTask{
			Title:      patchData.Title,
			Text:       patchData.Text,
			Priority:   patchData.Priority,
			ExpireDays: patchData.ExpireDays,
		}
		return kc.taskServices.Patch(ctx, patchedTask, patchData.ID, patchData.UserId, patchData.Role)
	//case shared.ChangeOperation:
	//	kc.taskServices.Patch()
	case domain.DeleteOperation:
		deleteData := shared.TaskDelete{}
		if err := json.Unmarshal(data, &deleteData); err != nil {
			return fmt.Errorf("couldn't extract data from struct: %s", err)
		}
		return kc.taskServices.Delete(ctx, deleteData.ID, deleteData.UserId, deleteData.Role)
	default:
		return fmt.Errorf("unknown operation passed in : %s", operation.Operation)
	}

}
