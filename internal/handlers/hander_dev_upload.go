package handlers

import (
	"test2/internal/db"
	"test2/internal/inserter"
	"test2/internal/stats"

	tele "gopkg.in/telebot.v4"
)

func HandleUpdatePortfolio(c tele.Context) error {
	portfolio := db.GetPortfolioOrCreate(c.Chat().ID)

	statsShare := stats.GetLastStatShare(portfolio.Operations)
	statsBond := stats.GetLastStatBond(portfolio.Operations)

	inserter.InsertIntoSheet(statsShare, statsBond)
	return c.Send("Данные загружены в таблицу")
}
