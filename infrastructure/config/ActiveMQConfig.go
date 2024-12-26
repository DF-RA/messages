package config

import (
	"fmt"
	"os"
)

type ActiveMQConfigImpl struct {
	Url string
}
type ActiveMQConfig interface {
	GetUrl() string
}

func NewActiveMQConfig() *ActiveMQConfigImpl {
	host := os.Getenv("ACTIVEMQ_HOST")
	port := os.Getenv("ACTIVEMQ_PORT_HTTP")
	user := os.Getenv("ACTIVEMQ_USER")
	password := os.Getenv("ACTIVEMQ_PASSWORD")
	url := fmt.Sprintf("http://%s:%s@%s:%s/api/message", user, password, host, port)
	return &ActiveMQConfigImpl{Url: url}
}

func (impl *ActiveMQConfigImpl) GetUrl() string {
	return impl.Url
}
