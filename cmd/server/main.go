package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"test2/internal/endpoints"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

	http.HandleFunc("/", mainPageHandler)
	http.HandleFunc("/grafik1", endpoints.Grafik1Handler)
	http.HandleFunc("/grafik2", endpoints.Grafik2Handler)
	http.HandleFunc("/grafik3", endpoints.Grafik3Handler)
	http.HandleFunc("/grafik4", endpoints.Grafik4Handler)
	http.HandleFunc("/grafik5", endpoints.Grafik5Handler)
	http.HandleFunc("/grafik6", endpoints.Grafik6Handler)

	log.Println("Starting server on http://localhost:8085 ...")
	log.Fatal(http.ListenAndServe(":8085", nil))
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Графики</title>
		<style>
			* {
				box-sizing: border-box;
				margin: 0;
				padding: 0;
			}
			body {
				font-family: 'Arial', sans-serif;
				background-color: #f4f7fc;
				color: #333;
				display: flex;
				justify-content: center;
				align-items: center;
				height: 100vh;
				padding: 20px;
			}
			.container {
				background-color: #ffffff;
				border-radius: 8px;
				box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
				padding: 40px;
				width: 100%;
				max-width: 600px;
				text-align: center;
			}
			h1 {
				font-size: 2rem;
				margin-bottom: 20px;
				color: #4CAF50;
			}
			ul {
				list-style: none;
				padding: 0;
			}
			ul li {
				margin-bottom: 20px;
			}
			ul li a {
				display: inline-block;
				padding: 15px 25px;
				background-color: #4CAF50;
				color: #fff;
				font-size: 1.2rem;
				text-decoration: none;
				border-radius: 6px;
				transition: background-color 0.3s ease;
			}
			ul li a:hover {
				background-color: #45a049;
			}
			@media (max-width: 600px) {
				h1 {
					font-size: 1.5rem;
				}
				ul li a {
					font-size: 1rem;
					padding: 12px 20px;
				}
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Выберите график</h1>
			<ul>
				<li><a href="/grafik1">пополнения</a></li>
				<li><a href="/grafik2">распределение активов</a></li>
				<li><a href="/grafik3">див+купоны</a></li>
				<li><a href="/grafik4">выплачено по акциям</a></li>
				<li><a href="/grafik5">самоокупаемость акций</a></li>
				<li><a href="/grafik6">будущие дивиденты</a></li>
			</ul>
		</div>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, html)
}
