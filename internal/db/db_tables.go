package db

import (
	"database/sql"
	"log"
)

const (
	dbPortfolio = "./portfolio.db"
	dbCache     = "./cache.db"
)

func InitTables() {
	log.Println("Starting init tables...")
	dbPortfolio, err := sql.Open("sqlite3", dbPortfolio)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPortfolio.Close()

	_, err = dbPortfolio.Exec(`
			CREATE TABLE IF NOT EXISTS portfolio (
				chat_id INTEGER NOT NULL,
				operations JSONB,
				money_operations JSONB,
				updated_at TEXT NOT NULL,
				time_period TEXT,
				PRIMARY KEY (chat_id)
			)
		`)
	if err != nil {
		log.Fatal(err)
	}

	dbCache, err := sql.Open("sqlite3", dbCache)
	if err != nil {
		log.Fatal(err)
	}
	defer dbCache.Close()
	_, err = dbCache.Exec(`
			CREATE TABLE IF NOT EXISTS cache (
			    key_name TEXT,
				dohod JSONB,
				moex_isin JSONB,
				moex_stock_bond JSONB,
				moex_stock_share JSONB,
				PRIMARY KEY (key_name)
			)
		`)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Init tables finished.")
}
