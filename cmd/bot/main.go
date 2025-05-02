package main

import (
	"log/slog"
	"os"
	"test2/internal/db"
	"test2/internal/handlers"
	"time"

	_ "github.com/mattn/go-sqlite3"
	// _ "modernc.org/sqlite"
	tele "gopkg.in/telebot.v4"
)

const (
	botTokenEnvName  = "TG_TOKEN_FIN_HEALTH"
	serverURLEnvName = "BACKEND_URL_FIN_HEALTH"
	serverURLDefault = "https://sedbrkuebrhfeyrbg.serveo.net/"

	btnMiniAppText        = "Mini App"
	btnPortfolioInfo2Text = "💵 баланс"
	btnPortfolioInfo3Text = "ℹ️ информация о портфеле"
	btnHelpText           = "❓ Помощь"
)

func initLogger() {
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					return slog.String("time", t.Format("15:04:05"))
				}
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
		Token:  os.Getenv(botTokenEnvName),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		slog.Error("Error starting bot", "error", err)
		return
	}

	b.Handle("/start", handleStartMsg)

	b.Handle(&tele.Btn{Text: btnPortfolioInfo2Text}, handlers.HandleStatsPortfolioTable)
	b.Handle(&tele.Btn{Text: btnPortfolioInfo3Text}, handlers.HandleInfoPortfolio)

	b.Handle(tele.OnDocument, handlers.HandleBrockerReportFile)

	b.Start()
}

func handleStartMsg(c tele.Context) error {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

	serverURL := os.Getenv(serverURLEnvName)
	if serverURL == "" {
		serverURL = serverURLDefault
	}

	btnMiniApp := menu.WebApp("Mini App", &tele.WebApp{URL: serverURL})
	btnPortfolioInfo2 := menu.Text(btnPortfolioInfo2Text)
	btnPortfolioInfo3 := menu.Text(btnPortfolioInfo3Text)

	btnHelp := menu.Text(btnHelpText)

	menu.Reply(
		menu.Row(btnMiniApp),
		menu.Row(btnPortfolioInfo2, btnPortfolioInfo3),
		menu.Row(btnHelp),
	)

	return c.Send(
		"Привет, я умею ...\nВыберите действие:",
		menu,
	)
}
