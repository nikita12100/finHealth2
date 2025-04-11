package db

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"test2/internal/models"
	"time"
)

type CacheDohod struct {
	Value   float64   `json:"value"`
	Created time.Time `json:"created"`
}

type CacheMoexIsin struct {
	Value   string    `json:"value"`
	Created time.Time `json:"created"`
}

type CacheMoexStock struct {
	Value   models.StockBondInfo `json:"value"`
	Created time.Time            `json:"created"`
}

func SaveCacheDohod(ticker string, cache CacheDohod) error {
	slog.Debug("SavingCacheDohod...")
	db, err := sql.Open("sqlite3", dbCache)
	if err != nil {
		return err
	}

	cacheJSON, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO cache (ticker, dohod)
		VALUES (?, ?)
		ON CONFLICT(ticker) DO UPDATE SET
			dohod = excluded.dohod
	`, ticker, cacheJSON)

	return err
}

func SaveCacheMoexIsin(isin string, cache CacheMoexIsin) error {
	slog.Debug("SaveCacheMoexIsin...")
	db, err := sql.Open("sqlite3", dbCache)
	if err != nil {
		return err
	}

	cacheJSON, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO cache (ticker, moex_isin)
		VALUES (?, ?)
		ON CONFLICT(ticker) DO UPDATE SET
			moex_isin = excluded.moex_isin
	`, isin, cacheJSON)

	return err
}

func SaveCacheMoexStock(ticker string, cache CacheMoexStock) error {
	slog.Debug("SaveCacheMoexStock...")
	db, err := sql.Open("sqlite3", dbCache)
	if err != nil {
		return err
	}

	cacheJSON, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO cache (ticker, moex_stock)
		VALUES (?, ?)
		ON CONFLICT(ticker) DO UPDATE SET
			moex_stock = excluded.moex_stock
	`, ticker, cacheJSON)

	return err
}

func GetCacheDohod(ticker string) (CacheDohod, error) {
	slog.Debug("GetCacheDohod...")
	db, _ := sql.Open("sqlite3", dbCache)

	var cache CacheDohod
	var cacheJSON []byte
	err := db.QueryRow(`
		SELECT dohod
		FROM cache
		WHERE ticker = ?
	`, ticker).Scan(&cacheJSON)
	if err != nil {
		return cache, err
	}

	err = json.Unmarshal(cacheJSON, &cache)
	if err != nil {
		return cache, err
	}

	return cache, nil
}

func GetCacheMoexIsin(isin string) (CacheMoexIsin, error) {
	slog.Debug("GetCacheMoexIsin...")
	db, _ := sql.Open("sqlite3", dbCache)

	var cache CacheMoexIsin
	var cacheJSON []byte
	err := db.QueryRow(`
		SELECT moex_isin
		FROM cache
		WHERE ticker = ?
	`, isin).Scan(&cacheJSON)
	if err != nil {
		return cache, err
	}

	err = json.Unmarshal(cacheJSON, &cache)
	if err != nil {
		return cache, err
	}

	return cache, nil
}

func GetCacheMoexStock(ticker string) (CacheMoexStock, error) {
	slog.Debug("GetCacheMoexStock...")
	db, _ := sql.Open("sqlite3", dbCache)

	var cache CacheMoexStock
	var cacheJSON []byte
	err := db.QueryRow(`
		SELECT moex_stock
		FROM cache
		WHERE ticker = ?
	`, ticker).Scan(&cacheJSON)
	if err != nil {
		return cache, err
	}

	err = json.Unmarshal(cacheJSON, &cache)
	if err != nil {
		return cache, err
	}

	return cache, nil
}
