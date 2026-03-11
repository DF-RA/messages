package config

import (
	"fmt"
	"os"
)

type ActiveMQConfigImpl struct {
	Broker   string
	User     string
	Password string
}

type ActiveMQConfig interface {
	GetBroker() string
	GetUser() string
	GetPassword() string
}

func NewActiveMQConfig() *ActiveMQConfigImpl {
	host := os.Getenv("ACTIVEMQ_HOST")
	port := os.Getenv("ACTIVEMQ_STOMP_PORT")
	user := os.Getenv("ACTIVEMQ_USER")
	password := os.Getenv("ACTIVEMQ_PASSWORD")
	if host == "" || port == "" || user == "" || password == "" {
		panic("ACTIVEMQ_HOST, ACTIVEMQ_STOMP_PORT, ACTIVEMQ_USER, ACTIVEMQ_PASSWORD envs are required")
	}
	return &ActiveMQConfigImpl{
		Broker:   fmt.Sprintf("%s:%s", host, port),
		User:     user,
		Password: password,
	}
}

func (impl *ActiveMQConfigImpl) GetBroker() string {
	return impl.Broker
}

func (impl *ActiveMQConfigImpl) GetUser() string {
	return impl.User
}

func (impl *ActiveMQConfigImpl) GetPassword() string {
	return impl.Password
}
