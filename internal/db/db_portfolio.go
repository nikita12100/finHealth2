package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"sort"
	"test2/internal/models"
)

func SavePortfolio(p models.Portfolio) error {
	log.Printf("SavePortfolio chatId=%v...\n", p.ChatId)
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

	_, err = db.Exec(`
		INSERT INTO portfolio (chat_id, name, operations, money_operations)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(chat_id, name) DO UPDATE SET
			operations = excluded.operations,
			money_operations = excluded.money_operations
	`, p.ChatId, p.Name, operationsJSON, moneyOperationsJSON)

	return err
}

func GetPortfolio(chatId int64, name string) (models.Portfolio, error) {
	log.Printf("GetPortfolio chatId=%v...\n", chatId)
	db, _ := sql.Open("sqlite3", dbPortfolio)

	var p models.Portfolio
	var operationsJSON []byte
	var moneyOperationsJSON []byte

	err := db.QueryRow(`
		SELECT chat_id, name, operations, money_operations
		FROM portfolio
		WHERE chat_id = ? AND name = ?
	`, chatId, name).Scan(&p.ChatId, &p.Name, &operationsJSON, &moneyOperationsJSON)
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

	return p, nil
}
