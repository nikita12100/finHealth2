package main

import (
	"log"
	"log/slog"
	"os"
	"test2/internal/db"
	"test2/internal/handlers"
	"time"

	_ "github.com/mattn/go-sqlite3"
	tele "gopkg.in/telebot.v4"
)

const (
	tokenName = "TG_TOKEN_FIN_HEALTH"
)

func initLogger() {
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize attribute display
			if a.Key == slog.TimeKey {
				return slog.Attr{} // Remove time for cleaner output
			}
			return a
		},
	})
	slog.SetDefault(slog.New(textHandler))

}

func main() {
	initLogger()
	slog.Info("Starting portfolio manager bot...")
	db.InitTables()

	pref := tele.Settings{
		Token:  os.Getenv(tokenName),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", handleStartMsg)
	b.Handle(&tele.Btn{Text: "📈 Статистика портфеля. Таблицы"}, handlers.HandleStatsPortfolioTable)
	b.Handle(&tele.Btn{Text: "📈 Статистика портфеля. Графики"}, handlers.HandleStatsPortfolioPlot)
	b.Handle(&tele.Btn{Text: "📝 [DEV]Записать данные"}, handlers.HandleUpdatePortfolio)
	b.Handle(tele.OnDocument, handlers.HandleBrockerReportFile)

	b.Start()
}

func handleStartMsg(c tele.Context) error {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

	btnPortfolioStatsTable := menu.Text("📈 Статистика портфеля. Таблицы")
	btnPortfolioStatsPlot := menu.Text("📈 Статистика портфеля. Графики")
	btnPortfolioUpdate := menu.Text("📝 [DEV]Записать данные")
	btnHelp := menu.Text("❓ Помощь")

	menu.Reply(
		menu.Row(btnPortfolioStatsTable),
		menu.Row(btnPortfolioStatsPlot),
		menu.Row(btnPortfolioUpdate),
		menu.Row(btnHelp),
	)

	return c.Send(
		"Привет, я умею ...\nВыберите действие:",
		menu,
	)
}
