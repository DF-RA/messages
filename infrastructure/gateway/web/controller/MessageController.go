package controller

import (
	"github.com/go-chi/chi/v5"
	useCases "messages/core/use-cases"
	"messages/infrastructure/config"
	"messages/infrastructure/gateway/web/dto"
	"net/http"
)

type MessageControllerImpl struct {
	SendTopicsAndQueues useCases.SendTopicsAndQueues
}
type MessageController interface{}

func NewMessageController(
	SendTopicsAndQueues useCases.SendTopicsAndQueues,
) chi.Router {
	router := chi.NewRouter()
	endpoints(router, &MessageControllerImpl{
		SendTopicsAndQueues,
	})
	return router
}

func endpoints(router chi.Router, impl *MessageControllerImpl) {
	router.Post("/handlers/send", impl.SendMessage)
}

func (impl *MessageControllerImpl) SendMessage(writer http.ResponseWriter, request *http.Request) {
	messages, err := config.GetBody[dto.MessagesRequestDto](request.Body)
	if err != nil {
		config.ErrorResponse(writer, http.StatusBadRequest, err)
		return
	}
	coreEntity, err := messages.ToCoreEntity()
	if err != nil {
		config.ErrorResponse(writer, http.StatusBadRequest, err)
		return
	}
	result, err := impl.SendTopicsAndQueues.Process(coreEntity)
	if err != nil {
		config.ErrorResponse(writer, http.StatusInternalServerError, err)
		return
	}
	config.ResponseWithoutBody[dto.MessagesResponseDto](writer, dto.ToMessagesResponseDto(result))
}
