package main

import (
	"log"
	"message-broker/internal/db"
	"message-broker/internal/queuing"
	"message-broker/internal/router"
	"message-broker/internal/service"
	"net/http"
)

func main() {

	mongoUri := "mongodb://localhost:27017"

	err := db.ConnectToMongoDB(mongoUri)

	if err != nil {
		log.Fatal(err)
		return
	}
	// Setup RabbitMQ connection and channel
	amqpURL := "amqp://guest:guest@localhost:5672/"
	if err := queuing.SetupAMQP(amqpURL); err != nil {
		log.Fatalf("Failed to setup RabbitMQ: %v\n", err)
	}
	go queuing.ConsumeMessages("SaveMessageQueue", service.CreateMessage)
	defer queuing.Cleanup()
	r := router.NewRouter()

	// Start the HTTP server
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("HTTP server error: %v\n", err)
	}

}
