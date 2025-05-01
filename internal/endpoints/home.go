package endpoints

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
	"test2/internal/common"
	"test2/internal/db"
	"test2/internal/models"
	"test2/internal/plotters"
	"test2/internal/stats"

	"github.com/go-echarts/go-echarts/v2/charts"
)

func getTemplate() *template.Template {
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Biba name</title>
		<script src="https://go-echarts.github.io/go-echarts-assets/assets/echarts.min.js"></script>
	</head>
	<body>
		<header>
			<h1><a href="/" style="text-decoration: none; color: inherit;">🏠 Home</a></h1>
		</header>
		<main>
			{{.Element}} {{.Script}}
			<style>
				.container {margin-top:30px; display: flex;justify-content: center;align-items: center;}
				.item {margin: auto;}
			</style>
		</main>
	</body>
	</html>
	`
	t := template.New("snippet")
	t, err := t.Parse(tmpl)
	if err != nil {
		slog.Error("Error parsing html template", "html", tmpl, "error", err)
		return nil
	}

	return t
}

func insertData(tmp *template.Template, chart *charts.Bar) *bytes.Buffer {
	chartSnippet := chart.RenderSnippet()
	data := struct {
		Element template.HTML
		Script  template.HTML
	}{
		Element: template.HTML(chartSnippet.Element),
		Script:  template.HTML(chartSnippet.Script),
	}

	var buf bytes.Buffer
	err := tmp.Execute(&buf, data)
	if err != nil {
		slog.Error("Error inserting data into HTML template", "error", err)
		return nil
	}
	return &buf
}

func HandleStatsReplenishment(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)
	statsPerMonth := stats.GetStatMoneyOperations(portfolio.MoneyOperations)

	chart := plotters.AddReplenishmentChart(statsPerMonth)
	
	tmplPage := getTemplate()
	page := insertData(tmplPage, chart)

	w.Write(page.Bytes())
}

func HandleStatsAllocations(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)

	statsShare := stats.GetLastStatShare(portfolio.Operations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})

	chart := plotters.AddAllocationsChart(statsShare)

	tmplPage := getTemplate()
	page := insertData(tmplPage, chart)

	w.Write(page.Bytes())
}

func HandleStatsDiv(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)
	statsPerMonth := stats.GetStatMoneyOperations(portfolio.MoneyOperations)

	chart := plotters.AddDivChart(statsPerMonth)
	
	tmplPage := getTemplate()
	page := insertData(tmplPage, chart)

	w.Write(page.Bytes())

	// sumDiv := 0
	// sumCoup := 0
	// for _, s := range statsPerMonth {
	// 	sumDiv += int(s.Dividends)
	// 	sumCoup += int(s.Coupon)
	// }

	// c.Send(fmt.Sprintf("Сумма купонов:%v, див:%v. Всего:%v", sumCoup, sumDiv, sumCoup+sumDiv))
}

func HandleStatsDivPerShare(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)
	statsDivPerTicker := stats.GetStatSumDivPerShare(portfolio.MoneyOperations)

	chart := plotters.AddSumDivPerShareChart(statsDivPerTicker)
	
	tmplPage := getTemplate()
	page := insertData(tmplPage, chart)

	w.Write(page.Bytes())
}

func HandleStatsDivPerShareCost(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)

	chart := plotters.AddSumPriceTotalWithDivChart(portfolio) // todo send stat
	
	tmplPage := getTemplate()
	page := insertData(tmplPage, chart)

	w.Write(page.Bytes())
}

func HandleStatsDivFuture(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)
	statsShare := stats.GetLastStatShare(portfolio.Operations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})

	chart := plotters.AddDivFutureChart(statsShare)
	
	tmplPage := getTemplate()
	page := insertData(tmplPage, chart)

	w.Write(page.Bytes())

	// c.Send("Итого по дивидентам в след 12мес.")
	// statsBond := stats.GetLastStatBond(portfolio.Operations)
	// statsBond = common.FilterValue(statsBond, func(stat models.StatsBond) bool {
	// 	return stat.Count != 0
	// })
	// divShareSum := 0.0
	// for _, stat := range statsShare {
	// 	divShareSum += stat.SumDiv
	// }
	// divBondSum := 0.0
	// for _, stat := range statsBond {
	// 	divBondSum += stat.Coup2025
	// }

	// report := fmt.Sprintf("Дивидентов: %.0f, в месяц: *%.0f*\n", divShareSum, (divShareSum/12)) +
	// 	fmt.Sprintf("Купонов: %.0f, в месяц: *%.0f*\n", divBondSum, (divBondSum/12)) +
	// 	fmt.Sprintf("Итого: %.0f, в месяц: *%.0f*\n", divBondSum+divShareSum, ((divBondSum+divShareSum)/12))
	// c.Send(report, tele.ModeMarkdown)
}
