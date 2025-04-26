package plotters

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"test2/internal/common"
	"test2/internal/models"
	"test2/internal/stats"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"

	tele "gopkg.in/telebot.v4"
)

type integerTicks struct{ Step int }

func (tic integerTicks) Ticks(min, max float64) []plot.Tick {
	var t []plot.Tick
	for i := math.Trunc(min); i <= max; i += float64(tic.Step) {
		t = append(t, plot.Tick{Value: i, Label: fmt.Sprint(i)})
	}
	return t
}

func initPlot(title string, yLabel string, ticks int) *plot.Plot {
	p := plot.New()
	p.Add(plotter.NewGrid())

	p.Title.Text = title
	p.Title.TextStyle.Font.Size = 16
	p.Title.Padding = vg.Length(10) // Space above title

	p.X.Padding = vg.Length(5) // Space below X axis
	p.X.Tick.Label.Rotation = math.Pi / 4
	p.X.Tick.Label.XAlign = text.XRight // Align to right of tick
	p.X.Tick.Label.YAlign = text.YTop   // Align above tick

	// p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02"}

	p.Y.Label.Text = yLabel
	p.Y.Tick.Marker = integerTicks{Step: ticks}

	p.Legend.Top = true              // Place at top
	p.Legend.Left = true             // Not on left side
	p.Legend.XOffs = 0               // Horizontal offset
	p.Legend.YOffs = -50             // Vertical offset
	p.Legend.Padding = vg.Length(10) // Space around legend
	p.Legend.TextStyle.Font.Size = 14

	return p
}

func renderPlot(plot *plot.Plot) (*bytes.Buffer, error) {
	writer, err := plot.WriterTo(36*vg.Centimeter, 27*vg.Centimeter, "png")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = writer.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func GetPlot[T any](
	title string,
	yLabel string,
	ticks int,
	data T,
	addData func(T, *plot.Plot) error,
) (*tele.Photo, error) {
	plot := initPlot(title, yLabel, ticks)
	err := addData(data, plot)
	if err != nil {
		return nil, err
	}

	plotBuffer, err := renderPlot(plot)
	if err != nil {
		return nil, err
	}

	return &tele.Photo{File: tele.FromReader(bytes.NewReader(plotBuffer.Bytes()))}, nil
}

func GetLine(data []models.StatsMoneyOperationSnapshoot) (*plotter.Line, error) {
	pts := make(plotter.XYs, len(data))
	for i, d := range data {
		pts[i].X = float64(d.Time.Unix()) // Convert time to float
		pts[i].Y = d.Replenishment
		// pts[i].Y = d.Coupon + d.Dividends
	}

	line, _ := plotter.NewLine(pts)
	line.Color = color.RGBA{R: 255, A: 255}

	return line, nil
}

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
