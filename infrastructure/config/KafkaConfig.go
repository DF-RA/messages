package config

import (
	"fmt"
	"os"
)

type KafkaConfigImpl struct {
	Broker string
}

type KafkaConfig interface {
	GetBroker() string
}

func NewKafkaConfig() *KafkaConfigImpl {
	host := os.Getenv("KAFKA_HOST")
	port := os.Getenv("KAFKA_PORT")

	if host == "" || port == "" {
		panic("KAFKA_HOST, KAFKA_PORT envs are required")
	}

	return &KafkaConfigImpl{
		Broker: fmt.Sprintf("%s:%s", host, port),
	}
}

func (impl *KafkaConfigImpl) GetBroker() string {
	return impl.Broker
}
