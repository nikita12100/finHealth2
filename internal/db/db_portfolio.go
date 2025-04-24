package db

import (
	"database/sql"
	"encoding/json"
	"errors"
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
		INSERT INTO portfolio (chat_id, operations, money_operations, updated_at, time_period)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(chat_id) DO UPDATE SET
			operations = excluded.operations,
			money_operations = excluded.money_operations,
			updated_at = excluded.updated_at,
			time_period = excluded.time_period
	`, p.ChatId, operationsJSON, moneyOperationsJSON, p.UpdatedAt.Format(time.RFC3339), timePeriod)

	return err
}

func GetPortfolioOrCreate(chatId int64) models.Portfolio {
	portfolio, err := getPortfolio(chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Not found portfolio, create new one", "chatId", chatId)
			return models.Portfolio{}
		} else {
			slog.Error("Error while fetch portfolio, create new one", "chatId", chatId, "error", err)
			return models.Portfolio{}
		}
	}
	return portfolio
}

func getPortfolio(chatId int64) (models.Portfolio, error) {
	slog.Info("Get portfolio ...", "chatId", chatId)
	db, _ := sql.Open("sqlite3", dbPortfolio)

	var p models.Portfolio
	var operationsJSON []byte
	var moneyOperationsJSON []byte
	var updatedAtStr string
	var timePeriodJSON []byte

	err := db.QueryRow(`
		SELECT chat_id, operations, money_operations, updated_at, time_period
		FROM portfolio
		WHERE chat_id = ?
	`, chatId).Scan(&p.ChatId, &operationsJSON, &moneyOperationsJSON, &updatedAtStr, &timePeriodJSON)
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
