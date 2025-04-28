package endpoints

import (
	"fmt"
	"net/http"
	"test2/internal/common"
	"test2/internal/db"
	"test2/internal/models"
	"test2/internal/plotters/echart"
	"test2/internal/stats"
)

func Grafik1Handler(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)
	statsPerMonth := stats.GetStatMoneyOperations(portfolio.MoneyOperations)

	htmlToInjectBytes := echart.AddReplenishmentChart(statsPerMonth)

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Chart Example</title>
		<script src="https://cdn.jsdelivr.net/npm/echarts@5.3.3/dist/echarts.min.js"></script>
	</head>
	<body>
		<h1>Go-ECharts Example</h1>
		<div id="chart-container"></div>
		<script>
			// Initialize the chart with the provided option
			var chart = echarts.init(document.getElementById('chart-container'));
			var option = %s
		</script>
		<a href="/">Вернуться на главную</a>
	</body>
	</html>
	`, htmlToInjectBytes.String())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, html)
}

func Grafik2Handler(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)

	statsShare := stats.GetLastStatShare(portfolio.Operations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})

	htmlToInjectBytes := echart.AddSumPriceTotalChart(statsShare)

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Chart Example</title>
		<script src="https://cdn.jsdelivr.net/npm/echarts@5.3.3/dist/echarts.min.js"></script>
	</head>
	<body>
		<h1>Go-ECharts Example</h1>
		<div id="chart-container"></div>
		<script>
			// Initialize the chart with the provided option
			var chart = echarts.init(document.getElementById('chart-container'));
			var option = %s
		</script>
		<a href="/">Вернуться на главную</a>
	</body>
	</html>
	`, htmlToInjectBytes.String())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, html)
}

func Grafik3Handler(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)
	statsPerMonth := stats.GetStatMoneyOperations(portfolio.MoneyOperations)

	htmlToInjectBytes := echart.AddCoupAndDivChart(statsPerMonth)

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Chart Example</title>
		<script src="https://cdn.jsdelivr.net/npm/echarts@5.3.3/dist/echarts.min.js"></script>
	</head>
	<body>
		<h1>Go-ECharts Example</h1>
		<div id="chart-container"></div>
		<script>
			// Initialize the chart with the provided option
			var chart = echarts.init(document.getElementById('chart-container'));
			var option = %s
		</script>
		<a href="/">Вернуться на главную</a>
	</body>
	</html>
	`, htmlToInjectBytes.String())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, html)
}

func Grafik4Handler(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)
	statsDivPerTicker := stats.GetStatMoneyOperationsSumDivPerTicker(portfolio.MoneyOperations)

	htmlToInjectBytes := echart.AddSumDivTotalChart(statsDivPerTicker)

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Chart Example</title>
		<script src="https://cdn.jsdelivr.net/npm/echarts@5.3.3/dist/echarts.min.js"></script>
	</head>
	<body>
		<h1>Go-ECharts Example</h1>
		<div id="chart-container"></div>
		<script>
			// Initialize the chart with the provided option
			var chart = echarts.init(document.getElementById('chart-container'));
			var option = %s
		</script>
		<a href="/">Вернуться на главную</a>
	</body>
	</html>
	`, htmlToInjectBytes.String())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, html)
}

func Grafik5Handler(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)

	htmlToInjectBytes := echart.AddSumPriceTotalWithDivChart(portfolio) // todo send stat

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Chart Example</title>
		<script src="https://cdn.jsdelivr.net/npm/echarts@5.3.3/dist/echarts.min.js"></script>
	</head>
	<body>
		<h1>Go-ECharts Example</h1>
		<div id="chart-container"></div>
		<script>
			// Initialize the chart with the provided option
			var chart = echarts.init(document.getElementById('chart-container'));
			var option = %s
		</script>
		<a href="/">Вернуться на главную</a>
	</body>
	</html>
	`, htmlToInjectBytes.String())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, html)
}

func Grafik6Handler(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)
	statsShare := stats.GetLastStatShare(portfolio.Operations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})

	htmlToInjectBytes := echart.AddSumDivFutureChart(statsShare)

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Chart Example</title>
		<script src="https://cdn.jsdelivr.net/npm/echarts@5.3.3/dist/echarts.min.js"></script>
	</head>
	<body>
		<h1>Go-ECharts Example</h1>
		<div id="chart-container"></div>
		<script>
			// Initialize the chart with the provided option
			var chart = echarts.init(document.getElementById('chart-container'));
			var option = %s
		</script>
		<a href="/">Вернуться на главную</a>
	</body>
	</html>
	`, htmlToInjectBytes.String())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, html)
}
