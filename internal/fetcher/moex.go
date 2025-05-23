package fetcher

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"test2/internal/common"
	"test2/internal/db"
	"test2/internal/models"
	"time"
)

const (
	ttlMOEX      = 1 * time.Hour
	ttlMOEXStock = 1 * time.Hour
)

func GetTickerByISINCached(isin string) (string, error) {
	return common.Cached(isin, ttlMOEX, db.GetCacheMoexIsin, getTickerByISIN, db.SaveCacheMoexIsin)
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
	return common.Cached(ticker, ttlMOEXStock, db.GetCacheMoexStockBond, getLastPriceBond, db.SaveCacheMoexStockBond)
}

func getLastPriceBond(ticker string) (models.StockBondInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://iss.moex.com/iss/engines/stock/markets/bonds/securities/%s.json?iss.meta=off&iss.only=securities&securities.columns=BOARDID,COUPONVALUE,FACEVALUE,COUPONPERIOD,FACEUNIT,PREVPRICE", ticker))
	if err != nil {
		slog.Error("Got error while GET url", ticker, err)
		return models.StockBondInfo{}, err
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
	if len(responseT.Securities.Data) == 0 || len(responseT.Securities.Data[0]) == 0 {
		slog.Error("wrong answer format moex Bond", "ticker", ticker)
		return models.StockBondInfo{}, err
	}
	for i := range len(responseT.Securities.Data) {
		if responseT.Securities.Data[i][0] == sType {
			tqbrIndex = i
			break
		}
	}

	if len(responseT.Securities.Data[tqbrIndex]) < 5 || len(responseT.Securities.Data[tqbrIndex]) == 0 {
		slog.Error("wrong answer format2 moex Bond", "ticker", ticker)
		return models.StockBondInfo{}, err
	}
	couponValue := responseT.Securities.Data[tqbrIndex][1].(float64)
	faceValue := responseT.Securities.Data[tqbrIndex][2].(float64)
	couponPeriod := responseT.Securities.Data[tqbrIndex][3].(float64)
	faceUnit := responseT.Securities.Data[tqbrIndex][4].(string)
	lastPrice := responseT.Securities.Data[tqbrIndex][5].(float64)

	return models.StockBondInfo{
		CouponValue:  couponValue,
		LastPrice:    lastPrice,
		FaceValue:    faceValue,
		CouponPeriod: couponPeriod,
		FaceUnit:     faceUnit,
	}, nil
}

func GetLastPriceShareCached(ticker string) (float64, error) {
	return common.Cached(ticker, ttlMOEXStock, db.GetCacheMoexStockShare, getLastPriceShare, db.SaveCacheMoexStockShare)
}

func getLastPriceShare(ticker string) (float64, error) {
	otcTickers := []string{"UDMN"}
	if slices.Contains(otcTickers, ticker) {
		return getLastPriceShareOTC(ticker)
	} else {
		return getLastPriceShareStock(ticker)
	}
}

func getLastPriceShareStock(ticker string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://iss.moex.com/iss/engines/stock/markets/shares/securities/%s.json?iss.meta=off&iss.only=securities,marketdata&marketdata.columns=BOARDID,LAST&securities.columns=BOARDID,PREVPRICE", ticker))
	if err != nil {
		slog.Error("Got error while GET url", ticker, err)
		return 0.0, err
	}
	defer resp.Body.Close()

	type response struct {
		Securities struct {
			Columns []string        `json:"columns"`
			Data    [][]interface{} `json:"data"`
		} `json:"securities"`
		Marketdata struct {
			Columns []string        `json:"columns"`
			Data    [][]interface{} `json:"data"`
		} `json:"marketdata"`
	}

	var responseT response
	if err := json.NewDecoder(resp.Body).Decode(&responseT); err != nil {
		slog.Error("Failed parse to json answer", "ticker", ticker, "error", err)
		return 0.0, err
	}

	if len(responseT.Marketdata.Data) == 0 {
		slog.Error("Empty marketdata for moex asnwer Share", "ticker", ticker)
		return 0.0, err
	}
	var tqbrIndex int
	for i := range len(responseT.Marketdata.Data) {
		if responseT.Marketdata.Data[i][0] == "TQBR" {
			tqbrIndex = i
			break
		}
	}

	if len(responseT.Marketdata.Data[tqbrIndex]) == 0 {
		slog.Error("wrong answer format2 moex Share", "ticker", ticker)
		return 0.0, err
	}

	if responseT.Marketdata.Data[tqbrIndex][1] == nil { // use securities
		return responseT.Securities.Data[tqbrIndex][1].(float64), nil
	} else {
		return responseT.Marketdata.Data[tqbrIndex][1].(float64), nil
	}
}

func getLastPriceShareOTC(ticker string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://iss.moex.com/iss/engines/otc/markets/shares/securities/%s.json?iss.meta=off&iss.only=securities&securities.columns=BOARDID,PREVPRICE", ticker))
	if err != nil {
		slog.Error("Got error while GET url", ticker, err)
		return 0.0, err
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
		slog.Error("Failed parse to json answer", "ticker", ticker, "error", err)
		return 0.0, err
	}

	if len(responseT.Securities.Data) == 0 {
		slog.Error("Empty securities for otc moex asnwer Share", "ticker", ticker)
		return 0.0, err
	}
	var mtqrIndex int
	for i := range len(responseT.Securities.Data) {
		if responseT.Securities.Data[i][0] == "MTQR" {
			mtqrIndex = i
			break
		}
	}

	if len(responseT.Securities.Data[mtqrIndex]) == 0 {
		slog.Error("wrong answer format2 moex Share", "ticker", ticker)
		return 0.0, err
	}
	return responseT.Securities.Data[mtqrIndex][1].(float64), nil
}

func GetLastPriceTOM(ticker string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://iss.moex.com/iss/engines/currency/markets/selt/boards/CETS/securities/%s.json?iss.meta=off&iss.only=marketdata&marketdata.columns=LAST", ticker))
	if err != nil {
		slog.Error("Got error while GET url", ticker, err)
		return 0.0, err
	}
	defer resp.Body.Close()

	type response struct {
		Marketdata struct {
			Columns []string        `json:"columns"`
			Data    [][]interface{} `json:"data"`
		} `json:"marketdata"`
	}

	var responseT response
	if err := json.NewDecoder(resp.Body).Decode(&responseT); err != nil {
		slog.Error("Failed parse to json answer", "ticker", ticker, "error", err)
		return 0.0, err
	}

	if len(responseT.Marketdata.Data) == 0 || len(responseT.Marketdata.Data[0]) == 0 || responseT.Marketdata.Data[0][0] == nil {
		slog.Error("wrong answer format moex TOM")
		return 0.0, err
	}

	return responseT.Marketdata.Data[0][0].(float64), nil
}
