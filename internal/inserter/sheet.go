package inserter

import (
	"context"
	"fmt"
	"log"
	"test2/internal/common"
	"test2/internal/models"
	"test2/internal/parser"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadsheetID   = "1lsmGWP02TFzMesZQz99S5EPRIpSJvadyCl646fOxqbQ"
	writeRangeShare = "script_share!A2:I"
	writeRangeBond  = "script_bond!A2:I"
	credentialsFile = "../gcred.json"
	WEIGHT_NORM     = 60000.0
)

func InsertIntoSheet(count map[string]int, avgPrice map[string]float64) {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	var valuesBonds [][]interface{}
	var valuesShares [][]interface{}
	countSorted := common.Sort(count)
	for _, key := range countSorted {
		if count[key] == 0 {
			continue
		}
		if len(key)> 5 && key[len(key)-4:] == "_TOM" {
			continue
		}
		if models.IsBond(key) {
			ticker := key
			count := count[ticker]
			avgPriceBuy := avgPrice[ticker] / 100
			couponValue, lastPrice, faceValue, couponPeriod, faceUnit, _ := parser.GetLastPriceBond(ticker)
			if faceUnit == "USD" {
				lastPrice = lastPrice * 100
				avgPriceBuy = avgPriceBuy * 100
			}
			var couponPeriodPerYear int
			if couponPeriod != 0.0 {
				couponPeriodPerYear = int(365.0 / couponPeriod)
			} else {
				couponPeriodPerYear = 0
			}
			lastPrice = lastPrice / 100
			coup2025 := float64(count) * couponValue * float64(couponPeriodPerYear)

			valuesBonds = append(valuesBonds, []interface{}{ticker, count, avgPriceBuy * faceValue, lastPrice * faceValue, couponValue, "nkd", couponPeriodPerYear, float64(count) * lastPrice * faceValue, coup2025})
		} else if models.IsShare(key) {
			ticker := key
			lastPrice, _ := parser.GetLastPriceShare(ticker)
			count := count[ticker]
			currSum := lastPrice * float64(count)
			weight := fmt.Sprintf("%.2f", currSum/WEIGHT_NORM)
			avgPriceBuy := fmt.Sprintf("%.2f", avgPrice[ticker])
			div, _ := parser.GetDivYield(ticker)
			sumDiv := div*float64(count)
			divPerc := (div / avgPrice[ticker]) * 100

			valuesShares = append(valuesShares, []interface{}{ticker, weight, count, avgPriceBuy, lastPrice, div, currSum, sumDiv, divPerc})
		}
	}

	valueRangeShare := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         valuesShares,
	}
	valueRangeBonds := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         valuesBonds,
	}

	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, writeRangeShare, valueRangeShare).
		ValueInputOption("RAW").
		Do()
	if err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, writeRangeBond, valueRangeBonds).
		ValueInputOption("RAW").
		Do()
	if err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}

	log.Println("Data written successfully!")
}
