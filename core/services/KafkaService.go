package services

import (
	"messages/core/entity"
	"messages/core/repository"
)

type KafkaServiceImpl struct {
	KafkaRepository repository.KafkaRepository
	MessageService  MessageService
}
type KafkaService interface {
	SendMessages(kafka entity.Kafka) (entity.KafkaOutput, error)
}

func NewKafkaService(
	KafkaRepository repository.KafkaRepository,
) *KafkaServiceImpl {
	return &KafkaServiceImpl{
		MessageService:  &MessageServiceImpl{},
		KafkaRepository: KafkaRepository,
	}
}

func (impl KafkaServiceImpl) SendMessages(kafka entity.Kafka) (entity.KafkaOutput, error) {
	return entity.KafkaOutput{
		Topics: impl.MessageService.ProcessMessage(kafka.Topics, impl.KafkaRepository.SendTopic),
	}, nil
}
