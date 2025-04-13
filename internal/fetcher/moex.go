package fetcher

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"test2/internal/db"
	"test2/internal/models"
	"time"
)

const (
	ttlMOEX      = 1 * time.Hour
	ttlMOEXStock = 1 * time.Hour
)

func GetTickerByISINCached(isin string) (string, error) {
	if entry, err := db.GetCacheMoexIsin(isin); err == nil {
		if time.Since(entry.Created) < ttlMOEX {
			return entry.Value, nil
		}
	}

	value, _ := getTickerByISIN(isin)

	db.SaveCacheMoexIsin(isin, db.CacheMoexIsin{
		Value:   value,
		Created: time.Now(),
	})

	return value, nil
}

func getTickerByISIN(isin string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://iss.moex.com/iss/securities.json?q=%s", isin))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	type response struct {
		Securities struct {
			Columns []string        `json:"columns"`
			Data    [][]interface{} `json:"data"`
		} `json:"securities"`
	}

	var responseT response
	if err := json.NewDecoder(resp.Body).Decode(&responseT); err != nil {
		return "", err
	}

	for i, name := range responseT.Securities.Columns {
		if name == "secid" {
			ticker := removeFixPrefix(responseT.Securities.Data[0][i].(string)) //FIXOZON -> OZON
			return fmt.Sprintf("%v", ticker), nil
		}
	}

	return "", fmt.Errorf("ticker not found")
}

func removeFixPrefix(input string) string {
	if strings.HasPrefix(input, "FIX") {
		return strings.TrimPrefix(input, "FIX")
	}
	return input
}

func GetLastPriceBondCached(ticker string) (models.StockBondInfo, error) {
	if entry, err := db.GetCacheMoexStock(ticker); err == nil {
		if time.Since(entry.Created) < ttlMOEXStock {
			return entry.Value, nil
		}
	}

	value, _ := getLastPriceBond(ticker)

	db.SaveCacheMoexStock(ticker, db.CacheMoexStock{
		Value:   value,
		Created: time.Now(),
	})

	return value, nil
}

func getLastPriceBond(ticker string) (models.StockBondInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://iss.moex.com/iss/engines/stock/markets/bonds/securities/%s.json", ticker))
	if err != nil {
		slog.Error("Got error while GET url", ticker, err)
		return models.StockBondInfo{}, err
	}
	defer resp.Body.Close()

	type response struct {
		Marketdata struct {
			Columns []string        `json:"columns"`
			Data    [][]interface{} `json:"data"`
		} `json:"securities"`
	}

	var responseT response
	if err := json.NewDecoder(resp.Body).Decode(&responseT); err != nil {
		slog.Error("Failed parse to json answer", ticker, err)
		return models.StockBondInfo{}, err
	}

	var tqbrIndex int
	var sType string
	if ticker[:2] == "SU" {
		sType = "TQOB"
	} else {
		sType = "TQCB"
	}
	for i, name := range responseT.Marketdata.Columns {
		if name == "BOARDID" {
			for j := range len(responseT.Marketdata.Data) {
				if responseT.Marketdata.Data[j][i] == sType {
					tqbrIndex = j
					break
				}
			}
		}
	}
	var couponValue, lastPrice, faceValue, couponPeriod float64
	var faceUnit string
	for i, name := range responseT.Marketdata.Columns {
		if name == "COUPONVALUE" {
			if len(responseT.Marketdata.Data) > 0 {
				couponValue = responseT.Marketdata.Data[tqbrIndex][i].(float64)
			} else {
				slog.Error("Not found COUPONVALUE in response", "ticker", ticker)
			}
		} else if name == "PREVPRICE" {
			if len(responseT.Marketdata.Data) > 0 {
				lastPrice = responseT.Marketdata.Data[tqbrIndex][i].(float64)
			} else {
				slog.Error("Not found PREVPRICE in response", "ticker", ticker)
			}
		} else if name == "FACEVALUE" {
			if len(responseT.Marketdata.Data) > 0 {
				faceValue = responseT.Marketdata.Data[tqbrIndex][i].(float64)
			} else {
				slog.Error("Not found FACEVALUE in response", "ticker", ticker)
			}
		} else if name == "COUPONPERIOD" {
			if len(responseT.Marketdata.Data) > 0 {
				couponPeriod = responseT.Marketdata.Data[tqbrIndex][i].(float64)
			} else {
				slog.Error("Not found COUPONPERIOD in response", "ticker", ticker)
			}
		} else if name == "FACEUNIT" {
			if len(responseT.Marketdata.Data) > 0 {
				faceUnit = responseT.Marketdata.Data[tqbrIndex][i].(string)
			} else {
				slog.Error("Not found FACEUNIT in response", "ticker", ticker)
			}
		}
	}
	return models.StockBondInfo{
		CouponValue:  couponValue,
		LastPrice:    lastPrice,
		FaceValue:    faceValue,
		CouponPeriod: couponPeriod,
		FaceUnit:     faceUnit,
	}, nil
}

func GetLastPriceShare(ticker string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://iss.moex.com/iss/engines/stock/markets/shares/securities/%s.json", ticker))
	if err != nil {
		slog.Error("Got error while GET url", ticker, err)
		return 0.0, err
	}
	defer resp.Body.Close()

	type response struct {
		Marketdata struct {
			Columns []string        `json:"columns"`
			Data    [][]interface{} `json:"data"`
		} `json:"securities"`
	}

	var responseT response
	if err := json.NewDecoder(resp.Body).Decode(&responseT); err != nil {
		slog.Error("Failed parse to json answer", "ticker", ticker, err)
		return 0.0, err
	}

	var tqbrIndex int
	for i, name := range responseT.Marketdata.Columns {
		if name == "BOARDID" {
			for j := range len(responseT.Marketdata.Data) {
				if responseT.Marketdata.Data[j][i] == "TQBR" {
					tqbrIndex = j
					break
				}
			}
		}
	}
	var lastPrice float64
	for i, name := range responseT.Marketdata.Columns {
		if name == "PREVPRICE" {
			if len(responseT.Marketdata.Data) > 0 {
				lastPrice = responseT.Marketdata.Data[tqbrIndex][i].(float64)
			} else {
				slog.Error("Not found PREVPRICE in response", "ticker", ticker)
			}
		}
	}
	return lastPrice, nil
}
