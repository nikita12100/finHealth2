package db

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"sort"
	"test2/internal/models"
	"time"
)

func SavePortfolio(p models.Portfolio) error {
	slog.Info("SavePortfolio...", "chatId", p.ChatId)
	db, err := sql.Open("sqlite3", dbPortfolio)
	if err != nil {
		return err
	}

	operations := p.Operations
	sort.Slice(operations, func(i, j int) bool {
		return operations[i].Date.Before(operations[j].Date)
	})
	operationsJSON, err := json.Marshal(operations)
	if err != nil {
		return err
	}

	moneyOperations := p.MoneyOperations
	sort.Slice(moneyOperations, func(i, j int) bool {
		return moneyOperations[i].Time.Before(moneyOperations[j].Time)
	})
	moneyOperationsJSON, err := json.Marshal(moneyOperations)
	if err != nil {
		return err
	}
	timePeriod, err := json.Marshal(p.TimePeriod)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO portfolio (chat_id, name, operations, money_operations, updated_at, time_period)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(chat_id, name) DO UPDATE SET
			operations = excluded.operations,
			money_operations = excluded.money_operations,
			updated_at = excluded.updated_at,
			time_period = excluded.time_period
	`, p.ChatId, p.Name, operationsJSON, moneyOperationsJSON, p.UpdatedAt.Format(time.RFC3339), timePeriod)

	return err
}

func GetPortfolio(chatId int64, name string) (models.Portfolio, error) {
	slog.Info("GetPortfolio...", "chatId", chatId)
	db, _ := sql.Open("sqlite3", dbPortfolio)

	var p models.Portfolio
	var operationsJSON []byte
	var moneyOperationsJSON []byte
	var updatedAtStr string
	var timePeriodJSON []byte

	err := db.QueryRow(`
		SELECT chat_id, name, operations, money_operations, updated_at, time_period
		FROM portfolio
		WHERE chat_id = ? AND name = ?
	`, chatId, name).Scan(&p.ChatId, &p.Name, &operationsJSON, &moneyOperationsJSON, &updatedAtStr, &timePeriodJSON)
	if err != nil {
		return p, err
	}

	err = json.Unmarshal(operationsJSON, &p.Operations)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(moneyOperationsJSON, &p.MoneyOperations)
	if err != nil {
		return p, err
	}
	p.UpdatedAt, err = time.Parse(time.RFC3339, updatedAtStr)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(timePeriodJSON, &p.TimePeriod)
	if err != nil {
		return p, err
	}

	return p, nil
}
