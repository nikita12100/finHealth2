package parser

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"test2/internal/fetcher"
	. "test2/internal/models"
	"time"

	"log/slog"
)

type BlockPos struct {
	posOperationsBegin int
	posOperationsEnd   int
}

const (
	OPERATION_BEGIN = "Информация о совершенных и исполненных сделках на конец отчетного периода"
	OPERATION_END   = "Информация о неисполненных сделках на конец отчетного периода"
)

func calcPos(rows [][]string) BlockPos {
	var positions BlockPos
	for i, row := range rows {
		if len(row) > 0 {
			if strings.Contains(row[0], OPERATION_BEGIN) {
				positions.posOperationsBegin = i
			}
			if strings.Contains(row[0], OPERATION_END) {
				positions.posOperationsEnd = i
			}
		}
	}
	return positions
}

func renameTicker(someTicker string) string {
	renames := map[string]string{
		"US87238U2033": "T",
		"US83418T1088": "CIAN",
		"YNDX":         "YDEX",
		"FIVE":         "X5",
	}

	if rename, exists := renames[someTicker]; exists {
		return rename
	}
	prefixPattern := `^[A-Z]{2}[0-9]{2}`
	re := regexp.MustCompile(prefixPattern) // todo check if it is share
	if re.MatchString(someTicker) {
		if rename, err := fetcher.GetTickerByISINCached(someTicker); err == nil {
			return rename
		} else {
			slog.Error("Got error from moex, ticker=%v, error:", someTicker, err)
		}
	}
	return someTicker
}

func FetchOperations(rows [][]string) []Operation {
	positions := calcPos(rows)

	var operations []Operation
	for i := positions.posOperationsBegin + 2; i < positions.posOperationsEnd; i++ {
		if rows[i][0] == "" {
			continue
		}
		var operation Operation
		operation.NoTransaction = rows[i][0]
		operation.IsBuy = strings.Contains(rows[i][6], "Покупка")

		if !strings.Contains(rows[i][10], "%") {
			operation.Ticker = renameTicker(rows[i][8])
		} else {
			operation.Ticker = rows[i][8]
		}

		operation.Currency = rows[i][10]
		operation.Date, _ = parseDateTime(rows[i][3] + "T" + rows[i][4])
		operation.Price = parsePrice(rows[i][9], operation.Ticker, operation.Date)
		operation.Count = parseCount(rows[i][11], operation.Ticker, operation.Date)

		if operation.Ticker != "" {
			operations = append(operations, operation)
		}
	}

	return operations
}


func parsePrice(priceS string, ticker string, date time.Time) float64 { // fix не работает почему-то
	price, _ := strconv.ParseFloat(priceS, 64)
	gmknFragmentation, _ := time.Parse("02.01.2006", "04.04.2024")
	if ticker == "GMKN" && date.Before(gmknFragmentation) {
		price /= 100
	}
	plzlFragmentation, _ := time.Parse("02.01.2006", "25.03.2025")
	if ticker == "PLZL" && date.Before(plzlFragmentation) {
		price /= 10
	}
	return price
}

func parseCount(countS string, ticker string, date time.Time) int { // fix не работает почему-то
	count, _ := strconv.Atoi(countS)
	gmknFragmentation, _ := time.Parse("02.01.2006", "04.04.2024")
	if ticker == "GMKN" && date.Before(gmknFragmentation) {
		count *= 100
	}
	plzlFragmentation, _ := time.Parse("02.01.2006", "25.03.2025")
	if ticker == "PLZL" && date.Before(plzlFragmentation) {
		count *= 10
	}
	return count
}

func parseDateTime(timeStr string) (time.Time, error) {
	layout := "02.01.2006T15:04:05"

	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		log.Printf("Error parsing time=%v error:%v\n", timeStr, err)
		return time.Now(), err
	}

	return parsedTime, nil
}
