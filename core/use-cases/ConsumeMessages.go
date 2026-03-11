package use_cases

import (
	"fmt"
	"log"
	"messages/core/entity"
	"messages/core/repository"
	"messages/core/services"
)

type ConsumeMessages interface {
	Start(source string, name string, destType string) error
	Stop() error
}

type ConsumeMessagesImpl struct {
	KafkaConsumer    repository.KafkaConsumerRepository
	ActiveMQConsumer repository.ActiveMQConsumerRepository
	EventStore       services.EventStore
	internalCh       chan entity.ConsumerEvent
}

func NewConsumeMessages(
	kafkaConsumer repository.KafkaConsumerRepository,
	activeMQConsumer repository.ActiveMQConsumerRepository,
	eventStore services.EventStore,
) *ConsumeMessagesImpl {
	return &ConsumeMessagesImpl{
		KafkaConsumer:    kafkaConsumer,
		ActiveMQConsumer: activeMQConsumer,
		EventStore:       eventStore,
	}
}

// Start begins consuming from the given source.
// source: "kafka" | "activemq"
// destType: "topic" | "queue" (only relevant for activemq)
func (uc *ConsumeMessagesImpl) Start(source string, name string, destType string) error {
	uc.internalCh = make(chan entity.ConsumerEvent, 128)

	// Relay: consumer writes to internalCh → Publish broadcasts to all SSE clients
	go func() {
		for event := range uc.internalCh {
			log.Printf("[relay] publishing event source=%s name=%s", event.Source, event.Name)
			uc.EventStore.Publish(event)
		}
	}()

	switch source {
	case "kafka":
		return uc.KafkaConsumer.Subscribe(name, uc.internalCh)
	case "activemq":
		return uc.ActiveMQConsumer.Subscribe(name, destType, uc.internalCh)
	default:
		close(uc.internalCh)
		return fmt.Errorf("unknown source: %s", source)
	}
}

func (uc *ConsumeMessagesImpl) Stop() error {
	kafkaErr := uc.KafkaConsumer.Unsubscribe()
	activeMQErr := uc.ActiveMQConsumer.Unsubscribe()
	// The consumer goroutine closes internalCh via defer when it exits,
	// which causes the relay goroutine to exit cleanly — no explicit close here.
	if kafkaErr != nil {
		return kafkaErr
	}
	return activeMQErr
}
