package interfaces

import "context"

type MessengerProducer interface {
	Publish(ctx context.Context, msg, topic string) error
	Close()
}
