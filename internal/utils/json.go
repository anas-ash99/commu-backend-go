package utils

import (
	"encoding/json"
	"net/http"
)

// ParseJSON Generic function to parse JSON into a struct
func ParseJSON[T any](jsonString string) (*T, error) {
	var result T
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RespondJSON sends a JSON response
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}
