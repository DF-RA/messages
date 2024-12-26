package impl

import (
	"bytes"
	"fmt"
	"messages/core/entity"
	"messages/infrastructure/config"
	"net/http"
)

type ActiveMqRepositoryImpl struct {
	ActiveMQConfig config.ActiveMQConfig
	Client         *http.Client
}

func NewActiveMqRepositoryImpl() *ActiveMqRepositoryImpl {
	client := &http.Client{}
	return &ActiveMqRepositoryImpl{
		ActiveMQConfig: config.NewActiveMQConfig(),
		Client:         client,
	}
}

func (impl *ActiveMqRepositoryImpl) SendTopic(message entity.Message) entity.MessageOutput {
	if err := impl.sendMessage(message.Name, "topic", message.Content); err != nil {
		return entity.MessageOutput{Name: message.Name, Status: err.Error()}
	}
	return entity.MessageOutput{Name: message.Name, Status: "Success"}
}

func (impl *ActiveMqRepositoryImpl) SendQueue(message entity.Message) entity.MessageOutput {
	if err := impl.sendMessage(message.Name, "queue", message.Content); err != nil {
		return entity.MessageOutput{Name: message.Name, Status: err.Error()}
	}
	return entity.MessageOutput{Name: message.Name, Status: "Success"}
}

func (impl *ActiveMqRepositoryImpl) sendMessage(destination string, typeMessage string, content string) error {
	request, err := buildRequest(impl.ActiveMQConfig.GetUrl(), destination, typeMessage, content)
	if err != nil {
		return err
	}
	response, err := impl.Client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Error sending message: %s", response.Status)
	}
	return nil
}

func buildRequest(url string, destination string, typeMessage string, content string) (*http.Request, error) {
	body := bytes.NewBuffer([]byte(content))
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	queryParam := request.URL.Query()
	queryParam.Add("destination", fmt.Sprintf("%s://%s", typeMessage, destination))
	request.URL.RawQuery = queryParam.Encode()

	return request, nil
}
