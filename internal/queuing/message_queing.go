package queuing

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

var (
	amqpConn    *amqp.Connection
	amqpChannel *amqp.Channel
)

func PublishToQueue(queueName string, message string) error {
	if amqpChannel == nil {
		return fmt.Errorf("AMQP channel is not initialized")
	}

	// Publish the message
	err := amqpChannel.Publish(
		"",        // Exchange
		queueName, // Routing key
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return fmt.Errorf("publish error: %v", err)
	}

	log.Printf("Message successfully published to queue '%s'\n", queueName)
	return nil
}

// SetupAMQP initializes the RabbitMQ connection and channel
func SetupAMQP(amqpURL string) error {
	var err error
	// Connect to RabbitMQ
	amqpConn, err = amqp.Dial(amqpURL)
	if err != nil {
		return fmt.Errorf("AMQP connection error: %v", err)
	}

	// Create a channel
	amqpChannel, err = amqpConn.Channel()
	if err != nil {
		return fmt.Errorf("AMQP channel error: %v", err)
	}
	log.Printf("AMQP connection opened\n")
	return nil
}

type Publish func(string, string) error

// ConsumeMessages listens for messages from a RabbitMQ queue and sends them to WebSocket
func ConsumeMessages(queueName string, preformAction func(data string)) {
	if amqpChannel == nil {
		log.Println("AMQP channel is not initialized")
		return
	}

	// Declare the queue
	_, err := amqpChannel.QueueDeclare(
		queueName,
		false, // Durable
		false, // Auto-deleted
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	// Consume messages from the queue
	msgs, err := amqpChannel.Consume(
		queueName,
		"",    // Consumer tag
		true,  // Auto-ack
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Printf("Queue consume error: %v\n", err)
		return
	}

	go func() {
		for d := range msgs {
			log.Printf("Received message from queue '%s': %s\n", queueName, d.Body)
			preformAction(string(d.Body))
		}
	}()
}

// Cleanup releases resources
func Cleanup() {
	if amqpChannel != nil {
		amqpChannel.Close()
	}
	if amqpConn != nil {
		amqpConn.Close()
	}
	log.Println("AMQP resources cleaned up")
}
