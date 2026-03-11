package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	useCases "messages/core/use-cases"
	"messages/core/services"
	"messages/infrastructure/config"
)

type ConsumerControllerImpl struct {
	ConsumeMessages useCases.ConsumeMessages
	EventStore      services.EventStore
}

type startRequestDto struct {
	Source   string `json:"source" validate:"required,oneof=kafka activemq"`
	Name     string `json:"name" validate:"required"`
	DestType string `json:"type" validate:"required,oneof=topic queue"`
}

func NewConsumerController(
	consumeMessages useCases.ConsumeMessages,
	eventStore services.EventStore,
) chi.Router {
	router := chi.NewRouter()
	consumerEndpoints(router, &ConsumerControllerImpl{
		ConsumeMessages: consumeMessages,
		EventStore:      eventStore,
	})
	return router
}

func consumerEndpoints(router chi.Router, impl *ConsumerControllerImpl) {
	router.Post("/start", impl.Start)
	router.Post("/stop", impl.Stop)
	router.Get("/events", impl.StreamEvents)
	router.Delete("/events", impl.ClearEvents)
}

func (impl *ConsumerControllerImpl) Start(writer http.ResponseWriter, request *http.Request) {
	body, err := config.GetBody[startRequestDto](request.Body)
	if err != nil {
		config.ErrorResponse(writer, http.StatusBadRequest, err)
		return
	}
	if err := impl.ConsumeMessages.Start(body.Source, body.Name, body.DestType); err != nil {
		config.ErrorResponse(writer, http.StatusInternalServerError, err)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (impl *ConsumerControllerImpl) Stop(writer http.ResponseWriter, request *http.Request) {
	if err := impl.ConsumeMessages.Stop(); err != nil {
		config.ErrorResponse(writer, http.StatusInternalServerError, err)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (impl *ConsumerControllerImpl) StreamEvents(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := writer.(http.Flusher)
	if !ok {
		config.ErrorResponse(writer, http.StatusInternalServerError, fmt.Errorf("streaming not supported"))
		return
	}

	// Send history first
	for _, event := range impl.EventStore.History() {
		data, _ := json.Marshal(event)
		fmt.Fprintf(writer, "data: %s\n\n", data)
	}
	flusher.Flush()

	// Subscribe to live events
	ch := impl.EventStore.Subscribe()
	defer impl.EventStore.Unsubscribe(ch)

	for {
		select {
		case event, ok := <-ch:
			if !ok {
				return
			}
			data, _ := json.Marshal(event)
			fmt.Fprintf(writer, "data: %s\n\n", data)
			flusher.Flush()
		case <-request.Context().Done():
			return
		}
	}
}

func (impl *ConsumerControllerImpl) ClearEvents(writer http.ResponseWriter, request *http.Request) {
	impl.EventStore.Clear()
	writer.WriteHeader(http.StatusNoContent)
}
