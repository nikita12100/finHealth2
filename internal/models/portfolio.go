package models

const (
	WEIGHT_NORM = 60000.0
)

type Portfolio struct {
	ChatId     int64       `json:"chat_id"`
	Name       string      `json:"name"`
	Operations []Operation `json:"operations"`
}
