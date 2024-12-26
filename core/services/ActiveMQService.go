package services

import (
	"messages/core/entity"
	"messages/core/repository"
)

type ActiveMQServiceImpl struct {
	ActiveMQRepository repository.ActiveMQRepository
	MessageService     MessageService
}
type ActiveMQService interface {
	SendMessages(activeMQ entity.ActiveMQ) (entity.ActiveMQOutput, error)
}

func NewActiveMQService(
	ActiveMQRepository repository.ActiveMQRepository,
) *ActiveMQServiceImpl {
	return &ActiveMQServiceImpl{
		MessageService:     &MessageServiceImpl{},
		ActiveMQRepository: ActiveMQRepository,
	}
}

func (impl ActiveMQServiceImpl) SendMessages(activeMQ entity.ActiveMQ) (entity.ActiveMQOutput, error) {
	return entity.ActiveMQOutput{
		Queues: impl.MessageService.ProcessMessage(activeMQ.Queues, impl.ActiveMQRepository.SendQueue),
		Topics: impl.MessageService.ProcessMessage(activeMQ.Topics, impl.ActiveMQRepository.SendTopic),
	}, nil
}
