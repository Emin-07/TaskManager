package producer

import (
	"context"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/Emin-07/TaskManager/internal/adapter/kafka/shared"
)

type KafkaProducer struct {
	cfg    *shared.KafkaConfig
	writer *kafka.Writer
}

func NewKafkaProducer() *KafkaProducer {
	cfg := shared.NewKafkaConfig()
	w := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Addr...),
		AllowAutoTopicCreation: true,
	}
	return &KafkaProducer{
		cfg:    cfg,
		writer: w,
	}
}

func (kp *KafkaProducer) Publish(data map[string]string, topic string) error {
	// Make a writer that publishes messages to topic-A.
	// The topic will be created if it is missing.

	messages := make([]kafka.Message, len(data))
	i := 0
	for key, val := range data {
		messages[i] = kafka.Message{Key: []byte(key), Value: []byte(val), Topic: topic}
		i++
	}

	var err error
	const retries = 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// attempt to create topic prior to publishing the message
		err = kp.writer.WriteMessages(ctx, messages...)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			return err
		}
		break
	}

	if err := kp.writer.Close(); err != nil {
		return err
	}
	return nil
}
