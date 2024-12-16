package model

type WsMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}
