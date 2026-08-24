package service

import (
	"context"

	"github.com/Emin-07/TaskManager/internal/core/port"
)

type MessageBrokerService struct {
	broker port.MessageBrokerIn
}

func (bs *MessageBrokerService) Publish(data map[string]string) error {
	return bs.broker.Publish(data)
}

func (bs *MessageBrokerService) Consume(ctx context.Context) error {
	errs := make(chan error)
	go func() {
		err := bs.broker.Consume(ctx)
		if err != nil {
			errs <- err
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errs:
			return err
		}
	}
}

