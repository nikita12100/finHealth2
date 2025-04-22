package handlers

import (
	"bytes"
	"fmt"
	"log"
	"test2/internal/common"
	"test2/internal/db"
	"test2/internal/models"
	"test2/internal/plotters"
	"test2/internal/stats"

	"gonum.org/v1/plot"
	tele "gopkg.in/telebot.v4"
)

func HandleStatsDivMain(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		log.Fatal(err)
	}
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
	btnDivPerShareCost := menu.Data("выплачено суммарно по акциям к стоимости акции", "btnDivPerShareCost")
	btnDivFuture := menu.Data("будущие дивиденты", "btnDivFuture")
	menu.Inline(
		menu.Row(btnDivPerShare),
		menu.Row(btnDivPerShareCost),
		menu.Row(btnDivFuture),
	)

	return c.Send("Еще графики:", menu)
}

func HandleStatsDivPerShare(c tele.Context) error {
	return c.Send("todo HandleStatsDivPerShare")
}

func HandleStatsDivPerShareCost(c tele.Context) error {
	// окупаемость
	return c.Send("todo HandleStatsDivPerShareCost")
}

// (две таблицы сюда поедут, вторая таблица как сообщение внизу)
func HandleStatsDivFuture(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		log.Fatal(err)
	}
	statsShare := stats.GetLastStatShare(portfolio.Operations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})

	photo, err := getPhoto("Див в след 12мес", "% к телу", 1, statsShare, plotters.AddHistogramSumDivFuture)
	if err != nil {
		return err
	}
	c.Send(photo, "Here's your photo!")

	return c.Send("DivFuture")
}

func HandleStatsReplenishmentMain(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		log.Fatal(err)
	}
	statsPerMonth := stats.GetStatMoneyOperations(portfolio.MoneyOperations)

	photo, err := getPhoto("Пополнения", "руб.", 50000, statsPerMonth, plotters.AddBarChart)
	if err != nil {
		return err
	}
	c.Send(photo, "Here's your photo!")

	return c.Send("Replenishments")
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
