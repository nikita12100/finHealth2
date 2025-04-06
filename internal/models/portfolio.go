package models

type Portfolio struct {
	ChatId     int64       `json:"chat_id"`
	Name       string      `json:"name"`
	Operations []Operation `json:"operations"`
}
