package endpoints

import (
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

	page := echart.AddReplenishmentChart(statsPerMonth)
	page.Render(w)
}

func Grafik2Handler(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)

	statsShare := stats.GetLastStatShare(portfolio.Operations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})

	page := echart.AddSumPriceTotalChart(statsShare)
	page.Render(w)
}

func Grafik3Handler(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)
	statsPerMonth := stats.GetStatMoneyOperations(portfolio.MoneyOperations)

	page := echart.AddCoupAndDivChart(statsPerMonth)
	page.Render(w)
}

func Grafik4Handler(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)
	statsDivPerTicker := stats.GetStatMoneyOperationsSumDivPerTicker(portfolio.MoneyOperations)

	page := echart.AddSumDivTotalChart(statsDivPerTicker)
	page.Render(w)
}

func Grafik5Handler(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)

	page := echart.AddSumPriceTotalWithDivChart(portfolio) // todo send stat
	page.Render(w)
}

func Grafik6Handler(w http.ResponseWriter, r *http.Request) {
	portfolio := db.GetPortfolioOrCreate(507097513)
	statsShare := stats.GetLastStatShare(portfolio.Operations)
	statsShare = common.FilterValue(statsShare, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})

	page := echart.AddSumDivFutureChart(statsShare)
	page.Render(w)
}
