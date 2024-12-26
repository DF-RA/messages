package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log"
	useCases "messages/core/use-cases"
	"messages/infrastructure/gateway/event/impl"
	"messages/infrastructure/gateway/web/controller"
	"net/http"
	"os"
)

func main() {
	// Envs
	if err := godotenv.Load("./devEnv/.env"); err != nil {
		log.Printf("Failed to load envs: %v", err)
	}

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
	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT env is required")
	}
	println("Server started on port " + port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return
	}
}
