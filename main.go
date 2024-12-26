package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	useCases "messages/core/use-cases"
	"messages/infrastructure/gateway/event/impl"
	"messages/infrastructure/gateway/web/controller"
	"net/http"
	"os"
)

func main() {
	// Routes
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Repository Impl
	kafkaRepository := impl.NewKafkaRepositoryImpl()
	activeMQRepository := impl.NewActiveMqRepositoryImpl()

	// Use cases
	sendTopicsAndQueues := useCases.NewSendTopicsAndQueues(
		kafkaRepository,
		activeMQRepository,
	)

	router.Mount("/messages", controller.NewMessageController(sendTopicsAndQueues))

	// Start server
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return
	}
}
