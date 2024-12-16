package model

type QueueMessage struct {
	QueueName string   `json:"queue_name"`
	ToClients []string `json:"to_clients"`
	Data      string   `json:"data"`
}
