package parser

import (
	"encoding/json"
	"fmt"
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

func GetLastPrice(ticker string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://iss.moex.com/iss/engines/stock/markets/shares/securities/%s.json", ticker))
	if err != nil {
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
	for i, name := range responseT.Marketdata.Columns {
		if name == "PREVPRICE" {
			if len(responseT.Marketdata.Data) > 0 {
				lastPrice := responseT.Marketdata.Data[tqbrIndex][i].(float64)
				return lastPrice, nil
			} else {
				return 0.0, fmt.Errorf("empty moex response")
			}
		}
	}

	return 0.0, fmt.Errorf("ticker not found")
}
