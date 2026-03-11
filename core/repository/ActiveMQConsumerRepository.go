package repository

import "messages/core/entity"

type ActiveMQConsumerRepository interface {
	Subscribe(name string, destType string, ch chan<- entity.ConsumerEvent) error
	Unsubscribe() error
}
