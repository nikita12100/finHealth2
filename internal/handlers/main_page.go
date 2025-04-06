package handlers

import (
	"log"
	"test2/internal/common"
	"test2/internal/db"
	"test2/internal/inserter"
	"test2/internal/models"
	"test2/internal/parser"

	"github.com/xuri/excelize/v2"
	tele "gopkg.in/telebot.v4"
)

const (
	brokerReportFile = "../broker_report.xlsx"
	listName         = "broker_rep"
)

func StatsPortfolio(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		log.Fatal(err)
	}

	// Stats>
	count := parser.CalcCount(portfolio.Operations)
	countSorted := common.Sort(parser.CalcCount(portfolio.Operations))
	avgPrice := parser.CalcAvgPrice(portfolio.Operations)
	// <Stats
	inserter.InsertIntoSheet(count, avgPrice, countSorted)
	return c.Send("Данные загружены в таблицу")
}

func UpdatePortfolio(c tele.Context) error {
	operations := fetchFromFile()
	portfolio := models.Portfolio{
		ChatId:     c.Chat().ID,
		Name:       "test",
		Operations: operations,
	}
	err := db.SavePortfolio(portfolio)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Portfolio saved successfully")

	return c.Send("Отчет успешно загружен из файла и сохранен")
}

func fetchFromFile() []models.Operation {
	f, err := excelize.OpenFile(brokerReportFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows, err := f.GetRows(listName)
	if err != nil {
		log.Fatal(err)
	}

	return parser.FetchOperations(rows)
}
