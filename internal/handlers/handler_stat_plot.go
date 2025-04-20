package handlers

import (
	"bytes"
	"log"
	"test2/internal/db"
	"test2/internal/models"
	"test2/internal/plotters"
	"test2/internal/stats"

	tele "gopkg.in/telebot.v4"
)

func HandleStatsPortfolioPlot(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		log.Fatal(err)
	}
	statsPerMonth := stats.GetReplenishmentPerMonth(portfolio.MoneyOperations)

	photo, err := getPhoto(statsPerMonth)
	if err != nil {
		return err
	}
	c.Send("Инвестиции в месяц...")
	c.Send(photo, "Here's your photo!")

	return c.Send("end plot")
}

func getPhoto(statsPerMonth []models.StatsMoneyOperationSnapshoot) (*tele.Photo, error){
	plot := plotters.InitPlot()
	err := plotters.AddHistogram2(statsPerMonth, plot)
	if err != nil {
		return nil, err
	}

	plotBuffer, err := plotters.RenderPlot(plot)
	if err != nil {
		return nil, err
	}

	return &tele.Photo{File: tele.FromReader(bytes.NewReader(plotBuffer.Bytes()))}, nil
}
