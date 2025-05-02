package main

import (
	"log/slog"
	"net/http"
	"os"
	"test2/internal/routes"
	"time"

	_ "github.com/mattn/go-sqlite3"
	// _ "modernc.org/sqlite"
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
	slog.Info("Starting portfolio server...")

	http.HandleFunc("/", routes.HomePageHandler)
	http.HandleFunc("/stat/replenishment", routes.HandleStatsReplenishment)
	http.HandleFunc("/stat/allocations", routes.HandleStatsAllocations)
	http.HandleFunc("/stat/div", routes.HandleStatsDiv)
	http.HandleFunc("/stat/div_per_share", routes.HandleStatsDivPerShare)
	http.HandleFunc("/stat/div_per_share_cost", routes.HandleStatsDivPerShareCost)
	http.HandleFunc("/stat/div_future", routes.HandleStatsDivFuture)

	slog.Info("Starting server on http://localhost:8085 ...")
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		slog.Error("Error ListenAndServe", "error", err)
	}
}
