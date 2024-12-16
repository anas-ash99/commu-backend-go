package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"message-broker/internal/repository"
	"message-broker/model"
	"net/http"
)

func UpsertUser(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = repository.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func GetUser(w http.ResponseWriter, r *http.Request) {

}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repository.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func AddFriend(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("userId is required"))
		return
	}

	friendId := chi.URLParam(r, "friendId")
	if friendId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("friendId is required"))
		return
	}
	if userId == friendId {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user id and friend id cannot be the same"))
		return
	}
	err := repository.AddFriend(userId, friendId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetFriends(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	fmt.Println(userId)
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("userId is required"))
		return
	}
	users, err := repository.GetFriendsForUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	res, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
