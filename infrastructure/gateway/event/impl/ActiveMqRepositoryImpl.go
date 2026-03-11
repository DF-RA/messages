package impl

import (
	"fmt"
	"messages/core/entity"
	"messages/infrastructure/config"

	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
)

type ActiveMqRepositoryImpl struct {
	ActiveMQConfig config.ActiveMQConfig
}

func NewActiveMqRepositoryImpl() *ActiveMqRepositoryImpl {
	return &ActiveMqRepositoryImpl{
		ActiveMQConfig: config.NewActiveMQConfig(),
	}
}

func (impl *ActiveMqRepositoryImpl) SendTopic(message entity.Message) entity.MessageOutput {
	if err := impl.sendMessage(message, "topic"); err != nil {
		return entity.MessageOutput{Name: message.Name, Status: err.Error()}
	}
	return entity.MessageOutput{Name: message.Name, Status: "Success"}
}

func (impl *ActiveMqRepositoryImpl) SendQueue(message entity.Message) entity.MessageOutput {
	if err := impl.sendMessage(message, "queue"); err != nil {
		return entity.MessageOutput{Name: message.Name, Status: err.Error()}
	}
	return entity.MessageOutput{Name: message.Name, Status: "Success"}
}

func (impl *ActiveMqRepositoryImpl) sendMessage(message entity.Message, typeMessage string) error {
	conn, err := stomp.Dial("tcp", impl.ActiveMQConfig.GetBroker(),
		stomp.ConnOpt.Login(impl.ActiveMQConfig.GetUser(), impl.ActiveMQConfig.GetPassword()),
	)
	if err != nil {
		return fmt.Errorf("error connecting to ActiveMQ: %w", err)
	}
	defer conn.Disconnect()

	destination := fmt.Sprintf("/%s/%s", typeMessage, message.Name)
	opts := buildSendOptions(message.Headers)

	if err := conn.Send(destination, "application/json", []byte(message.Content), opts...); err != nil {
		return fmt.Errorf("error sending message to %s: %w", destination, err)
	}
	return nil
}

func buildSendOptions(headers []entity.Header) []func(*frame.Frame) error {
	opts := make([]func(*frame.Frame) error, 0, len(headers))
	for _, h := range headers {
		opts = append(opts, stomp.SendOpt.Header(h.Key, h.Value))
	}
	return opts
}
