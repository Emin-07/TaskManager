package port

import (
	"context"
)

type MessageBrokerIn interface {
	Publish(data map[string]string) error
	Consume(ctx context.Context) error
}

type MessageBrokerOut interface {
	Publish(data map[string]string) error
	Consume(ctx context.Context) error
}
