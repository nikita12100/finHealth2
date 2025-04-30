package echart

import (
	"test2/internal/common"
	"test2/internal/models"
	"test2/internal/stats"

	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func AddReplenishmentChart(stats []models.StatsMoneyOperationSnapshoot) *components.Page {
	values := make([]opts.BarData, len(stats))
	labels := make([]string, len(stats))

	for i, stat := range stats {
		values[i] = opts.BarData{Value: stat.Replenishment}
		labels[i] = stat.Time.Format("06/01")
	}

	data := ColumnsData{
		Data:   values,
		Labels: labels,
		Legend: []string{"пополнения"},
	}

	options := ChartOptions{
		Title:      "Пополнения",
		Subtitle:   "пополнения брокерского счета",
		EnableZoom: true,
	}

	return getChart(data, options)
}

func AddCoupAndDivChart(stats []models.StatsMoneyOperationSnapshoot) *components.Page {
	values := make([]opts.BarData, len(stats))
	values2 := make([]opts.BarData, len(stats))
	labels := make([]string, len(stats))

	for i, stat := range stats {
		values[i] = opts.BarData{Value: stat.Coupon}
		values2[i] = opts.BarData{Value: stat.Dividends}
		labels[i] = stat.Time.Format("06/01")
	}

	data := ColumnsData{
		Data:   values,
		Data2:  values2,
		Legend: []string{"купоны", "дивиденты"},
		Labels: labels,
	}

	options := ChartOptions{
		Title:      "Пассивный доход",
		Subtitle:   "поступление дивидентов и купонов",
		EnableZoom: true,
	}

	return getChart(data, options)
}

func AddSumPriceTotalChart(stats map[string]models.StatsShare) *components.Page {
	values := make([]opts.BarData, len(stats))
	labels := make([]string, len(stats))

	statsKV := common.SortValue(stats, func(i, j models.StatsShare) bool {
		return i.SumPriceTotal > j.SumPriceTotal
	})

	i := 0
	for _, kv := range statsKV {
		values[i] = opts.BarData{Value: kv.Value.SumPriceTotal}
		labels[i] = kv.Key
		i++
	}

	data := ColumnsData{
		Data:   values,
		Labels: labels,
		Legend: []string{"текущая цена"},
	}

	options := ChartOptions{
		Title:    "Распределение акций",
		Subtitle: "текущие цены",
	}

	return getChart(data, options)
}

func AddSumDivTotalChart(stats map[string]float64) *components.Page {
	values := make([]opts.BarData, len(stats))
	labels := make([]string, len(stats))

	statsKV := common.SortValue(stats, func(i, j float64) bool {
		return i > j
	})

	i := 0
	for _, kv := range statsKV {
		values[i] = opts.BarData{Value: kv.Value}
		labels[i] = kv.Key
		i++
	}

	data := ColumnsData{
		Data:   values,
		Labels: labels,
		Legend: []string{"дивиденты"},
	}

	options := ChartOptions{
		Title:    "Выплачено двидентов",
		Subtitle: "",
	}

	return getChart(data, options)
}

func AddSumPriceTotalWithDivChart(p models.Portfolio) *components.Page {
	statsDivPerTicker := stats.GetStatMoneyOperationsSumDivPerTicker(p.MoneyOperations)
	stats := stats.GetLastStatShare(p.Operations)
	stats = common.FilterValue(stats, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})

	values := make([]opts.BarData, len(stats))
	values2 := make([]opts.BarData, len(stats))
	labels := make([]string, len(stats))

	statsKV := common.SortValue(stats, func(i, j models.StatsShare) bool {
		return i.SumPriceTotal > j.SumPriceTotal
	})

	i := 0
	for _, kv := range statsKV {
		values[i] = opts.BarData{Value: statsDivPerTicker[kv.Key]}
		values2[i] = opts.BarData{Value: kv.Value.SumPriceTotal - statsDivPerTicker[kv.Key]}
		labels[i] = kv.Key
		i++
	}

	data := ColumnsData{
		Data:   values,
		Data2:  values2,
		Legend: []string{"дивиденты", "остальное"},
		Labels: labels,
	}

	options := ChartOptions{
		Title:    "Самоокупаемость акций",
		Subtitle: "стоимость акции в портфеле к выплаченным по ней дивидентам",
	}

	return getChart(data, options)
}

func AddSumDivFutureChart(stats map[string]models.StatsShare) *components.Page {
	values := make([]opts.BarData, len(stats))
	labels := make([]string, len(stats))

	statsKV := common.SortValue(stats, func(i, j models.StatsShare) bool {
		return i.DivPerc > j.DivPerc
	})

	i := 0
	for _, kv := range statsKV {
		values[i] = opts.BarData{
			Value: kv.Value.DivPerc,
		}
		labels[i] = kv.Key
		i++
	}

	data := ColumnsData{
		Data:   values,
		Labels: labels,
		Legend: []string{"ожидаемые див"},
	}

	options := ChartOptions{
		Title:    "Будущие дивиденты",
		Subtitle: "ожидаемые дивиденты в след 12 месяцев",
	}

	return getChart(data, options)
}
