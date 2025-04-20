package handlers

import (
	"bytes"
	"log"
	"test2/internal/db"
	"test2/internal/plotters"
	"test2/internal/stats"

	tele "gopkg.in/telebot.v4"
)

func HandleStatsPortfolioPlot(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		log.Fatal(err)
	}

	c.Send("Инвестиции в месяц...")
	replenishments := stats.GetReplenishmentPerMonth(portfolio.MoneyOperations)
	plotBuffer, err := plotters.CreatePlot(replenishments)
	if err != nil {
		return err
	}
	photo := &tele.Photo{File: tele.FromReader(bytes.NewReader(plotBuffer.Bytes()))}
	c.Send(photo, "Here's your photo!")

	return c.Send("end plot")
}