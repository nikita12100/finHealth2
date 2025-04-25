package parser

import (
	"fmt"
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
	posOperationsBegin      int
	posOperationsEnd        int
	posMoneyOperationsBegin int
	posMoneyOperationsEnd   int
}

const (
	OPERATION_BEGIN       = "Информация о совершенных и исполненных сделках на конец отчетного периода"
	OPERATION_END         = "Информация о неисполненных сделках на конец отчетного периода"
	MONEY_OPERATION_BEGIN = "Операции с денежными средствами"
	MONEY_OPERATION_END   = "Движение по ценным бумагам инвестора"
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
			if strings.Contains(row[0], MONEY_OPERATION_BEGIN) {
				positions.posMoneyOperationsBegin = i
			}
			if strings.Contains(row[0], MONEY_OPERATION_END) {
				positions.posMoneyOperationsEnd = i
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

func FetchData(rows [][]string) ([]Operation, []MoneyOperation) {
	positions := calcPos(rows)

	operations := fetchOperations(rows, positions)
	moneyOperations := fetchMoneyOperations(rows, positions)

	return operations, moneyOperations
}

func fetchOperations(rows [][]string, positions BlockPos) []Operation {
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

func parsePrice(priceS string, ticker string, date time.Time) float64 {
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

func parseCount(countS string, ticker string, date time.Time) int {
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

func parseDate(timeStr string) (time.Time, error) {
	layout := "02.01.2006"

	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		log.Printf("Error parsing time=%v error:%v\n", timeStr, err)
		return time.Now(), err
	}

	return parsedTime, nil
}

func fetchMoneyOperations(rows [][]string, positions BlockPos) []MoneyOperation {
	var moneyOperations []MoneyOperation
	skipStart := true
	for i := positions.posMoneyOperationsBegin + 2; i < positions.posMoneyOperationsEnd; i++ {
		if skipStart && len(rows[i]) > 1 {
			continue
		} else {
			skipStart = false
		}
		if len(rows[i]) <= 1 || (rows[i][1] == "" && rows[i][9] == "") || rows[i][9] == "Дата исполнения" {
			continue
		}
		var moneyOperation MoneyOperation
		moneyOperation.Time, _ = parseDate(rows[i][9])
		moneyOperation.OperationType = parseMoneyOperationType(rows[i][14])
		if moneyOperation.OperationType == Unknown {
			slog.Warn("Unknown OperationType for", "row", rows[i])
		}
		moneyOperation.AmountIn, _ = strconv.ParseFloat(rows[i][19], 64)
		moneyOperation.AmountOut, _ = strconv.ParseFloat(rows[i][24], 64)
		if len(rows[i]) == 28 {
			moneyOperation.CommentRaw = rows[i][27]

			if moneyOperation.OperationType == Dividends {
				moneyOperation.Comment = parseMoneyOperationComment(moneyOperation.CommentRaw)
			}
		}

		moneyOperations = append(moneyOperations, moneyOperation)
	}

	return moneyOperations
}

func parseMoneyOperationType(s string) OperationType {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "пополнение счета":
		return Replenishment
	case "покупка/продажа":
		return BuyOrSell
	case "комиссия за сделки":
		return Commision
	case "dvp/rvp":
		return DVP
	case "dfp/rfp":
		return DFP
	case "выплата купонов":
		return Coupon
	case "погашение облигации":
		return BondExpire
	case "налог":
		return Tax
	case "налог (по итогу года)":
		return TaxYear
	case "налог (дивиденды)":
		return TaxDiv
	case "выплата дивидендов":
		return Dividends
	case "вывод средств":
		return Withdraw
	case "репо":
		return Repo
	case "частичное погашение облигации (амортизация номинала)":
		return Redemption
	default:
		return Unknown
	}
}

func parseMoneyOperationComment(row string) CommentMoneyOperation {
	isin, count, err := parseMoneyOperationCommentRow(row)
	if err != nil {
		return CommentMoneyOperation{}
	}
	ticker, err := fetcher.GetTickerByISINCached(isin)
	if err != nil {
		return CommentMoneyOperation{}
	}

	return CommentMoneyOperation{Ticker: ticker, Count: count}
}

func parseMoneyOperationCommentRow(row string) (isin string, count int, err error) {
	re := regexp.MustCompile(`^([A-Z0-9]+)\/.*\/\s([0-9]+)\s*шт\.`)
	matches := re.FindStringSubmatch(row)
	if len(matches) < 3 {
		slog.Warn("Failed to find matches for MoneyOperationComment", "row", row)
		return "", 0, fmt.Errorf("invalid row format")
	}

	isin = strings.TrimSpace(matches[1])
	count, err = strconv.Atoi(matches[2])
	if err != nil {
		slog.Warn("Failed to find parse count for MoneyOperationComment", "row", row)
		return "", 0, fmt.Errorf("failed to parse count: %v", err)
	}

	return isin, count, nil
}
