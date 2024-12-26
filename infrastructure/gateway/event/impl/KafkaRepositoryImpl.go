package impl

import (
	"context"
	"github.com/segmentio/kafka-go"
	"messages/core/entity"
	"messages/infrastructure/config"
	"time"
)

type KafkaRepositoryImpl struct {
	KafkaConfig config.KafkaConfig
}

func NewKafkaRepositoryImpl() *KafkaRepositoryImpl {
	return &KafkaRepositoryImpl{
		KafkaConfig: config.NewKafkaConfig(),
	}
}

func (impl *KafkaRepositoryImpl) SendTopic(message entity.Message) entity.MessageOutput {
	conn, err := kafka.DialLeader(context.Background(), "tcp", impl.KafkaConfig.GetBroker(), message.Name, 0)
	if err != nil {
		return entity.MessageOutput{Name: message.Name, Status: err.Error()}
	}
	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return entity.MessageOutput{Name: message.Name, Status: err.Error()}
	}
	if _, err := conn.WriteMessages(kafka.Message{Value: []byte(message.Content)}); err != nil {
		return entity.MessageOutput{Name: message.Name, Status: err.Error()}
	}
	if err := conn.Close(); err != nil {
		return entity.MessageOutput{Name: message.Name, Status: err.Error()}
	}
	return entity.MessageOutput{Name: message.Name, Status: "Success"}
}
