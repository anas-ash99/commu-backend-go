package handler

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"message-broker/internal/repository"
	"message-broker/internal/ws"
	"message-broker/model"
	"net/http"
)

// Upgrader to handle WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections for simplicity
		return true
	},
}

// HandleWebsocketConnections upgrades the HTTP connection to WebSocket and listens for messages
func HandleWebsocketConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v\n", err)
		return
	}
	defer conn.Close()

	clientID := r.URL.Query().Get("clientID")
	// Store the connection associated with the client ID
	ws.ChatClients.Store(clientID, conn)
	log.Printf("Client %v connected!", clientID)

	// Clean up when the client disconnects
	defer func() {
		log.Printf("Client %v disconnected!", clientID)
		ws.ChatClients.Delete(clientID)
	}()

	// Listen for incoming WebSocket messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v\n", err)
			break
		}
		//log.Printf("Received message from WebSocket: %s\n", msg)
		var wsMessage model.WsMessage
		err = json.Unmarshal(msg, &wsMessage)
		if err != nil {
			log.Printf("WebSocket unmarshal error: %v\n", err)
		}
		log.Printf("WebSocket received message: %v\n", wsMessage)
		if wsMessage.Type == "update_user" {
			handleUpdateUser(wsMessage.Data)
		}
	}
}

func handleUpdateUser(userJson string) {
	var user model.User
	err := json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		log.Printf("WebSocket unmarshal error: %v\n", err)
		return
	}
	log.Printf("WebSocket received user: %v\n", user)
	err = repository.UpdateUser(&user)
	if err != nil {
		log.Printf("WebSocket update error: %v\n", err)
		return
	}

	for _, friend := range user.Friends {
		var wsMessage model.WsMessage = model.WsMessage{
			Type: "update_user",
			Data: userJson,
		}
		re, err := json.Marshal(wsMessage)
		if err != nil {
			log.Printf("WebSocket marshal error: %v\n", err)
			return
		}
		ws.SendTooAllPlatforms(friend, string(re))

	}

}

func GetConnectedClients(w http.ResponseWriter, r *http.Request) {
	clientsIds := ws.GetConnectedClients()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(clientsIds); err != nil {
		http.Error(w, "Failed to encode client IDs", http.StatusInternalServerError)
	}
}

func HandleUsersWebsocketConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v\n", err)
		return
	}
	defer conn.Close()

	clientID := r.URL.Query().Get("clientID")
	// Store the connection associated with the client ID
	ws.ChatClients.Store(clientID, conn)
	log.Printf("Client %v connected!", clientID)

	// Clean up when the client disconnects
	defer func() {
		log.Printf("Client %v disconnected!", clientID)
		ws.ChatClients.Delete(clientID)
	}()

	// Listen for incoming WebSocket messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v\n", err)
			break
		}
		log.Printf("Received message from WebSocket: %s\n", msg)
		var queueMsg model.QueueMessage
		err = json.Unmarshal(msg, &queueMsg)
		ws.BroadcastMessage(msg)

	}

}
