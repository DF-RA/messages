package repository

import "messages/core/entity"

type ActiveMQRepository interface {
	SendTopic(message entity.Message) entity.MessageOutput
	SendQueue(message entity.Message) entity.MessageOutput
}
