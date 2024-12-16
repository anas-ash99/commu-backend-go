package repository

import (
	"message-broker/internal/db"
	"message-broker/model"
)

func GetMessagesByKey(value string, key string) ([]*model.Message, error) {

	messages, err := db.FineMessagesByField(value, key)
	if messages == nil {
		messages = []*model.Message{}
	}
	return messages, err
}
