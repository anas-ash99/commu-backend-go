package router

import (
	"github.com/go-chi/chi/v5"
	"message-broker/internal/handler"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()
	// chats
	r.Post("/chat", handler.CreateChat)
	r.Get("/chat/{userId}", handler.GetChatsForUser)
	r.Post("/message", handler.CreateMessage)
	r.Get("/message", handler.GetMessageByKey)

	// users
	r.Post("/user", handler.UpsertUser)
	r.Get("/user", handler.GetUsers)
	r.Patch("/user/friend/{userId}/{friendId}", handler.AddFriend)
	r.Get("/user/friend/{userId}", handler.GetFriends)

	r.Get("/ws", handler.HandleWebsocketConnections)
	r.Get("/ws-clients", handler.GetConnectedClients)
	return r
}
