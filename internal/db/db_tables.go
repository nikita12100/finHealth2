package db

import (
	"database/sql"
	"log"
)

const (
	dbName = "./portfolio.db"
)

func InitTables() {
	log.Println("Starting init tables...")
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS portfolio (
				chat_id INTEGER NOT NULL,
				name TEXT NOT NULL,
				operations JSONB,
				PRIMARY KEY (chat_id, name)
			)
		`)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Init tables finished.")
}
