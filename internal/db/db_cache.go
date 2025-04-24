package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"test2/internal/common"
	"test2/internal/models"
)

func saveCacheDB[K any, CacheV any](key K, value CacheV, nameValueDB string) error {
	db, err := sql.Open("sqlite3", dbCache)
	if err != nil {
		return err
	}

	cacheJSON, err := json.Marshal(value)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
		INSERT INTO cache (ticker, %s)
		VALUES (?, ?)
		ON CONFLICT(ticker) DO UPDATE SET
			%s = excluded.%s
	`, nameValueDB, nameValueDB, nameValueDB)

	_, err = db.Exec(query, key, cacheJSON)

	return err
}

func SaveCacheDohod(ticker string, cache common.Cache[float64]) error {
	return saveCacheDB(ticker, cache, "dohod")
}

func SaveCacheMoexIsin(isin string, cache common.Cache[string]) error {
	return saveCacheDB(isin, cache, "moex_isin")
}

func SaveCacheMoexStockBond(ticker string, cache common.Cache[models.StockBondInfo]) error {
	return saveCacheDB(ticker, cache, "moex_stock_bond")
}

func SaveCacheMoexStockShare(ticker string, cache common.Cache[float64]) error {
	return saveCacheDB(ticker, cache, "moex_stock_share")
}

func getCache[K any, CacheV any](ticker K, nameValueDB string) (CacheV, error) {
	db, _ := sql.Open("sqlite3", dbCache)

	var cache CacheV
	var cacheJSON []byte
	query := fmt.Sprintf(`
		SELECT %s
		FROM cache
		WHERE ticker = ?
	`, nameValueDB)
	err := db.QueryRow(query, ticker).Scan(&cacheJSON)
	if err != nil {
		return cache, err
	}

	err = json.Unmarshal(cacheJSON, &cache)
	if err != nil {
		return cache, err
	}

	return cache, nil
}

func GetCacheDohod(ticker string) (common.Cache[float64], error) {
	return getCache[string, common.Cache[float64]](ticker, "dohod")
}

func GetCacheMoexIsin(isin string) (common.Cache[string], error) {
	return getCache[string, common.Cache[string]](isin, "moex_isin")
}

func GetCacheMoexStockBond(ticker string) (common.Cache[models.StockBondInfo], error) {
	return getCache[string, common.Cache[models.StockBondInfo]](ticker, "moex_stock_bond")
}

func GetCacheMoexStockShare(ticker string) (common.Cache[float64], error) {
	return getCache[string, common.Cache[float64]](ticker, "moex_stock_share")
}
