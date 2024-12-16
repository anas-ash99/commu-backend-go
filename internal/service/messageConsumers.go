package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"message-broker/internal/db"
	"message-broker/internal/utils"
	"message-broker/internal/ws"
	"message-broker/model"
	"strings"
)

func CreateMessage(msg string) {

	queueMsg, err := utils.ParseJSON[model.QueueMessage](msg)
	if err != nil {
		return
	}
	message, err := utils.ParseJSON[model.Message](queueMsg.Data)
	if err != nil {
		fmt.Printf("error while parsing json %v\n", err.Error())
		return
	}
	if strings.TrimSpace(message.ID) == "" {
		message.ID = uuid.New().String()
	}
	message.DeliveredAt = utils.GetCurrentTime()
	err = db.InsertMessage(message)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	str, err := json.Marshal(message)
	var wsMessage = model.WsMessage{
		Type: "chat_message",
		Data: string(str),
	}
	wsStr, _ := json.Marshal(wsMessage)
	for _, client := range queueMsg.ToClients {
		ws.SendTooAllPlatforms(client, string(wsStr))
	}
}
