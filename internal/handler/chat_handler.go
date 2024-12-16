package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"message-broker/internal/repository"
	"message-broker/model"
	"net/http"
)

type Chat interface {
}

func GetChatsForUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := chi.URLParam(r, "userId")
	chats, err := repository.GetChatsByUserId(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	response, _ := json.Marshal(chats)

	w.Write(response)

}

func CreateChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var chat *model.Chat
	var err error

	err = json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = repository.CreateChat(chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res, err := json.Marshal(chat)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)

}
