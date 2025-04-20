package handlers

import (
	"log"
	"test2/internal/db"
	"test2/internal/inserter"
	"test2/internal/stats"

	tele "gopkg.in/telebot.v4"
)

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
