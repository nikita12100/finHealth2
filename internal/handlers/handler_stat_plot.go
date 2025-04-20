package handlers

import (
	"bytes"
	"fmt"
	"log"
	"test2/internal/db"
	"test2/internal/models"
	"test2/internal/plotters"
	"test2/internal/stats"

	"gonum.org/v1/plot"
	tele "gopkg.in/telebot.v4"
)

func HandleStatsPortfolioPlot(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		log.Fatal(err)
	}
	statsPerMonth := stats.GetReplenishmentPerMonth(portfolio.MoneyOperations)

	photo, err := getPhoto(statsPerMonth, plotters.InitPlotCouponAndDiv, plotters.AddHistogramCoupAndDiv)
	if err != nil {
		return err
	}
	c.Send(photo, "Here's your photo!")
	sumDiv := 0
	sumCoup := 0
	for _, s := range statsPerMonth {
		sumDiv += int(s.Dividends)
		sumCoup += int(s.Coupon)
	}

	return c.Send(fmt.Sprintf("Сумма купонов:%v, див:%v. Всего:%v", sumCoup, sumDiv, sumCoup+sumDiv))
}

func HandleStatsPortfolioPlotReplenishment(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		log.Fatal(err)
	}
	statsPerMonth := stats.GetReplenishmentPerMonth(portfolio.MoneyOperations)

	photo, err := getPhoto(statsPerMonth, plotters.InitPlotReplenishment, plotters.AddHistogram)
	if err != nil {
		return err
	}
	c.Send(photo, "Here's your photo!")

	return c.Send("Replenishments")
}

func getPhoto(
	statsPerMonth []models.StatsMoneyOperationSnapshoot,
	initPlot func() *plot.Plot,
	addData func([]models.StatsMoneyOperationSnapshoot, *plot.Plot) error,
) (*tele.Photo, error) {
	plot := initPlot()
	err := addData(statsPerMonth, plot)
	if err != nil {
		return nil, err
	}

	plotBuffer, err := plotters.RenderPlot(plot)
	if err != nil {
		return nil, err
	}

	return &tele.Photo{File: tele.FromReader(bytes.NewReader(plotBuffer.Bytes()))}, nil
}
