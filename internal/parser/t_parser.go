package parser

import (
	"regexp"
	"strconv"
	"strings"
	. "test2/internal/models"
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
	}

	if rename, exists := renames[someTicker]; exists {
		return rename
	}
	if rename, err := getTickerByISIN(someTicker); err == nil {
		return rename
	}
	return someTicker
}

// todo GMKN 100 -> 400
func FetchOperations(rows [][]string) []Operation {
	positions := calcPos(rows)

	var operations []Operation
	for i := positions.posOperationsBegin + 2; i < positions.posOperationsEnd; i++ {
		if len(rows[i]) < 10 {
			continue
		}
		var operation Operation
		operation.IsBuy = strings.Contains(rows[i][6], "Покупка")

		////
		prefixPattern := `^[A-Z]{2}[0-9]{2}`
		re := regexp.MustCompile(prefixPattern)
		if !strings.Contains(rows[i][10], "%") && re.MatchString(rows[i][8]) {
			operation.Ticker = renameTicker(rows[i][8])
		} else {
			operation.Ticker = rows[i][8]
		}
		////
		operation.Price, _ = strconv.ParseFloat(rows[i][9], 64)
		operation.Count, _ = strconv.Atoi(rows[i][11])
		operation.Currency = rows[i][10]

		if operation.Ticker != "" {
			operations = append(operations, operation)
		}
	}

	return operations
}

func CalcCount(operations []Operation) map[string]int {
	countMap := make(map[string]int)
	for _, operation := range operations {
		if operation.IsBuy {
			countMap[operation.Ticker] += operation.Count
		} else {
			countMap[operation.Ticker] -= operation.Count
		}
	}
	return countMap
}

func CalcAvgPrice(operations []Operation) map[string]float64 {
	sumPrice := make(map[string]float64)
	avgCount := make(map[string]int)
	for _, operation := range operations {
		if operation.IsBuy {
			sumPrice[operation.Ticker] += operation.Price * float64(operation.Count)
			avgCount[operation.Ticker] += operation.Count
		}
	}
	avgPrice := make(map[string]float64)
	for ticker, sumPrice := range sumPrice {
		avgPrice[ticker] = sumPrice / float64(avgCount[ticker])
	}
	return avgPrice
}

func removeFixPrefix(input string) string {
	if strings.HasPrefix(input, "FIX") {
		return strings.TrimPrefix(input, "FIX")
	}
	return input
}
