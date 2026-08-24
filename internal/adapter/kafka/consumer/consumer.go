package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/segmentio/kafka-go"

	"github.com/Emin-07/TaskManager/internal/adapter/kafka/shared"
	"github.com/Emin-07/TaskManager/internal/core/domain"
	"github.com/Emin-07/TaskManager/internal/core/port"
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

func (kc *KafkaConsumer) Consume(ctx context.Context, msgsCh chan<- domain.Message) error {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			m, err := kc.reader.FetchMessage(ctx)
			if err != nil {
				break
			}
			//msgsCh <- domain.Message{Key: m.Key, Val: m.Value, Topic: m.Topic}
			// TODO: after testing, based on topic convert it and send into db
			switch m.Topic {
			case "tasks":
				operation := shared.MsgOperation{}
				if err = json.Unmarshal(m.Value, &operation); err != nil {
					return fmt.Errorf("couldn't extract operation from struct: %s", err)
				}
			case "users":
				fmt.Print("Users")
			default:
				return fmt.Errorf("unknown topic in message : %s", m.Topic)
			}
			fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
			if err = kc.reader.CommitMessages(ctx, m); err != nil {
				return err
			}
		}

	}

	if err := kc.reader.Close(); err != nil {
		return err
	}
	return nil
}

func (kc *KafkaConsumer) handleMessage(ctx context.Context, data []byte, topic string) error {
	operation := shared.MsgOperation{}
	if err := json.Unmarshal(data, &operation); err != nil {
		return fmt.Errorf("couldn't extract operation from struct: %s", err)
	}
	switch topic {
	case "tasks":
		switch strings.ToLower(operation.Operation) {
		case shared.CreateOperation:
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
		case shared.PatchOperation:
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
		case shared.DeleteOperation:
			deleteData := shared.TaskDelete{}
			if err := json.Unmarshal(data, &deleteData); err != nil {
				return fmt.Errorf("couldn't extract data from struct: %s", err)
			}
			return kc.taskServices.Delete(ctx, deleteData.ID, deleteData.UserId, deleteData.Role)
		default:
			return fmt.Errorf("unknown operation passed in : %s", operation.Operation)
		}
	case "users":
		switch strings.ToLower(operation.Operation) {
		case shared.CreateOperation:
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
		case shared.PatchOperation:
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
		case shared.DeleteOperation:
			deleteData := shared.UserDelete{}
			if err := json.Unmarshal(data, &deleteData); err != nil {
				return fmt.Errorf("couldn't extract data from struct: %s", err)
			}
			return kc.userServices.Delete(ctx, deleteData.ID)
		default:
			return fmt.Errorf("unknown operation passed in : %s", operation.Operation)
		}
	}
	return nil
}
