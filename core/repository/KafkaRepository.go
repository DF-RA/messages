package repository

import "messages/core/entity"

type KafkaRepository interface {
	SendTopic(message entity.Message) entity.MessageOutput
}
