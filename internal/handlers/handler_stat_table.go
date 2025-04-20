package handlers

import (
	"fmt"
	"log"
	"strconv"
	"test2/internal/common"
	"test2/internal/db"
	"test2/internal/models"
	"test2/internal/stats"

	tele "gopkg.in/telebot.v4"
)

func HandleStatsPortfolioTable(c tele.Context) error {
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

	c.Send("[DEV] Итого сумма")
	c.Send(sumTotalTable, tele.ModeMarkdown)

	return c.Send("end table")
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
