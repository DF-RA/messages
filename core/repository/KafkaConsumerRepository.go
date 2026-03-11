package repository

import "messages/core/entity"

type KafkaConsumerRepository interface {
	Subscribe(topic string, ch chan<- entity.ConsumerEvent) error
	Unsubscribe() error
}
