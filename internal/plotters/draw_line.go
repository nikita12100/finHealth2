package plotters

import (
	"fmt"
	"image/color"
	"log"
	"test2/internal/common"
	"test2/internal/models"
	"test2/internal/stats"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)


func AddBarChart(stats []models.StatsMoneyOperationSnapshoot, plot *plot.Plot) error {
	pts := make(plotter.Values, len(stats))
	labels := make([]string, len(stats))

	var labelsText []string
	var labelsPos []plotter.XY
	for i, stat := range stats {
		labels[i] = stat.Time.Format("06/01")
		pts[i] = stat.Replenishment

		rLabel := fmt.Sprintf("%.0f", pts[i])
		labelsText = append(labelsText, rLabel)

		rxPos := float64(i) - 0.25
		ryPos := pts[i] + 50

		rXY := plotter.XY{X: rxPos, Y: ryPos}
		labelsPos = append(labelsPos, rXY)
	}

	hist, err := plotter.NewBarChart(pts, 20)
	if err != nil {
		return err
	}

	plot.Add(hist)
	plot.NominalX(labels...)

	rl, err := plotter.NewLabels(plotter.XYLabels{
		XYs:    labelsPos,
		Labels: labelsText,
	},
	)
	if err != nil {
		log.Fatalf("could not creates labels plotter: %+v", err)
	}
	for i := range rl.TextStyle {
		rl.TextStyle[i].Color = color.Black
		rl.TextStyle[i].Font.Size = 12
	}

	plot.Add(rl)

	return nil
}

func AddBarChartCoupAndDiv(stats []models.StatsMoneyOperationSnapshoot, plot *plot.Plot) error {
	ptsCoupon := make(plotter.Values, len(stats))
	ptsDiv := make(plotter.Values, len(stats))
	labels := make([]string, len(stats))

	var labelsText []string
	var labelsPos []plotter.XY
	for i, stat := range stats {
		labels[i] = stat.Time.Format("06/01")
		ptsCoupon[i] = stat.Coupon
		ptsDiv[i] = stat.Dividends

		rLabel := fmt.Sprintf("%.0f", ptsCoupon[i]+ptsDiv[i])
		labelsText = append(labelsText, rLabel)

		rxPos := float64(i) - 0.2
		ryPos := ptsCoupon[i] + ptsDiv[i] + 50

		rXY := plotter.XY{X: rxPos, Y: ryPos}
		labelsPos = append(labelsPos, rXY)
	}

	barsCoupon, err := plotter.NewBarChart(ptsCoupon, 1*vg.Centimeter)
	if err != nil {
		return err
	}
	barsDiv, err := plotter.NewBarChart(ptsDiv, 1*vg.Centimeter)
	if err != nil {
		return err
	}

	barsDiv.StackOn(barsCoupon)
	barsCoupon.Color = color.RGBA{B: 255, A: 255}
	barsDiv.Color = color.RGBA{R: 255, A: 255}

	plot.Add(barsCoupon, barsDiv)
	plot.Legend.Add("купоны", barsCoupon)
	plot.Legend.Add("дивиденты", barsDiv)

	plot.NominalX(labels...)

	rl, err := plotter.NewLabels(plotter.XYLabels{
		XYs:    labelsPos,
		Labels: labelsText,
	},
	)
	if err != nil {
		log.Fatalf("could not creates labels plotter: %+v", err)
	}
	for i := range rl.TextStyle {
		rl.TextStyle[i].Color = color.Black
		rl.TextStyle[i].Font.Size = 12
	}

	plot.Add(rl)

	return nil
}

func AddHistogramSumPriceTotal(stats map[string]models.StatsShare, plot *plot.Plot) error {
	pts := make(plotter.XYs, len(stats))
	labels := make([]string, len(stats))

	statsKV := common.SortValue(stats, func(i, j models.StatsShare) bool {
		return i.SumPriceTotal > j.SumPriceTotal
	})

	i := 0
	for _, kv := range statsKV {
		labels[i] = kv.Key
		pts[i].X = float64(i)
		pts[i].Y = kv.Value.SumPriceTotal

		i++
	}

	hist, err := plotter.NewHistogram(pts, len(stats)*2)
	if err != nil {
		return err
	}

	plot.Add(hist)
	plot.NominalX(labels...)

	return nil
}

func AddHistogramSumDivTotal(stats map[string]float64, plot *plot.Plot) error {
	pts := make(plotter.XYs, len(stats))
	labels := make([]string, len(stats))

	statsKV := common.SortValue(stats, func(i, j float64) bool {
		return i > j
	})

	i := 0
	for _, kv := range statsKV {
		labels[i] = kv.Key
		pts[i].X = float64(i)
		pts[i].Y = kv.Value

		i++
	}

	hist, err := plotter.NewHistogram(pts, len(stats)*2)
	if err != nil {
		return err
	}

	plot.Add(hist)
	plot.NominalX(labels...)

	return nil
}

func AddHistogramSumDivFuture(stats map[string]models.StatsShare, plot *plot.Plot) error {
	pts := make(plotter.XYs, len(stats))
	labels := make([]string, len(stats))

	statsKV := common.SortValue(stats, func(i, j models.StatsShare) bool {
		return i.DivPerc > j.DivPerc
	})

	var labelsText []string
	var labelsPos []plotter.XY
	i := 0
	for _, kv := range statsKV {
		labels[i] = kv.Key
		pts[i].X = float64(i)
		pts[i].Y = kv.Value.DivPerc

		rLabel := fmt.Sprintf("%.0f", kv.Value.SumDiv)
		labelsText = append(labelsText, rLabel)

		rxPos := pts[i].X - 0.2
		ryPos := pts[i].Y + 0.3 // высота лейбла от столбца

		rXY := plotter.XY{X: rxPos, Y: ryPos}
		labelsPos = append(labelsPos, rXY)

		i++
	}

	hist, err := plotter.NewHistogram(pts, len(stats)*2)
	if err != nil {
		return err
	}

	plot.Add(hist)
	plot.NominalX(labels...)

	rl, err := plotter.NewLabels(plotter.XYLabels{
		XYs:    labelsPos,
		Labels: labelsText,
	},
	)
	if err != nil {
		log.Fatalf("could not creates labels plotter: %+v", err)
	}
	for i := range rl.TextStyle {
		rl.TextStyle[i].Color = color.Black
		rl.TextStyle[i].Font.Size = 12
	}

	plot.Add(rl)

	return nil
}

func AddHistogramSumPriceTotalWithDiv(p models.Portfolio, plot *plot.Plot) error {
	statsDivPerTicker := stats.GetStatMoneyOperationsSumDivPerTicker(p.MoneyOperations)
	stats := stats.GetLastStatShare(p.Operations)
	stats = common.FilterValue(stats, func(stat models.StatsShare) bool {
		return stat.Count != 0
	})


	ptsSum := make(plotter.Values, len(stats))
	ptsDiv := make(plotter.Values, len(stats))
	labels := make([]string, len(stats))

	statsKV := common.SortValue(stats, func(i, j models.StatsShare) bool {
		return i.SumPriceTotal > j.SumPriceTotal
	})

	var labelsText []string
	var labelsPos []plotter.XY
	i := 0
	for _, kv := range statsKV {
		labels[i] = kv.Key
		// ptsSum[i].X = float64(i)
		ptsSum[i] = kv.Value.SumPriceTotal - statsDivPerTicker[kv.Key]
		ptsDiv[i] = statsDivPerTicker[kv.Key]

		rLabel := fmt.Sprintf("%.1f%%", (statsDivPerTicker[kv.Key] / kv.Value.SumPriceTotal) * 100)
		labelsText = append(labelsText, rLabel)

		rxPos := float64(i) - 0.2
		ryPos := ptsSum[i] + ptsDiv[i] + 10 // высота лейбла от столбца

		rXY := plotter.XY{X: rxPos, Y: ryPos}
		labelsPos = append(labelsPos, rXY)
		i++
	}

	barsSum, err := plotter.NewBarChart(ptsSum, 1*vg.Centimeter)
	if err != nil {
		return err
	}
	barsDiv, err := plotter.NewBarChart(ptsDiv, 1*vg.Centimeter)
	if err != nil {
		return err
	}

	// barsDiv.StackOn(barsSum)
	barsSum.StackOn(barsDiv)
	barsSum.Color = color.RGBA{B: 255, A: 255}
	barsDiv.Color = color.RGBA{R: 255, A: 255}

	plot.Add(barsSum, barsDiv)
	plot.Legend.Add("остальное", barsSum)
	plot.Legend.Add("дивиденты", barsDiv)

	plot.NominalX(labels...)

	rl, err := plotter.NewLabels(plotter.XYLabels{
		XYs:    labelsPos,
		Labels: labelsText,
	},
	)
	if err != nil {
		log.Fatalf("could not creates labels plotter: %+v", err)
	}
	for i := range rl.TextStyle {
		rl.TextStyle[i].Color = color.Black
		rl.TextStyle[i].Font.Size = 12
	}

	plot.Add(rl)

	return nil
}
