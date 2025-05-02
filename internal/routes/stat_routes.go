package routes

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"test2/internal/common"
	"test2/internal/db"
	"test2/internal/models"
	"test2/internal/plotters"
	"test2/internal/stats"
)

func insertDataChart(tmp *template.Template, dataRaw string) *bytes.Buffer {
	dataHtml := template.HTML(dataRaw)

	data := struct{ Data template.HTML }{Data: dataHtml}

	var buf bytes.Buffer
	err := tmp.Execute(&buf, data)
	if err != nil {
		slog.Error("Error inserting chart data into HTML template", "error", err)
		return nil
	}
	return &buf
}

func HandleStatsReplenishment(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handle /stat/replenishment")

	portfolio := db.GetPortfolioOrCreate(507097513)
	statsPerMonth := stats.GetStatMoneyOperations(portfolio.MoneyOperations)

	chart := plotters.AddReplenishmentChart(statsPerMonth)

	tmplPage := getTemplate("./static/stat_page.html")
	chartSnippet := chart.RenderSnippet()
	page := insertDataChart(tmplPage, chartSnippet.Element+chartSnippet.Script)

	w.Write(page.Bytes())
}

func HandleStatsAllocations(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handle /stat/allocations")

	portfolio := db.GetPortfolioOrCreate(507097513)

	statsShare := stats.GetLastStatShare(portfolio.Operations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})

	chart := plotters.AddAllocationsChart(statsShare)

	tmplPage := getTemplate("./static/stat_page.html")
	chartSnippet := chart.RenderSnippet()
	page := insertDataChart(tmplPage, chartSnippet.Element+chartSnippet.Script)

	w.Write(page.Bytes())
}

func HandleStatsDiv(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handle /stat/div")

	portfolio := db.GetPortfolioOrCreate(507097513)

	// ************************* render1 *************************
	statsPerMonth := stats.GetStatMoneyOperations(portfolio.MoneyOperations)
	chartDiv := plotters.AddDivChart(statsPerMonth)
	chartDivSnippet := chartDiv.RenderSnippet()

	sumDiv := 0
	sumCoup := 0
	for _, s := range statsPerMonth {
		sumDiv += int(s.Dividends)
		sumCoup += int(s.Coupon)
	}
	report := fmt.Sprintf("Сумма купонов:%v, див:%v. Всего:%v", sumCoup, sumDiv, sumCoup+sumDiv)
	render1 := chartDivSnippet.Element + chartDivSnippet.Script + report

	// ************************* render2 *************************
	statsDivPerTicker := stats.GetStatSumDivPerShare(portfolio.MoneyOperations)
	chartDivPerShare := plotters.AddSumDivPerShareChart(statsDivPerTicker)
	chartDivPerSSnippet := chartDivPerShare.RenderSnippet()
	render2 := chartDivPerSSnippet.Element + chartDivPerSSnippet.Script

	// ************************* render3 *************************
	chartPriceToDiv := plotters.AddPriceToDivChart(portfolio)
	chartPriceToDivSnippet := chartPriceToDiv.RenderSnippet()
	render3 := chartPriceToDivSnippet.Element + chartPriceToDivSnippet.Script

	// ************************* render4 *************************
	statsShare := stats.GetLastStatShare(portfolio.Operations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})

	statsBond := stats.GetLastStatBond(portfolio.Operations)
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

	report = "Итого по дивидентам в след 12мес.</br>"
	report += fmt.Sprintf("Дивидентов: %.0f, в месяц: <b>%.0f</b></br>", divShareSum, (divShareSum/12)) +
		fmt.Sprintf("Купонов: %.0f, в месяц: <b>%.0f</b></br>", divBondSum, (divBondSum/12)) +
		fmt.Sprintf("Итого: %.0f, в месяц: <b>%.0f</b></br>", divBondSum+divShareSum, ((divBondSum+divShareSum)/12))

	chartDivFuture := plotters.AddDivFutureChart(statsShare)
	chartDivFutureSnippet := chartDivFuture.RenderSnippet()
	render4 := chartDivFutureSnippet.Element + chartDivFutureSnippet.Script + report

	// ************************* render all *************************
	tmplPage := getTemplate("./static/stat_page.html")
	page := insertDataChart(tmplPage, render1+render2+render3+render4)

	w.Write(page.Bytes())
}
