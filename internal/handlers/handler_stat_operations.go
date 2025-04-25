package handlers

import (
	"bytes"
	"fmt"
	"test2/internal/common"
	"test2/internal/db"
	"test2/internal/models"
	"test2/internal/plotters"
	"test2/internal/stats"

	"gonum.org/v1/plot"
	tele "gopkg.in/telebot.v4"
)

func HandleStatsDivMain(c tele.Context) error {
	portfolio := db.GetPortfolioOrCreate(c.Chat().ID)

	statsPerMonth := stats.GetStatMoneyOperations(portfolio.MoneyOperations)

	photo, err := getPhoto("Пассивный доход", "руб.", 1000, statsPerMonth, plotters.AddBarChartCoupAndDiv)
	if err != nil {
		return err
	}
	c.Send("выплаченно суммарно по месяцам")
	c.Send(photo, "Here's your photo!")
	sumDiv := 0
	sumCoup := 0
	for _, s := range statsPerMonth {
		sumDiv += int(s.Dividends)
		sumCoup += int(s.Coupon)
	}

	c.Send(fmt.Sprintf("Сумма купонов:%v, див:%v. Всего:%v", sumCoup, sumDiv, sumCoup+sumDiv))

	menu := &tele.ReplyMarkup{}
	btnDivPerShare := menu.Data("выплаченно по акциям", "btnDivPerShare")
	btnDivPerShareCost := menu.Data("[todo]выплачено суммарно по акциям к стоимости акции", "btnDivPerShareCost")
	btnDivFuture := menu.Data("будущие дивиденты", "btnDivFuture")
	menu.Inline(
		menu.Row(btnDivPerShare),
		menu.Row(btnDivPerShareCost),
		menu.Row(btnDivFuture),
	)

	return c.Send("Еще графики:", menu)
}

func HandleStatsDivPerShare(c tele.Context) error {
	c.RespondText("готовим графики")
	c.Bot().Edit(c.Message(), "Дивиденты по акциям:", &tele.ReplyMarkup{})

	portfolio := db.GetPortfolioOrCreate(c.Chat().ID)

	statsDivPerTicker := stats.GetStatMoneyOperationsSumDivPerTicker(portfolio.MoneyOperations)

	photo, err := getPhoto("Див по акциям", "руб.", 500, statsDivPerTicker, plotters.AddHistogramSumDivTotal)
	if err != nil {
		return err
	}

	return c.Send(photo, "Here's your photo!")
}

func HandleStatsDivPerShareCost(c tele.Context) error {
	c.RespondText("готовим графики")
	c.Bot().Edit(c.Message(), "Самооккупаемость акций:", &tele.ReplyMarkup{})

	// окупаемость
	return c.Send("todo HandleStatsDivPerShareCost")
}

func HandleStatsDivFuture(c tele.Context) error {
	c.RespondText("готовим графики")
	c.Bot().Edit(c.Message(), "Будущие дивиднты:", &tele.ReplyMarkup{})

	portfolio := db.GetPortfolioOrCreate(c.Chat().ID)
	statsShare := stats.GetLastStatShare(portfolio.Operations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})

	photo, err := getPhoto("Ожидаемые див в след 12мес", "% к средней цене покупки", 1, statsShare, plotters.AddHistogramSumDivFuture)
	if err != nil {
		return err
	}
	c.Send(photo, "Here's your photo!")

	c.Send("Итого по дивидентам в след 12мес.")
	statsBond := stats.GetLastStatBond(portfolio.Operations)
	statsBond = common.FilterValue(statsBond, func(stat models.StatsBond) bool {
		return stat.Count != 0
	})
	divShareSum := 0.0
	for _, stat := range statsShare {
		divShareSum += stat.SumDiv
	}
	divBondSum := 0.0
	for _, stat := range statsBond {
		divBondSum += stat.Coup2025
	}

	report := fmt.Sprintf("Дивидентов: %.0f, в месяц: *%.0f*\n", divShareSum, (divShareSum/12)) +
		fmt.Sprintf("Купонов: %.0f, в месяц: *%.0f*\n", divBondSum, (divBondSum/12)) +
		fmt.Sprintf("Итого: %.0f, в месяц: *%.0f*\n", divBondSum+divShareSum, ((divBondSum+divShareSum)/12))
	c.Send(report, tele.ModeMarkdown)

	return nil
}

func HandleStatsReplenishmentMain(c tele.Context) error {
	portfolio := db.GetPortfolioOrCreate(c.Chat().ID)

	statsPerMonth := stats.GetStatMoneyOperations(portfolio.MoneyOperations)

	photo, err := getPhoto("Пополнения", "руб.", 50000, statsPerMonth, plotters.AddBarChart)
	if err != nil {
		return err
	}
	c.Send(photo, "Here's your photo!")

	return nil
}

func getPhoto[T any](
	title string,
	yLabel string,
	ticks int,
	data T,
	addData func(T, *plot.Plot) error,
) (*tele.Photo, error) {
	plot := plotters.InitPlot(title, yLabel, ticks)
	err := addData(data, plot)
	if err != nil {
		return nil, err
	}

	plotBuffer, err := plotters.RenderPlot(plot)
	if err != nil {
		return nil, err
	}

	return &tele.Photo{File: tele.FromReader(bytes.NewReader(plotBuffer.Bytes()))}, nil
}
