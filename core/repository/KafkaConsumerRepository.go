package repository

import (
	"context"
	"messages/core/entity"
)

type KafkaConsumerRepository interface {
	// Subscribe starts consuming from the given topic.
	// Events are sent to the returned channel until ctx is cancelled.
	Subscribe(ctx context.Context, topic string) (<-chan entity.ConsumerEvent, error)
}
