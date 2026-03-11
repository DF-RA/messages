package repository

import (
	"context"
	"messages/core/entity"
)

type ActiveMQConsumerRepository interface {
	// Subscribe starts consuming from the given destination.
	// Events are sent to the returned channel until ctx is cancelled.
	Subscribe(ctx context.Context, name string, destType string) (<-chan entity.ConsumerEvent, error)
}
