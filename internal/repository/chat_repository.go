package repository

import (
	"github.com/google/uuid"
	"message-broker/internal/db"
	"message-broker/internal/utils"
	"message-broker/model"
	"strings"
)

func GetChatsByUserId(userId string) ([]*model.Chat, error) {

	chats, err := db.FindChatsForUser(userId)
	if err != nil {
		return nil, err
	}
	if chats == nil {
		return []*model.Chat{}, nil
	}
	return chats, nil
}

func GetChatById(chatId string) (*model.Chat, error) {
	return nil, nil
}

func CreateChat(chat *model.Chat) error {
	if strings.TrimSpace(chat.ID) == "" {
		chat.ID = uuid.New().String()
	}
	chat.CreatedAt = utils.GetCurrentTime()
	err := db.InsertChat(chat)
	if err != nil {
		return err
	}
	return nil
}
