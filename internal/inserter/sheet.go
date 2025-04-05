package inserter

import (
	"context"
	"fmt"
	"log"
	"test2/internal/parser"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadsheetID   = "1lsmGWP02TFzMesZQz99S5EPRIpSJvadyCl646fOxqbQ"
	writeRange      = "script_info!A2:H"
	credentialsFile = "../gcred.json"
)

func InsertIntoSheet(data map[string]int, avgPrice map[string]float64, countSorted []string) {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	var values [][]interface{}
	for _, key := range countSorted {
		ticker := key
		avgPriceBuy := fmt.Sprintf("%.2f", avgPrice[ticker])
		lastPrice, _ := parser.GetLastPrice(ticker)
		count := data[ticker]
		values = append(values, []interface{}{key, "w", count, avgPriceBuy, lastPrice, "???", lastPrice*float64(count), "sumDiv"})
	}

	valueRange := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         values,
	}

	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, writeRange, valueRange).
		ValueInputOption("RAW").
		Do()
	if err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}

	log.Println("Data written successfully!")
}
