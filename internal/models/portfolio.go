package models

import "time"

type FromTo struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Portfolio struct {
	ChatId          int64            `json:"chat_id"`
	Name            string           `json:"name"`
	Operations      []Operation      `json:"operations"`
	MoneyOperations []MoneyOperation `json:"money_operations"`
	UpdatedAt       time.Time        `json:"updated_at"`
	TimePeriod      FromTo           `json:"time_period"`
}
