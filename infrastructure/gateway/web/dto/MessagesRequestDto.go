package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"messages/core/entity"
)

type MessagesRequestDto struct {
	Messages []MessageRequestDto `json:"messages" validate:"min=1"`
}

type MessageRequestDto struct {
	Name     string      `json:"name" validate:"required"`
	Platform string      `json:"platform" validate:"required,oneof=kafka activemq"`
	Type     string      `json:"type" validate:"required,oneof=topic queue"`
	Content  interface{} `json:"content" validate:"required"`
}

func (dto MessagesRequestDto) ToCoreEntity() (entity.HandlerMessages, error) {
	kafkaTopicsMessages := make([]entity.Message, 0)
	activeMQQueuesMessages := make([]entity.Message, 0)
	activeMQTopicsMessages := make([]entity.Message, 0)

	for _, message := range dto.Messages {
		contentMarshal, _ := json.Marshal(message.Content)
		contentString := string(contentMarshal)
		switch message.Platform {
		case "kafka":
			switch message.Type {
			case "topic":
				kafkaTopicsMessages = append(kafkaTopicsMessages, entity.Message{Name: message.Name, Content: contentString})
			case "queue":
				return entity.HandlerMessages{}, errors.New(fmt.Sprintf("Kafka does not support queues: %s", message.Name))
			}
		case "activemq":
			switch message.Type {
			case "queue":
				activeMQQueuesMessages = append(activeMQQueuesMessages, entity.Message{Name: message.Name, Content: contentString})
			case "topic":
				activeMQTopicsMessages = append(activeMQTopicsMessages, entity.Message{Name: message.Name, Content: contentString})
			}
		}
	}
	return entity.HandlerMessages{
		Kafka: entity.Kafka{Topics: kafkaTopicsMessages},
		ActiveMQ: entity.ActiveMQ{
			Queues: activeMQQueuesMessages,
			Topics: activeMQTopicsMessages,
		},
	}, nil
}
