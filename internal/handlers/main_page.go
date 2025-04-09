package handlers

import (
	// "database/sql"
	// "errors"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
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

func HandleStatsPortfolio(c tele.Context) error {
	portfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		log.Fatal(err)
	}

	count := portfolio.GetCountPerTicker()
	avgPrice := portfolio.GetAvgBuyPricePerTicker()

	inserter.InsertIntoSheet(count, avgPrice)
	return c.Send("Данные загружены в таблицу")
}

func HandleUpdatePortfolio(c tele.Context) error {
	return c.Send("Загрузите отчет формата .xlsx ...")
}

func fetchFromBuf(buf io.Reader) []models.Operation {
	f, err := excelize.OpenReader(buf)
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

func HandleBrockerReportFile(c tele.Context) error {
	file := c.Message().Document
	slog.Info("Got file", "name", file.FileName, "size(KB)", file.FileSize/1024)

	if file.FileName[len(file.FileName)-5:] != ".xlsx" {
		return c.Send("Неправильный тип файла, поддерживаемый тип .xlsx")
	}

	reader, err := c.Bot().File(&file.File)
	if err != nil {
		slog.Warn("Failed to download file")
		return c.Send("Failed to download file: " + err.Error())
	}
	defer reader.Close()
	// FILE DONE 
	operations := fetchFromBuf(reader)
	// PARSE DONE
	var resPortfolio models.Portfolio

	optPortfolio, err := db.GetPortfolio(c.Chat().ID, "test")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Portfolio not found - use default
			resPortfolio = models.Portfolio{
				ChatId:     c.Chat().ID,
				Name:       "test",
				Operations: operations,
			}
		} else {
			slog.Error("Error getting portfolio", err)
			return err
		}
	} else {
		resPortfolio = models.Portfolio{
			ChatId:     c.Chat().ID,
			Name:       "test",
			Operations: common.UnionOperation(optPortfolio.Operations, operations),
		}
	}

	err = db.SavePortfolio(resPortfolio)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Portfolio saved successfully")

	c.Send("Отчет успешно загружен из файла и сохранен")

	return c.Send(fmt.Sprintf("res_len=%v", len(operations)))
}
