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
	"time"

	"github.com/xuri/excelize/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	tele "gopkg.in/telebot.v4"
)

const (
	listName = "broker_rep"
)

/*
Сколько я вкладываю каждый месяц
*/
func HandleStatsPortfolio(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		log.Fatal(err)
	}

	topShareWeightTable, topShareDivTable, divTable, sumTotalTable := printStatReport(&portfolio)

	c.Send("Топ акций по весу")
	c.Send(topShareWeightTable, tele.ModeMarkdown)

	c.Send("Топ акций по дивидентам")
	c.Send(topShareDivTable, tele.ModeMarkdown)

	c.Send("Итого по дивидентам")
	c.Send(divTable, tele.ModeMarkdown)

	c.Send("Итого сумма")
	c.Send(sumTotalTable, tele.ModeMarkdown)

	c.Send("Инвестиции в месяц...")

	return c.Send("end")
}

func printStatReport(p *models.Portfolio) (string, string, string, string) {
	statsShare := stats.GetLastStatShare(p.Operations)
	statsBond := stats.GetLastStatBond(p.Operations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})
	statsBond = common.FilterValue(statsBond, func(stat models.StatsBond) bool {
		return stat.Count != 0
	})

	topShareWeightTable := prepareTopShareWeightTable(statsShare)
	topShareDivTable := prepareTopShareDivTable(statsShare)
	divTable := prepareDivTable(statsShare, statsBond)
	sumTotalTable := prepareSumTotalTable(statsShare, statsBond)

	return topShareWeightTable, topShareDivTable, divTable, sumTotalTable
}

func prepareTopShareWeightTable(statsShare map[string]models.StatsShare) string {
	topShareWeightHeaders := []string{"№", "Ticker", fmt.Sprintf("W=%.0f", models.WEIGHT_NORM), "Sum"}
	var topShareWeightRows [][]string
	for i, kv := range common.SortValue(statsShare, func(i, j models.StatsShare) bool {
		return i.Weight > j.Weight
	}) {
		row := []string{strconv.Itoa(i + 1), kv.Key, fmt.Sprintf("%.2f", kv.Value.Weight), fmt.Sprintf("%.0f", kv.Value.SumPriceTotal)}
		topShareWeightRows = append(topShareWeightRows, row)
	}
	return common.PrintTable4(topShareWeightHeaders, topShareWeightRows)
}

func prepareTopShareDivTable(statsShare map[string]models.StatsShare) string {
	topShareDivHeaders := []string{"№", "Ticker", "%", "sumDiv"}
	var topShareDivRows [][]string
	for i, kv := range common.SortValue(statsShare, func(i, j models.StatsShare) bool {
		return i.DivPerc > j.DivPerc
	}) {
		row := []string{strconv.Itoa(i + 1), kv.Key, fmt.Sprintf("%.2f", kv.Value.DivPerc), fmt.Sprintf("%.0f", kv.Value.SumDiv)}
		topShareDivRows = append(topShareDivRows, row)
	}
	return common.PrintTable4(topShareDivHeaders, topShareDivRows)
}

func prepareDivTable(statsShare map[string]models.StatsShare, statsBond map[string]models.StatsBond) string {
	divHeaders := []string{"type", "month", "year"}
	var divRows [][]string
	divShareSum := 0.0
	for _, stat := range statsShare {
		divShareSum += stat.SumDiv
	}
	row := []string{"share", fmt.Sprintf("%.0f", (divShareSum / 12)), fmt.Sprintf("%.0f", divShareSum)}
	divRows = append(divRows, row)
	divBondSum := 0.0
	for _, stat := range statsBond {
		divBondSum += stat.Coup2025
	}
	row = []string{"bond", fmt.Sprintf("%.0f", (divBondSum / 12)), fmt.Sprintf("%.0f", divBondSum)}
	divRows = append(divRows, row)
	row = []string{"total", fmt.Sprintf("%.0f", ((divBondSum + divShareSum) / 12)), fmt.Sprintf("%.0f", divBondSum+divShareSum)}
	divRows = append(divRows, row)
	return common.PrintTable3(divHeaders, divRows)
}

func prepareSumTotalTable(statsShare map[string]models.StatsShare, statsBond map[string]models.StatsBond) string {
	sumTotalHeaders := []string{"type", "value"}
	var sumTotalRows [][]string
	shareSum := 0.0
	for _, stat := range statsShare {
		shareSum += stat.SumPriceTotal
	}
	row := []string{"share", fmt.Sprintf("%.0f", shareSum)}
	sumTotalRows = append(sumTotalRows, row)
	bondSum := 0.0
	for _, stat := range statsBond {
		bondSum += stat.SumPriceTotal
	}
	row = []string{"bond", fmt.Sprintf("%.0f", bondSum)}
	sumTotalRows = append(sumTotalRows, row)

	row = []string{"Э_TOOOM", fmt.Sprintf("116000")}
	sumTotalRows = append(sumTotalRows, row)
	row = []string{"total", fmt.Sprintf("%.0f", shareSum+bondSum+116000+140000)}
	sumTotalRows = append(sumTotalRows, row)

	return common.PrintTable2(sumTotalHeaders, sumTotalRows)
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

	// var res []stats.TimeValue
	// curSum := 0.0
	// var curMonth time.Time
	// for i, op := range newMoneyOperations {
	// 	if op.OperationType == models.Replenishment {
	// 		if curMonth.IsZero() {
	// 			curMonth = op.Time
	// 		}
	// 		if op.Time.Month() == curMonth.Month() {
	// 			curSum += op.AmountIn
	// 		} else {
	// 			res = append(res, stats.TimeValue{Time: curMonth, Value: curSum})
	// 			curSum = 0
	// 			curSum += op.AmountIn
	// 			curMonth = op.Time
	// 		}

	// 	}
	// 	if i == (len(newMoneyOperations) - 1) {
	// 		res = append(res, stats.TimeValue{Time: curMonth, Value: curSum})
	// 	}
	// }
	// for _, kv := range res {
	// 	fmt.Printf("wefwefwe date:%v-%v=%.0f\n", kv.Time.Year(), kv.Time.Month(), kv.Value)
	// }

	// pts := make(plotter.XYs, len(res))
	// for i, d := range res {
	// 	pts[i].X = float64(d.Time.Unix()) // Convert time to float
	// 	pts[i].Y = d.Value
	// }

	// p := plot.New()
	// p.Title.Text = "Time Series"
	// p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02"}
	// p.Y.Label.Text = "Values"

	// line, _ := plotter.NewLine(pts)
	// p.Add(line)
	// p.Save(8*vg.Inch, 4*vg.Inch, "timeseries.png")

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
