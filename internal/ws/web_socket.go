package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// ChatClients A thread-safe map to store active WebSocket connections
var ChatClients = sync.Map{} // Map of client ID to WebSocket connection

type WriteMessageCallBack func(msg string, clients []string)

func WriteMessage(message string, clientID string) {
	value, ok := ChatClients.Load(clientID)
	if !ok {
		fmt.Errorf("client %v not found", clientID)
	}

	if value == nil {
		fmt.Errorf("client %v not found", clientID)
		return
	}
	conn := value.(*websocket.Conn)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Printf("WebSocket write error for client %v: %v\n", clientID, err)
		conn.Close()
		ChatClients.Delete(clientID) // Remove the client if the connection fails

	}
	fmt.Printf("WebSocket write success for client %v\n", clientID)

}

func SendTooAllPlatforms(clientId string, message string) {
	android := clientId + "-android"
	ios := clientId + "-ios"
	web := clientId + "-web"
	desktop := clientId + "-desktop"
	WriteMessage(message, android)
	WriteMessage(message, ios)
	WriteMessage(message, web)
	WriteMessage(message, desktop)
}

func GetConnectedClients() []string {
	var clientsIds []string
	ChatClients.Range(func(key, value interface{}) bool {
		clientID := key.(string) // Convert key to string
		clientsIds = append(clientsIds, clientID)
		return true
	})
	return clientsIds
}

func BroadcastMessage(msg []byte) {
	ChatClients.Range(func(key, value interface{}) bool {
		conn := value.(*websocket.Conn)
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Printf("Error broadcasting message to %v: %v\n", key, err)
		}
		return true // continue iteration
	})
}
