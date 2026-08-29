package port

import (
	"context"
)

type Publisher interface {
	Publish(data map[string]string, topic string) error
}
type Consumer interface {
	Consume(ctx context.Context) error
}

type MessageBrokerIn interface {
	Publisher
	Consumer
}

type MessageBrokerOut interface {
	Publisher
	Consumer
}
