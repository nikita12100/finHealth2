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
	tokenName             = "TG_TOKEN_FIN_HEALTH"
	btnPortfolioInfo1Text = "📊 распределение активов"
	btnPortfolioInfo2Text = "💵 баланс"
	btnPortfolioInfo3Text = "ℹ️ информация о портфеле"

	btnPortfolioStat1Text = "📈 пополнения"
	btnPortfolioStat2Text = "📈 дивиденты+купоны"
	btnPortfolioStat3Text = "[DEV] ???"

	btnHelpText = "❓ Помощь"
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

	b.Handle(&tele.Btn{Text: btnPortfolioInfo1Text}, handlers.HandleStatsPortfolioAllocations)
	b.Handle(&tele.Btn{Text: btnPortfolioInfo2Text}, handlers.HandleStatsPortfolioTable)
	b.Handle(&tele.Btn{Text: btnPortfolioInfo3Text}, handlers.HandleStatsPortfolioTable)

	b.Handle(&tele.Btn{Text: btnPortfolioStat1Text}, handlers.HandleStatsReplenishmentMain)
	b.Handle(&tele.Btn{Text: btnPortfolioStat2Text}, handlers.HandleStatsDivMain)
	b.Handle(&tele.Btn{Text: btnPortfolioStat3Text}, handlers.HandleUpdatePortfolio)

	b.Handle(&tele.Btn{Unique: "btnDivPerShare"}, handlers.HandleStatsDivPerShare)
	b.Handle(&tele.Btn{Unique: "btnDivPerShareCost"}, handlers.HandleStatsDivPerShareCost)
	b.Handle(&tele.Btn{Unique: "btnDivFuture"}, handlers.HandleStatsDivFuture)

	b.Handle(tele.OnDocument, handlers.HandleBrockerReportFile)

	b.Start()
}

func handleStartMsg(c tele.Context) error {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

	btnPortfolioInfo1 := menu.Text(btnPortfolioInfo1Text)
	btnPortfolioInfo2 := menu.Text(btnPortfolioInfo2Text)
	btnPortfolioInfo3 := menu.Text(btnPortfolioInfo3Text)

	btnPortfolioStat1 := menu.Text(btnPortfolioStat1Text)
	btnPortfolioStat2 := menu.Text(btnPortfolioStat2Text)
	btnPortfolioStat3 := menu.Text(btnPortfolioStat3Text)
	btnHelp := menu.Text(btnHelpText)

	menu.Reply(
		menu.Row(btnPortfolioInfo1, btnPortfolioStat1),
		menu.Row(btnPortfolioInfo2, btnPortfolioStat2),
		menu.Row(btnPortfolioInfo3, btnPortfolioStat3),
		menu.Row(btnHelp),
	)

	return c.Send(
		"Привет, я умею ...\nВыберите действие:",
		menu,
	)
}
