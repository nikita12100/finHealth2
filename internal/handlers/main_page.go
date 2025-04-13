package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strconv"
	"test2/internal/common"
	"test2/internal/db"
	"test2/internal/fetcher"
	"test2/internal/inserter"
	"test2/internal/models"
	"test2/internal/parser"
	"test2/internal/stats"

	"github.com/xuri/excelize/v2"
	tele "gopkg.in/telebot.v4"
)

const (
	listName = "broker_rep"
)

/*
топ акции по весу, по %див
пассивнфй доход
сумма акций, сумма облиг, всего баланс
разбивка по валютам
Сколько я вкладываю каждый месяц
*/
func HandleStatsPortfolio(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		log.Fatal(err)
	}

	topShareWeightTable, topShareDivTable := printStatReport(&portfolio)
	
	c.Send("Топ акций по весу")
	c.Send(topShareWeightTable, tele.ModeMarkdown)
	
	c.Send("Топ акций по дивидентам")
	c.Send(topShareDivTable, tele.ModeMarkdown)

	return c.Send("end")
}

func printStatReport(p *models.Portfolio) (string, string) {
	statsShare := stats.GetLastStatShare(p.Operations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})

	topShareWeightHeaders := []string{"№", "Ticker", fmt.Sprintf("W=%.0f", models.WEIGHT_NORM), "Sum"}
	var topShareWeightRows [][]string
	for i, kv := range common.SortValue(statsShare, func(i, j models.StatsShare) bool {
		return i.Weight > j.Weight
	}) {
		row := []string{strconv.Itoa(i + 1), kv.Key, fmt.Sprintf("%.2f", kv.Value.Weight), fmt.Sprintf("%.0f", kv.Value.SumPriceTotal)}
		topShareWeightRows = append(topShareWeightRows, row)
	}
	topShareWeightTable := common.PrintTable4(topShareWeightHeaders, topShareWeightRows)

	topShareDivHeaders := []string{"№", "Ticker", "%", "sumDiv"}
	var topShareDivRows [][]string
	for i, kv := range common.SortValue(statsShare, func(i, j models.StatsShare) bool {
		return i.DivPerc > j.DivPerc
	}) {
		row := []string{strconv.Itoa(i + 1), kv.Key, fmt.Sprintf("%.2f", kv.Value.DivPerc), fmt.Sprintf("%.0f", kv.Value.SumDiv)}
		topShareDivRows = append(topShareDivRows, row)
	}
	topShareDivTable := common.PrintTable4(topShareDivHeaders, topShareDivRows)

	return topShareWeightTable, topShareDivTable
}

func HandleUpdatePortfolio(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		log.Fatal(err)
	}

	statsShare := stats.GetLastStatShare(portfolio.Operations)
	statsBond := stats.GetLastStatBond(portfolio.Operations)

	inserter.InsertIntoSheet(statsShare, statsBond)
	return c.Send("Данные загружены в таблицу")
}

func fetchFromBuf(buf io.Reader) []models.Operation {
	f, err := excelize.OpenReader(buf)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows, err := f.GetRows(listName)
	if err != nil {
		log.Fatal(err)
	}

	return parser.FetchOperations(rows)
}

func fetchFileData(c tele.Context) ([]models.Operation, error) {
	file := c.Message().Document
	slog.Info("Got file", "name", file.FileName, "size(KB)", file.FileSize/1024)

	if file.FileName[len(file.FileName)-5:] != ".xlsx" {
		return nil, fmt.Errorf("неправильный тип файла, поддерживаемый тип .xlsx")
	}

	reader, err := c.Bot().File(&file.File)
	if err != nil {
		slog.Warn("Failed to download file")
		return nil, fmt.Errorf("ошибка скачивания: %v", err.Error())
	}
	defer reader.Close()

	return fetchFromBuf(reader), nil
}

func HandleBrockerReportFile(c tele.Context) error {
	newOperations, err := fetchFileData(c)
	if err != nil {
		c.Send(err)
	}

	var resOperations []models.Operation
	optOldPortfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			resOperations = newOperations
		} else {
			slog.Error("Error getting portfolio", "error", err)
			return err
		}
	} else {
		resOperations = common.UnionOperation(optOldPortfolio.Operations, newOperations)
	}

	err = db.SavePortfolio(models.Portfolio{
		ChatId:     c.Chat().ID,
		Name:       "test",
		Operations: resOperations,
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
