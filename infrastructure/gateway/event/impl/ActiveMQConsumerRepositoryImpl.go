package impl

import (
	"fmt"
	"messages/core/entity"
	"messages/infrastructure/config"
	"time"

	"github.com/go-stomp/stomp/v3"
	"github.com/google/uuid"
)

type ActiveMQConsumerRepositoryImpl struct {
	ActiveMQConfig config.ActiveMQConfig
	conn           *stomp.Conn
	sub            *stomp.Subscription
}

func NewActiveMQConsumerRepositoryImpl() *ActiveMQConsumerRepositoryImpl {
	return &ActiveMQConsumerRepositoryImpl{
		ActiveMQConfig: config.NewActiveMQConfig(),
	}
}

func (impl *ActiveMQConsumerRepositoryImpl) Subscribe(name string, destType string, ch chan<- entity.ConsumerEvent) error {
	conn, err := stomp.Dial("tcp", impl.ActiveMQConfig.GetBroker(),
		stomp.ConnOpt.Login(impl.ActiveMQConfig.GetUser(), impl.ActiveMQConfig.GetPassword()),
	)
	if err != nil {
		return fmt.Errorf("error connecting to ActiveMQ: %w", err)
	}
	impl.conn = conn

	destination := fmt.Sprintf("/%s/%s", destType, name)
	sub, err := conn.Subscribe(destination, stomp.AckAuto)
	if err != nil {
		conn.Disconnect()
		return fmt.Errorf("error subscribing to %s: %w", destination, err)
	}
	impl.sub = sub

	go func() {
		defer close(ch)
		for msg := range sub.C {
			if msg.Err != nil {
				return
			}
			ch <- entity.ConsumerEvent{
				ID:         uuid.NewString(),
				Source:     "activemq",
				Name:       name,
				Payload:    string(msg.Body),
				ReceivedAt: time.Now(),
			}
		}
	}()

	return nil
}

func (impl *ActiveMQConsumerRepositoryImpl) Unsubscribe() error {
	if impl.sub != nil {
		if err := impl.sub.Unsubscribe(); err != nil {
			return err
		}
		impl.sub = nil
	}
	if impl.conn != nil {
		if err := impl.conn.Disconnect(); err != nil {
			return err
		}
		impl.conn = nil
	}
	return nil
}
