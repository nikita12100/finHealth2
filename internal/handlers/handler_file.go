package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"test2/internal/common"
	"test2/internal/db"
	"test2/internal/models"
	"test2/internal/parser"
	"test2/internal/stats"
	"time"

	"github.com/xuri/excelize/v2"
	tele "gopkg.in/telebot.v4"
)

const (
	listName = "broker_rep"
)

func fetchFromBuf(buf io.Reader) ([]models.Operation, []models.MoneyOperation) {
	f, err := excelize.OpenReader(buf)
	if err != nil {
		slog.Error("Error while opening buffer file", "error", err)
		return nil, nil
	}
	defer f.Close()

	rows, err := f.GetRows(listName)
	if err != nil {
		slog.Error("Error while fetching rows from file", "error", err)
		return nil, nil
	}

	operations, moneyOperations := parser.FetchData(rows)
	return operations, moneyOperations
}

func fetchFileData(c tele.Context) ([]models.Operation, []models.MoneyOperation, error) {
	file := c.Message().Document
	slog.Info("Got file", "name", file.FileName, "size(KB)", file.FileSize/1024)

	if file.FileName[len(file.FileName)-5:] != ".xlsx" {
		slog.Warn("Wrong file format", "name", file.FileName)
		return nil, nil, fmt.Errorf("неправильный тип файла, поддерживаемый тип .xlsx")
	}

	reader, err := c.Bot().File(&file.File)
	if err != nil {
		slog.Warn("Failed to download file")
		return nil, nil, fmt.Errorf("ошибка скачивания: %v", err.Error())
	}
	defer reader.Close()

	operations, moneyOperations := fetchFromBuf(reader)
	return operations, moneyOperations, nil
}

func HandleBrockerReportFile(c tele.Context) error {
	newOperations, newMoneyOperations, err := fetchFileData(c)
	if err != nil {
		slog.Warn("Error fetching file data", "chatId", c.Chat().ID, "error", err)
		return c.Send("Ошибка чтения файла")
	}

	optOldPortfolio := db.GetPortfolioOrCreate(c.Chat().ID)
	resOperations := common.UnionOperation(optOldPortfolio.Operations, newOperations)
	resMoneyOperations := common.UnionOperation(optOldPortfolio.MoneyOperations, newMoneyOperations)
	resTimePeriod := optOldPortfolio.TimePeriod.ExtendTimePeriod(models.DateRange{
		Start: newOperations[0].Date,
		End:   newOperations[len(newOperations)-1].Date,
	})

	portfolio := models.Portfolio{
		ChatId:          c.Chat().ID,
		Operations:      resOperations,
		MoneyOperations: resMoneyOperations,
		UpdatedAt:       time.Now(),
		TimePeriod:      resTimePeriod,
	}

	if err = db.SavePortfolio(portfolio); err != nil {
		slog.Error("Failed save portfolio", "chatId", c.Chat().ID, "error", err)
	}
	slog.Info("Portfolio saved successfully", "chatId", c.Chat().ID)

	report := printDiffReport(optOldPortfolio.Operations, newOperations)
	return c.Send(report, tele.ModeMarkdown)
}

func printDiffReport(oldOperations []models.Operation, newOperations []models.Operation) string {
	if common.SlicesContainsAll(oldOperations, newOperations) {
		return "Нечего обновлять, все операции уже в портфеле"
	}
	diffOperations := common.SlicesDifference(oldOperations, newOperations)

	statsShare := stats.GetLastStatShare(diffOperations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})
	statsBond := stats.GetLastStatBond(diffOperations)
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

	report := fmt.Sprintf("Дивидентов: +%.0f, в месяц: *+%.0f*\n", divShareSum, (divShareSum/12)) +
		fmt.Sprintf("Купонов: +%.0f, в месяц: *+%.0f*\n", divBondSum, (divBondSum/12)) +
		fmt.Sprintf("Итого: +%.0f, в месяц: *+%.0f*\n", divBondSum+divShareSum, ((divBondSum+divShareSum)/12))

	return report
}
