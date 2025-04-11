package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"test2/internal/models"
)

func SavePortfolio(p models.Portfolio) error {
	log.Printf("SavePortfolio chatId=%v...\n", p.ChatId)
	db, err := sql.Open("sqlite3", dbPortfolio)
	if err != nil {
		return err
	}

	operationsJSON, err := json.Marshal(p.Operations)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO portfolio (chat_id, name, operations)
		VALUES (?, ?, ?)
		ON CONFLICT(chat_id, name) DO UPDATE SET
			operations = excluded.operations
	`, p.ChatId, p.Name, operationsJSON)

	return err
}

func GetPortfolio(chatId int64, name string) (models.Portfolio, error) {
	log.Printf("GetPortfolio chatId=%v...\n", chatId)
	db, _ := sql.Open("sqlite3", dbPortfolio)

	var p models.Portfolio
	var operationsJSON []byte

	err := db.QueryRow(`
		SELECT chat_id, name, operations
		FROM portfolio
		WHERE chat_id = ? AND name = ?
	`, chatId, name).Scan(&p.ChatId, &p.Name, &operationsJSON)
	if err != nil {
		return p, err
	}

	err = json.Unmarshal(operationsJSON, &p.Operations)
	if err != nil {
		return p, err
	}

	return p, nil
}
