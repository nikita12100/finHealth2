package inserter

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"test2/internal/common"
	"test2/internal/fetcher"
	"test2/internal/models"

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
			bondInfo, _ := fetcher.GetLastPriceBondCached(ticker)
			if bondInfo.FaceUnit == "USD" {
				bondInfo.LastPrice = bondInfo.LastPrice * 100
				avgPriceBuy = avgPriceBuy * 100
			}
			var couponPeriodPerYear int
			if bondInfo.CouponPeriod != 0.0 {
				couponPeriodPerYear = int(365.0 / bondInfo.CouponPeriod)
			} else {
				couponPeriodPerYear = 0
			}
			bondInfo.LastPrice = bondInfo.LastPrice / 100
			coup2025 := float64(count) * bondInfo.CouponValue * float64(couponPeriodPerYear)
			if bondInfo.CouponValue == 0.0 && bondInfo.LastPrice == 0.0 && bondInfo.FaceValue == 0.0 {
				slog.Warn("Skip bond in case of potintial expire", "bond", ticker)
				continue
			}

			valuesBonds = append(valuesBonds, []interface{}{ticker, count, avgPriceBuy * bondInfo.FaceValue, bondInfo.LastPrice * bondInfo.FaceValue, bondInfo.CouponValue, "nkd", couponPeriodPerYear, float64(count) * bondInfo.LastPrice * bondInfo.FaceValue, coup2025})
		} else if models.IsShare(key) {
			ticker := key
			lastPrice, _ := fetcher.GetLastPriceShare(ticker)
			count := count[ticker]
			currSum := lastPrice * float64(count)
			weight := fmt.Sprintf("%.2f", currSum/WEIGHT_NORM)
			avgPriceBuy := fmt.Sprintf("%.2f", avgPrice[ticker])
			div, _ := fetcher.GetDivYieldCached(ticker)
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
