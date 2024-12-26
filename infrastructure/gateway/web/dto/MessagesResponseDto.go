package dto

import (
	"messages/core/entity"
)

type MessagesResponseDto struct {
	Messages []MessageResponseDto `json:"messages"`
}

type MessageResponseDto struct {
	Name     string `json:"name"`
	Platform string `json:"platform"`
	Type     string `json:"type"`
	Status   string `json:"status"`
}

func ToMessagesResponseDto(coreEntity *entity.ResultMessages) MessagesResponseDto {
	result := MessagesResponseDto{
		Messages: []MessageResponseDto{},
	}
	if coreEntity == nil {
		return result
	}

	for _, topic := range coreEntity.Kafka.Topics {
		result.Messages = append(result.Messages, MessageResponseDto{
			Name:     topic.Name,
			Platform: "Kafka",
			Type:     "Topic",
			Status:   topic.Status,
		})
	}

	for _, queue := range coreEntity.ActiveMQ.Queues {
		result.Messages = append(result.Messages, MessageResponseDto{
			Name:     queue.Name,
			Platform: "ActiveMQ",
			Type:     "Queue",
			Status:   queue.Status,
		})
	}

	for _, topic := range coreEntity.ActiveMQ.Topics {
		result.Messages = append(result.Messages, MessageResponseDto{
			Name:     topic.Name,
			Platform: "ActiveMQ",
			Type:     "Topic",
			Status:   topic.Status,
		})
	}

	return result
}
