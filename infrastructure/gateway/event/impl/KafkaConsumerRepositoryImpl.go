package impl

import (
	"context"
	"log"
	"messages/core/entity"
	"messages/infrastructure/config"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumerRepositoryImpl struct {
	KafkaConfig config.KafkaConfig
}

func NewKafkaConsumerRepositoryImpl() *KafkaConsumerRepositoryImpl {
	return &KafkaConsumerRepositoryImpl{
		KafkaConfig: config.NewKafkaConfig(),
	}
}

func (impl *KafkaConsumerRepositoryImpl) Subscribe(ctx context.Context, topic string) (<-chan entity.ConsumerEvent, error) {
	// No GroupID: reads directly from the partition without consumer group coordination.
	// This avoids rebalancing delays and starts receiving messages immediately.
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{impl.KafkaConfig.GetBroker()},
		Topic:       topic,
		StartOffset: kafka.LastOffset,
	})

	ch := make(chan entity.ConsumerEvent, 128)

	log.Printf("[kafka-consumer] subscribed to topic=%s broker=%s", topic, impl.KafkaConfig.GetBroker())

	go func() {
		defer reader.Close()
		defer close(ch)
		for {
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					log.Printf("[kafka-consumer] context cancelled, stopping consumer for topic=%s", topic)
					return
				}
				log.Printf("[kafka-consumer] read error on topic=%s: %v", topic, err)
				continue
			}
			log.Printf("[kafka-consumer] received message on topic=%s payload=%s", topic, string(msg.Value))
			ch <- entity.ConsumerEvent{
				ID:         uuid.NewString(),
				Source:     "kafka",
				Name:       topic,
				Payload:    string(msg.Value),
				ReceivedAt: time.Now(),
			}
		}
	}()

	return ch, nil
}
