package use_cases

import (
	"messages/core/entity"
	"messages/core/repository"
	"messages/core/services"
)

type SendTopicsAndQueuesImpl struct {
	ActiveMQService services.ActiveMQService
	KafkaService    services.KafkaService
}
type SendTopicsAndQueues interface {
	Process(newMessages entity.HandlerMessages) (*entity.ResultMessages, error)
}

func NewSendTopicsAndQueues(
	KafkaRepository repository.KafkaRepository,
	ActiveMQRepository repository.ActiveMQRepository,
) *SendTopicsAndQueuesImpl {
	return &SendTopicsAndQueuesImpl{
		KafkaService:    services.NewKafkaService(KafkaRepository),
		ActiveMQService: services.NewActiveMQService(ActiveMQRepository),
	}
}

func (impl SendTopicsAndQueuesImpl) Process(newMessages entity.HandlerMessages) (*entity.ResultMessages, error) {
	kafkaResult, err := impl.KafkaService.SendMessages(newMessages.Kafka)
	if err != nil {
		return nil, err
	}
	activemqResult, err := impl.ActiveMQService.SendMessages(newMessages.ActiveMQ)
	if err != nil {
		return nil, err
	}
	return &entity.ResultMessages{
		Kafka:    kafkaResult,
		ActiveMQ: activemqResult,
	}, nil
}
