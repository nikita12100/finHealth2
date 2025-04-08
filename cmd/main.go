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
	b.Handle(&tele.Btn{Text: "📈 Статистика портфеля"}, handlers.HandleStatsPortfolio)
	b.Handle(&tele.Btn{Text: "📝 Обновить портфель"}, handlers.HandleUpdatePortfolio)

	// b.Handle(&tele.Btn{Text: "😕 Confused"}, func(c tele.Context) error {
	// 	inlineMenu := &tele.ReplyMarkup{}
	// 	btnHelp := inlineMenu.Data("Get Help", "help_btn")
	// 	inlineMenu.Inline(inlineMenu.Row(btnHelp))

	// 	return c.Send(
	// 		"Let me help you! Click below:",
	// 		inlineMenu,
	// 	)
	// })

	b.Start()
}

func handleStartMsg(c tele.Context) error {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

	btnPortfolioStats := menu.Text("📈 Статистика портфеля")
	btnPortfolioUpdate := menu.Text("📝 Обновить портфель")
	btnHelp := menu.Text("❓ Помощь")

	menu.Reply(
		menu.Row(btnPortfolioStats),
		menu.Row(btnPortfolioUpdate),
		menu.Row(btnHelp),
	)

	return c.Send(
		"Привет, я умею ...\nВыберите действие:",
		menu,
	)
}
