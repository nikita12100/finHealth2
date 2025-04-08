package parser

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

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

// fetch currency from FACEUNIT
func GetLastPriceBond(ticker string) (float64, float64, float64, float64, string, error) {
	if ticker[:2] == "SU" {
		return getLastPrice(fmt.Sprintf("https://iss.moex.com/iss/engines/stock/markets/bonds/securities/%s.json", ticker), "TQOB")
	} else {
		return getLastPrice(fmt.Sprintf("https://iss.moex.com/iss/engines/stock/markets/bonds/securities/%s.json", ticker), "TQCB")
	}
}

func GetLastPriceShare(ticker string) (float64, error) {
	_, lastPrice, _, _, _, err := getLastPrice(fmt.Sprintf("https://iss.moex.com/iss/engines/stock/markets/shares/securities/%s.json", ticker), "TQBR")
	return lastPrice, err
}

func getLastPrice(url string, sType string) (float64, float64, float64, float64, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Got error while GET url", err)
		return 0.0, 0.0, 0.0, 0.0, "", err
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
		slog.Error("Failed parse to json answer", err)
		return 0.0, 0.0, 0.0, 0.0, "", err
	}

	var tqbrIndex int
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
				slog.Error("Not found COUPONVALUE in response")
			}
		} else if name == "PREVPRICE" {
			if len(responseT.Marketdata.Data) > 0 {
				lastPrice = responseT.Marketdata.Data[tqbrIndex][i].(float64)
			} else {
				slog.Error("Not found PREVPRICE in response")
			}
		} else if name == "FACEVALUE" {
			if len(responseT.Marketdata.Data) > 0 {
				faceValue = responseT.Marketdata.Data[tqbrIndex][i].(float64)
			} else {
				slog.Error("Not found FACEVALUE in response")
			}
		} else if name == "COUPONPERIOD" {
			if len(responseT.Marketdata.Data) > 0 {
				couponPeriod = responseT.Marketdata.Data[tqbrIndex][i].(float64)
			} else {
				slog.Error("Not found COUPONPERIOD in response")
			}
		} else if name == "FACEUNIT" {
			if len(responseT.Marketdata.Data) > 0 {
				faceUnit = responseT.Marketdata.Data[tqbrIndex][i].(string)
			} else {
				slog.Error("Not found FACEUNIT in response")
			}
		}
	}
	return couponValue, lastPrice, faceValue, couponPeriod, faceUnit, nil
}
