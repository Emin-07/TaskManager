package service

import (
	"context"

	"github.com/Emin-07/TaskManager/internal/core/port"
)

type MessageBrokerService struct {
	broker port.MessageBrokerIn
}

func NewMessageBrokerService(p port.Publisher, c port.Consumer) MessageBrokerService {
	type combined struct {
		port.Publisher
		port.Consumer
	}
	return MessageBrokerService{
		broker: combined{p, c},
	}
}

func (bs MessageBrokerService) Publish(data map[string]string, topic string) error {
	return bs.broker.Publish(data, topic)
}

// Consume runs the underlying broker consumer in a goroutine and fans its
// completion (success or error) back onto a single buffered channel, so the
// caller can select between context cancellation and consumer completion.
func (bs MessageBrokerService) Consume(ctx context.Context) error {
	errs := make(chan error, 1)
	go func() {
		if err := bs.broker.Consume(ctx); err != nil {
			errs <- err
		}
		close(errs)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err, ok := <-errs:
		if !ok {
			return nil
		}
		return err
	}
}
