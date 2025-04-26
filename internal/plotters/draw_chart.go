package plotters

import (
	"fmt"
	"image/color"
	"log/slog"
	"test2/internal/common"
	"test2/internal/models"
	"test2/internal/stats"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func AddReplenishmentChart(stats []models.StatsMoneyOperationSnapshoot, plot *plot.Plot) error {
	xys := make([]plotter.XY, len(stats))
	labels := make([]string, len(stats))
	columnLabels := make([]string, len(stats))
	columnLabelPos := make([]plotter.XY, len(stats))

	for i, stat := range stats {
		xys[i] = plotter.XY{
			X: float64(i),
			Y: stat.Replenishment,
		}
		labels[i] = stat.Time.Format("06/01")
		columnLabels[i] = fmt.Sprintf("%.0f", stat.Replenishment)
		columnLabelPos[i] = plotter.XY{
			X: float64(i) - 0.25,
			Y: stat.Replenishment + 50,
		}
	}

	data := ColumnsData{
		XYs:            xys,
		Labels:         labels,
		ColumnLabels:   columnLabels,
		ColumnLabelPos: columnLabelPos,
		BarWidth:       20,
	}

	if err := addBarChart(data, plot); err != nil {
		slog.Error("Failed to add chart for StatsMoneyOperationSnapshoot", "error", err)
		return err
	}
	return nil
}

func AddCoupAndDivChart(stats []models.StatsMoneyOperationSnapshoot, plot *plot.Plot) error {
	xys := make([]plotter.XY, len(stats))
	xys2 := make([]plotter.XY, len(stats))
	labels := make([]string, len(stats))
	columnLabels := make([]string, len(stats))
	columnLabelPos := make([]plotter.XY, len(stats))

	for i, stat := range stats {
		xys[i] = plotter.XY{
			X: float64(i),
			Y: stat.Coupon,
		}
		xys2[i] = plotter.XY{
			X: float64(i),
			Y: stat.Dividends,
		}
		labels[i] = stat.Time.Format("06/01")
		columnLabels[i] = fmt.Sprintf("%.0f", stat.Coupon+stat.Dividends)
		columnLabelPos[i] = plotter.XY{
			X: float64(i) - 0.2,
			Y: stat.Coupon + stat.Dividends + 50,
		}
	}

	data := ColumnsData{
		XYs:            xys,
		XYs2:           xys2,
		Legend:         []string{"купоны", "дивиденты"},
		Labels:         labels,
		ColumnLabels:   columnLabels,
		ColumnLabelPos: columnLabelPos,
		BarWidth:       1 * vg.Centimeter,
	}

	if err := addBarChart(data, plot); err != nil {
		slog.Error("Failed to add chart for AddBarChartCoupAndDiv", "error", err)
		return err
	}
	return nil
}

func AddSumPriceTotalWithDivChart(p models.Portfolio, plot *plot.Plot) error {
	statsDivPerTicker := stats.GetStatMoneyOperationsSumDivPerTicker(p.MoneyOperations)
	stats := stats.GetLastStatShare(p.Operations)
	stats = common.FilterValue(stats, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})

	xys := make([]plotter.XY, len(stats))
	xys2 := make([]plotter.XY, len(stats))
	labels := make([]string, len(stats))
	columnLabels := make([]string, len(stats))
	columnLabelPos := make([]plotter.XY, len(stats))

	statsKV := common.SortValue(stats, func(i, j models.StatsShare) bool {
		return i.SumPriceTotal > j.SumPriceTotal
	})

	i := 0
	for _, kv := range statsKV {
		xys[i] = plotter.XY{
			X: float64(i),
			Y: statsDivPerTicker[kv.Key],
		}
		xys2[i] = plotter.XY{
			X: float64(i),
			Y: kv.Value.SumPriceTotal - statsDivPerTicker[kv.Key],
		}
		labels[i] = kv.Key
		columnLabels[i] = fmt.Sprintf("%.1f%%", (statsDivPerTicker[kv.Key]/kv.Value.SumPriceTotal)*100)
		columnLabelPos[i] = plotter.XY{
			X: float64(i) - 0.2,
			Y: kv.Value.SumPriceTotal + 10,
		}
		i++
	}

	data := ColumnsData{
		XYs:            xys,
		XYs2:           xys2,
		Legend:         []string{"дивиденты", "остальное"},
		Colors:         []color.Color{color.RGBA{R: 255, A: 255}, color.RGBA{B: 255, A: 255}},
		Labels:         labels,
		ColumnLabels:   columnLabels,
		ColumnLabelPos: columnLabelPos,
		BarWidth:       1 * vg.Centimeter,
	}

	if err := addBarChart(data, plot); err != nil {
		slog.Error("Failed to add chart for AddHistogramSumPriceTotalWithDiv", "error", err)
		return err
	}
	return nil
}

func AddSumPriceTotalChart(stats map[string]models.StatsShare, plot *plot.Plot) error {
	xys := make([]plotter.XY, len(stats))
	labels := make([]string, len(stats))

	statsKV := common.SortValue(stats, func(i, j models.StatsShare) bool {
		return i.SumPriceTotal > j.SumPriceTotal
	})

	i := 0
	for _, kv := range statsKV {
		xys[i] = plotter.XY{
			X: float64(i),
			Y: kv.Value.SumPriceTotal,
		}
		labels[i] = kv.Key
		i++
	}

	data := ColumnsData{
		XYs:    xys,
		Labels: labels,
	}

	if err := addHistogramChart(data, plot); err != nil {
		slog.Error("Failed to add chart for AddHistogramSumPriceTotal", "error", err)
		return err
	}
	return nil
}

func AddSumDivTotalChart(stats map[string]float64, plot *plot.Plot) error {
	xys := make([]plotter.XY, len(stats))
	labels := make([]string, len(stats))

	statsKV := common.SortValue(stats, func(i, j float64) bool {
		return i > j
	})

	i := 0
	for _, kv := range statsKV {
		xys[i] = plotter.XY{
			X: float64(i),
			Y: kv.Value,
		}
		labels[i] = kv.Key
		i++
	}

	data := ColumnsData{
		XYs:    xys,
		Labels: labels,
	}

	if err := addHistogramChart(data, plot); err != nil {
		slog.Error("Failed to add chart for AddHistogramSumDivTotal", "error", err)
		return err
	}
	return nil
}

func AddSumDivFutureChart(stats map[string]models.StatsShare, plot *plot.Plot) error {
	xys := make([]plotter.XY, len(stats))
	labels := make([]string, len(stats))
	columnLabels := make([]string, len(stats))
	columnLabelPos := make([]plotter.XY, len(stats))

	statsKV := common.SortValue(stats, func(i, j models.StatsShare) bool {
		return i.DivPerc > j.DivPerc
	})

	i := 0
	for _, kv := range statsKV {
		xys[i] = plotter.XY{
			X: float64(i),
			Y: kv.Value.DivPerc,
		}
		labels[i] = kv.Key
		columnLabels[i] = fmt.Sprintf("%.0f", kv.Value.SumDiv)
		columnLabelPos[i] = plotter.XY{
			X: float64(i) - 0.2,
			Y: kv.Value.DivPerc + 0.3,
		}
		i++
	}

	data := ColumnsData{
		XYs:            xys,
		Labels:         labels,
		ColumnLabels:   columnLabels,
		ColumnLabelPos: columnLabelPos,
	}

	if err := addHistogramChart(data, plot); err != nil {
		slog.Error("Failed to add chart for AddHistogramSumDivFuture", "error", err)
		return err
	}
	return nil
}
