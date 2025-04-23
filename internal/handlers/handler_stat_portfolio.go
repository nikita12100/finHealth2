package handlers

import (
	"fmt"
	"log"
	"test2/internal/common"
	"test2/internal/db"
	"test2/internal/models"
	"test2/internal/plotters"
	"test2/internal/stats"

	tele "gopkg.in/telebot.v4"
)

func HandleStatsPortfolioAllocations(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		log.Fatal(err)
	}

	statsShare := stats.GetLastStatShare(portfolio.Operations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})

	photo, err := getPhoto("Распределение активов", "руб.", 10000, statsShare, plotters.AddHistogramSumPriceTotal)
	if err != nil {
		return err
	}

	return c.Send(photo, "Here's your photo!")
}

func HandleStatsPortfolioTable(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		log.Fatal(err)
	}

	statsShare := stats.GetLastStatShare(portfolio.Operations)
	statsBond := stats.GetLastStatBond(portfolio.Operations)
	statsTOM := stats.GetLastStatTOM(portfolio.Operations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})
	statsBond = common.FilterValue(statsBond, func(stat models.StatsBond) bool {
		return stat.Count != 0
	})
	statsTOM = common.FilterValue(statsTOM, func(stat models.StatsTOM) bool {
		return stat.Count != 0
	})

	sumTotalTable := prepareSumTotalTable(statsShare, statsBond, statsTOM)

	c.Send("Итого баланс")
	c.Send(sumTotalTable, tele.ModeMarkdown)

	return nil
}

// FIX _TOM+UDMN
func prepareSumTotalTable(
	statsShare map[string]models.StatsShare,
	statsBond map[string]models.StatsBond,
	statsTOM map[string]models.StatsTOM,
) string {
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

	TOMSum := 0.0
	for ticker, stat := range statsTOM {
		if ticker != "CNYRUB_TOM" { // FIXME покупка облиг за юани не уменьшает кол-во бумаг
			TOMSum += stat.SumPriceTotal
		}
	}
	row = []string{"gold", fmt.Sprintf("%.0f", TOMSum)}
	sumTotalRows = append(sumTotalRows, row)
	row = []string{"total", fmt.Sprintf("%.0f", shareSum+bondSum+TOMSum+55000)} // +UDMN
	sumTotalRows = append(sumTotalRows, row)

	return common.PrintTable2(sumTotalHeaders, sumTotalRows)
}

func HandleInfoPortfolio(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		log.Fatal(err)
	}

	report := fmt.Sprintf("Первая операция: %s\n", portfolio.TimePeriod.From.Format("2006-01-02")) +
		fmt.Sprintf("Последняя операция: %s\n", portfolio.TimePeriod.To.Format("2006-01-02"))
	c.Send(report)

	return nil
}
