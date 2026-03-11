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
	Headers  []Header    `json:"headers"`
	Content  interface{} `json:"content" validate:"required"`
}

type Header struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

func (dto MessagesRequestDto) ToCoreEntity() (entity.HandlerMessages, error) {
	kafkaTopicsMessages := make([]entity.Message, 0)
	activeMQQueuesMessages := make([]entity.Message, 0)
	activeMQTopicsMessages := make([]entity.Message, 0)

	for _, message := range dto.Messages {
		headers := make([]entity.Header, len(message.Headers))
		for i := range message.Headers {
			headers[i] = entity.Header{Key: message.Headers[i].Key, Value: message.Headers[i].Value}
		}
		contentMarshal, _ := json.Marshal(message.Content)
		contentString := string(contentMarshal)
		switch message.Platform {
		case "kafka":
			switch message.Type {
			case "topic":
				kafkaTopicsMessages = append(kafkaTopicsMessages, entity.Message{Name: message.Name, Content: contentString, Headers: headers})
			case "queue":
				return entity.HandlerMessages{}, errors.New(fmt.Sprintf("Kafka does not support queues: %s", message.Name))
			}
		case "activemq":
			switch message.Type {
			case "queue":
				activeMQQueuesMessages = append(activeMQQueuesMessages, entity.Message{Name: message.Name, Content: contentString, Headers: headers})
			case "topic":
				activeMQTopicsMessages = append(activeMQTopicsMessages, entity.Message{Name: message.Name, Content: contentString, Headers: headers})
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
