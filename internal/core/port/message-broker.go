package port

import "context"

type MessageBrokerIn interface {
	Publish(data map[string]string) error
	Consume(ctx context.Context, msgsCh chan<- [][]byte) error
}

type MessageBrokerOut interface {
	Publish(data map[string]string) error
	Consume(ctx context.Context, msgsCh chan<- [][]byte) error
}
