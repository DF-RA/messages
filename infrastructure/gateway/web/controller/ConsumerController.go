package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"messages/core/repository"
	"messages/infrastructure/config"
)

type ConsumerControllerImpl struct {
	KafkaConsumer    repository.KafkaConsumerRepository
	ActiveMQConsumer repository.ActiveMQConsumerRepository
}

func NewConsumerController(
	kafkaConsumer repository.KafkaConsumerRepository,
	activeMQConsumer repository.ActiveMQConsumerRepository,
) chi.Router {
	router := chi.NewRouter()
	impl := &ConsumerControllerImpl{
		KafkaConsumer:    kafkaConsumer,
		ActiveMQConsumer: activeMQConsumer,
	}
	router.Get("/events", impl.StreamEvents)
	return router
}

func (impl *ConsumerControllerImpl) StreamEvents(writer http.ResponseWriter, request *http.Request) {
	source := request.URL.Query().Get("source")
	name := request.URL.Query().Get("name")
	destType := request.URL.Query().Get("type")

	if source == "" || name == "" {
		config.ErrorResponse(writer, http.StatusBadRequest, fmt.Errorf("source and name are required"))
		return
	}
	if source != "kafka" && source != "activemq" {
		config.ErrorResponse(writer, http.StatusBadRequest, fmt.Errorf("source must be kafka or activemq"))
		return
	}
	if source == "activemq" && destType == "" {
		config.ErrorResponse(writer, http.StatusBadRequest, fmt.Errorf("type is required for activemq"))
		return
	}

	flusher, ok := writer.(http.Flusher)
	if !ok {
		config.ErrorResponse(writer, http.StatusInternalServerError, fmt.Errorf("streaming not supported"))
		return
	}

	ctx := request.Context()

	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	switch source {
	case "kafka":
		eventCh, err := impl.KafkaConsumer.Subscribe(ctx, name)
		if err != nil {
			fmt.Fprintf(writer, "event: error\ndata: %s\n\n", err.Error())
			flusher.Flush()
			return
		}
		for event := range eventCh {
			data, _ := json.Marshal(event)
			fmt.Fprintf(writer, "data: %s\n\n", data)
			flusher.Flush()
		}

	case "activemq":
		eventCh, err := impl.ActiveMQConsumer.Subscribe(ctx, name, destType)
		if err != nil {
			fmt.Fprintf(writer, "event: error\ndata: %s\n\n", err.Error())
			flusher.Flush()
			return
		}
		for event := range eventCh {
			data, _ := json.Marshal(event)
			fmt.Fprintf(writer, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
