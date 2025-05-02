package routes

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

func insertDataChart(tmp *template.Template, chart *charts.Bar) *bytes.Buffer {
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
	page := insertDataChart(tmplPage, chart)

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
	page := insertDataChart(tmplPage, chart)

	w.Write(page.Bytes())
}

func HandleStatsDiv(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handle /stat/div")

	portfolio := db.GetPortfolioOrCreate(507097513)
	statsPerMonth := stats.GetStatMoneyOperations(portfolio.MoneyOperations)

	chart := plotters.AddDivChart(statsPerMonth)

	tmplPage := getTemplate("./static/stat_page.html")
	page := insertDataChart(tmplPage, chart)

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
	slog.Info("Handle /stat/div_per_share")

	portfolio := db.GetPortfolioOrCreate(507097513)
	statsDivPerTicker := stats.GetStatSumDivPerShare(portfolio.MoneyOperations)

	chart := plotters.AddSumDivPerShareChart(statsDivPerTicker)

	tmplPage := getTemplate("./static/stat_page.html")
	page := insertDataChart(tmplPage, chart)

	w.Write(page.Bytes())
}

func HandleStatsDivPerShareCost(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handle /stat/div_per_share_cost")
	
	portfolio := db.GetPortfolioOrCreate(507097513)

	chart := plotters.AddSumPriceTotalWithDivChart(portfolio) // todo send stat

	tmplPage := getTemplate("./static/stat_page.html")
	page := insertDataChart(tmplPage, chart)

	w.Write(page.Bytes())
}

func HandleStatsDivFuture(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handle /stat/div_future")
	
	portfolio := db.GetPortfolioOrCreate(507097513)
	statsShare := stats.GetLastStatShare(portfolio.Operations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})

	chart := plotters.AddDivFutureChart(statsShare)

	tmplPage := getTemplate("./static/stat_page.html")
	page := insertDataChart(tmplPage, chart)

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
