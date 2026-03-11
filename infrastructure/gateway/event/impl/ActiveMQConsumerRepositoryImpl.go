package impl

import (
	"context"
	"fmt"
	"messages/core/entity"
	"messages/infrastructure/config"
	"time"

	"github.com/go-stomp/stomp/v3"
	"github.com/google/uuid"
)

type ActiveMQConsumerRepositoryImpl struct {
	ActiveMQConfig config.ActiveMQConfig
}

func NewActiveMQConsumerRepositoryImpl() *ActiveMQConsumerRepositoryImpl {
	return &ActiveMQConsumerRepositoryImpl{
		ActiveMQConfig: config.NewActiveMQConfig(),
	}
}

func (impl *ActiveMQConsumerRepositoryImpl) Subscribe(ctx context.Context, name string, destType string) (<-chan entity.ConsumerEvent, error) {
	conn, err := stomp.Dial("tcp", impl.ActiveMQConfig.GetBroker(),
		stomp.ConnOpt.Login(impl.ActiveMQConfig.GetUser(), impl.ActiveMQConfig.GetPassword()),
	)
	if err != nil {
		return nil, fmt.Errorf("error connecting to ActiveMQ: %w", err)
	}

	destination := fmt.Sprintf("/%s/%s", destType, name)
	sub, err := conn.Subscribe(destination, stomp.AckAuto)
	if err != nil {
		conn.Disconnect()
		return nil, fmt.Errorf("error subscribing to %s: %w", destination, err)
	}

	ch := make(chan entity.ConsumerEvent, 128)

	go func() {
		defer close(ch)
		defer conn.Disconnect()
		for {
			select {
			case msg, ok := <-sub.C:
				if !ok {
					return
				}
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
			case <-ctx.Done():
				sub.Unsubscribe()
				return
			}
		}
	}()

	return ch, nil
}
