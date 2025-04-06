package main

import (
	"log"
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

func main() {
	log.Println("Starting portfolio manager bot...")
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
	b.Handle(&tele.Btn{Text: "📈 Статистика портфеля"}, handlers.StatsPortfolio)
	b.Handle(&tele.Btn{Text: "📝 Обновить портфель"}, handlers.UpdatePortfolio)

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
