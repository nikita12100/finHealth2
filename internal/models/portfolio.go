package models

import "time"

type DateRange struct {
	Start time.Time `json:"from"`
	End   time.Time `json:"to"`
}

type Portfolio struct {
	ChatId          int64            `json:"chat_id"`
	Operations      []Operation      `json:"operations"`
	MoneyOperations []MoneyOperation `json:"money_operations"`
	UpdatedAt       time.Time        `json:"updated_at"`
	TimePeriod      DateRange        `json:"time_period"`
}

func (old DateRange) ExtendTimePeriod(new DateRange) DateRange {
	left := old.Start
	if old.Start.IsZero() || new.Start.Before(old.Start) {
		left = new.Start
	}
	right := old.End
	if old.End.IsZero() || new.End.After(old.End) {
		right = new.End
	}
	return DateRange{
		Start: left,
		End:   right,
	}
}