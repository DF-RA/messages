package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	useCases "messages/core/use-cases"
	"messages/infrastructure/gateway/event/impl"
	"messages/infrastructure/gateway/web/controller"
)

func main() {
	// Envs
	if err := godotenv.Load("./devEnv/.env"); err != nil {
		log.Printf("Failed to load envs: %v", err)
	}

	// Routes
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Producer
	kafkaRepository := impl.NewKafkaRepositoryImpl()
	activeMQRepository := impl.NewActiveMqRepositoryImpl()
	sendTopicsAndQueues := useCases.NewSendTopicsAndQueues(kafkaRepository, activeMQRepository)
	router.Mount("/messages", controller.NewMessageController(sendTopicsAndQueues))

	// Consumer (per-client, lifecycle tied to SSE connection)
	kafkaConsumer := impl.NewKafkaConsumerRepositoryImpl()
	activeMQConsumer := impl.NewActiveMQConsumerRepositoryImpl()
	router.Mount("/consumer", controller.NewConsumerController(kafkaConsumer, activeMQConsumer))

	// Static files (embedded)
	subFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("Failed to load static files: %v", err)
	}
	router.Handle("/*", http.FileServer(http.FS(subFS)))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT env is required")
	}
	println("Server started on port " + port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
