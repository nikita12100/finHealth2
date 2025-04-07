package handlers

import (
	// "database/sql"
	// "errors"
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
	brokerReportFile23    = "../2023.xlsx"
	brokerReportFile24    = "../2024.xlsx"
	brokerReportFile25    = "../2025.xlsx"
	brokerReportFile23_24 = "../broker_report.xlsx"
	listName              = "broker_rep"
)

func StatsPortfolio(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test_2")
	if err != nil {
		log.Fatal(err)
	}

	// Stats>
	count := parser.CalcCount(portfolio.Operations)
	countSorted := common.Sort(count)
	avgPrice := parser.CalcAvgPrice(portfolio.Operations)
	// <Stats
	inserter.InsertIntoSheet(count, avgPrice, countSorted)
	return c.Send("Данные загружены в таблицу")
}

func UpdatePortfolio(c tele.Context) error {
	var resPortfolio models.Portfolio

	fileOperations23 := fetchFromFile(brokerReportFile23)
	fileOperations24 := fetchFromFile(brokerReportFile24)
	fileOperations25 := fetchFromFile(brokerReportFile25)
	resPortfolio = models.Portfolio{
		ChatId:     c.Chat().ID,
		Name:       "test_2",
		Operations: common.UnionOperation(common.UnionOperation(fileOperations23, fileOperations24), fileOperations25),
	}
	// optPortfolio, err := db.GetPortfolio(c.Chat().ID, "test_2") // how handle error
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		// Portfolio not found - use default
	// 		resPortfolio = models.Portfolio{
	// 			ChatId:     c.Chat().ID,
	// 			Name:       "test_2",
	// 			Operations: fileOperations,
	// 		}
	// 	} else {
	// 		log.Printf("Error getting portfolio: %v", err)
	// 		return err
	// 	}
	// } else {
	// 	resPortfolio = models.Portfolio{
	// 		ChatId:     c.Chat().ID,
	// 		Name:       "test_2",
	// 		Operations: common.UnionOperation(optPortfolio.Operations, fileOperations),
	// 	}
	// }

	err := db.SavePortfolio(resPortfolio)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Portfolio saved successfully")

	return c.Send("Отчет успешно загружен из файла и сохранен")
}

func fetchFromFile(fileName string) []models.Operation {
	f, err := excelize.OpenFile(fileName)
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
