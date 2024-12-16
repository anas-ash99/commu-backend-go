package handler

import (
	"encoding/json"
	"message-broker/internal/queuing"
	"message-broker/internal/repository"
	"message-broker/model"
	"net/http"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {

	var msg model.QueueMessage
	err := json.NewDecoder(r.Body).Decode(&msg)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	msgString, err := json.Marshal(msg)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = queuing.PublishToQueue("SaveMessageQueue", string(msgString))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)

}

func GetMessageByKey(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	value := query.Get("value")
	if value == "" {
		http.Error(w, "value is empty", http.StatusBadRequest)
		return
	}
	key := query.Get("key")
	if key == "" {
		http.Error(w, "key is empty", http.StatusBadRequest)
		return
	}
	messages, err := repository.GetMessagesByKey(value, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	response, err := json.Marshal(messages)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
