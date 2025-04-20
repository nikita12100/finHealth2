package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"test2/internal/common"
	"test2/internal/db"
	"test2/internal/fetcher"
	"test2/internal/models"
	"test2/internal/parser"

	"github.com/xuri/excelize/v2"
	tele "gopkg.in/telebot.v4"
)

const (
	listName = "broker_rep"
)

func fetchFromBuf(buf io.Reader) ([]models.Operation, []models.MoneyOperation) {
	f, err := excelize.OpenReader(buf)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows, err := f.GetRows(listName)
	if err != nil {
		log.Fatal(err)
	}

	operations, moneyOperations := parser.FetchData(rows)
	return operations, moneyOperations
}

func fetchFileData(c tele.Context) ([]models.Operation, []models.MoneyOperation, error) {
	file := c.Message().Document
	slog.Info("Got file", "name", file.FileName, "size(KB)", file.FileSize/1024)

	if file.FileName[len(file.FileName)-5:] != ".xlsx" {
		return nil, nil, fmt.Errorf("неправильный тип файла, поддерживаемый тип .xlsx")
	}

	reader, err := c.Bot().File(&file.File)
	if err != nil {
		slog.Warn("Failed to download file")
		return nil, nil, fmt.Errorf("ошибка скачивания: %v", err.Error())
	}
	defer reader.Close()

	operations, moneyOperations := fetchFromBuf(reader)
	return operations, moneyOperations, nil
}

func HandleBrockerReportFile(c tele.Context) error {
	newOperations, newMoneyOperations, err := fetchFileData(c)
	if err != nil {
		c.Send(err)
	}

	var resOperations []models.Operation
	var resMoneyOperations []models.MoneyOperation
	optOldPortfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			resOperations = newOperations
			resMoneyOperations = newMoneyOperations
		} else {
			slog.Error("Error getting portfolio", "error", err)
			return err
		}
	} else {
		resOperations = common.UnionOperation(optOldPortfolio.Operations, newOperations)
		resMoneyOperations = common.UnionOperation(optOldPortfolio.MoneyOperations, newMoneyOperations)
	}

	err = db.SavePortfolio(models.Portfolio{
		ChatId:          c.Chat().ID,
		Name:            "test",
		Operations:      resOperations,
		MoneyOperations: resMoneyOperations,
	})
	if err != nil {
		slog.Error("Failed save portfolio", "error", err)
	}
	slog.Info("Portfolio saved successfully")

	c.Send("Отчет успешно загружен из файла и сохранен")

	return c.Send(printDiffReport(newOperations), tele.ModeMarkdown)
}

func printDiffReport(operations []models.Operation) string {
	sumPrice := make(map[string]float64)
	countPerTicker := make(map[string]int)
	for _, operation := range operations {
		if operation.IsBuy {
			sumPrice[operation.Ticker] += operation.Price * float64(operation.Count)
			countPerTicker[operation.Ticker] += operation.Count
		} else {
			sumPrice[operation.Ticker] -= operation.Price * float64(operation.Count)
			countPerTicker[operation.Ticker] -= operation.Count
		}
	}

	sum := 0.0
	divSum := 0.0
	for ticker, sumPrice := range sumPrice {
		sum += sumPrice
		div, _ := fetcher.GetDivYieldCached(ticker) // прогноз на след 12мес.
		divSum += (div * float64(countPerTicker[ticker]) / 12.0)
	}

	report := "Отчет:\n"
	for ticker, count := range countPerTicker {
		report += fmt.Sprintf("*%v* +%v шт. на %.1f\n", ticker, count, sumPrice[ticker])
	}
	report += fmt.Sprintf("\nсуммарно куплено на %.0fр.\n", sum)
	report += fmt.Sprintf("\nпассивнй доход в месяц +%.2f\n", divSum)

	return report
}
