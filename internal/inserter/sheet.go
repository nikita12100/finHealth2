package inserter

import (
	"context"
	"log"
	"log/slog"
	"test2/internal/common"
	"test2/internal/models"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadsheetID   = "1lsmGWP02TFzMesZQz99S5EPRIpSJvadyCl646fOxqbQ"
	writeRangeShare = "script_share!A2:I"
	writeRangeBond  = "script_bond!A2:I"
	credentialsFile = "../gcred.json"
)

func InsertIntoSheet(statsShare map[string]models.StatsShare, statsBond map[string]models.StatsBond) {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	var valuesBonds [][]interface{}
	var valuesShares [][]interface{}
	countSortedShare := common.SortKey(statsShare)
	countSortedBond := common.SortKey(statsBond)
	for _, ticker := range countSortedShare {
		if statsShare[ticker].Count == 0 {
			continue
		}
		if models.IsCurrency(ticker) {
			continue
		}
		if models.IsShare(ticker) {
			valuesShares = append(valuesShares, []interface{}{
				ticker,
				statsShare[ticker].Weight,
				statsShare[ticker].Count,
				statsShare[ticker].AvgPriceBuy,
				statsShare[ticker].LastPrice,
				statsShare[ticker].Div,
				statsShare[ticker].SumPriceTotal,
				statsShare[ticker].SumDiv,
				statsShare[ticker].DivPerc,
			})
		}
	}
	for _, ticker := range countSortedBond {
		if statsBond[ticker].Count == 0 {
			continue
		}
		if models.IsCurrency(ticker) {
			continue
		}
		if models.IsBond(ticker) {
			if statsBond[ticker].CouponValue == 0.0 && statsBond[ticker].LastPrice == 0.0 && statsBond[ticker].AvgPriceBuy == 0.0 {
				slog.Warn("Skip bond in case of potintial expire", "bond", ticker)
				continue
			}

			valuesBonds = append(valuesBonds, []interface{}{
				ticker,
				statsBond[ticker].Count,
				statsBond[ticker].AvgPriceBuy,
				statsBond[ticker].LastPrice,
				statsBond[ticker].CouponValue,
				"nkd",
				statsBond[ticker].CouponPeriodPerYear,
				statsBond[ticker].SumPriceTotal,
				statsBond[ticker].Coup2025,
			})
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
