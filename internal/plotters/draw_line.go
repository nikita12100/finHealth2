package plotters

import (
	"bytes"
	"image/color"
	"test2/internal/models"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func InitPlot() *plot.Plot {
	p := plot.New()
	p.Title.Text = "Time Series"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02"}
	p.Y.Label.Text = "Values"

	return p
}

func RenderPlot(plot *plot.Plot) (*bytes.Buffer, error) {
	writer, err := plot.WriterTo(8*vg.Inch, 4*vg.Inch, "png")
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

func AddHistogram(stats []models.StatsMoneyOperationSnapshoot, plot *plot.Plot) error {
	pts := make(plotter.Values, len(stats))
	labels := make([]string, len(stats))
	for i, stat := range stats {
		labels[i] = stat.Time.Month().String()
		pts[i] = stat.Replenishment
	}

	hist, err := plotter.NewBarChart(pts, 20)
	if err != nil {
		return err
	}

	plot.Add(hist)
	plot.Legend.Add("Replenishment", hist)

	plot.NominalX(labels...)

	return nil
}

func AddHistogram2(stats []models.StatsMoneyOperationSnapshoot, plot *plot.Plot) error {
	ptsCoupon := make(plotter.Values, len(stats))
	ptsDiv := make(plotter.Values, len(stats))
	labels := make([]string, len(stats))
	for i, stat := range stats {
		labels[i] = stat.Time.Month().String()
		ptsCoupon[i] = stat.Coupon
		ptsDiv[i] = stat.Dividends
	}

	barsBase, err := plotter.NewBarChart(ptsCoupon, 20)
	barsTop, err := plotter.NewBarChart(ptsDiv, 20)
	if err != nil {
		return err
	}

	barsTop.StackOn(barsBase) // Stack on top of base
	barsBase.Color = color.RGBA{B: 255, A: 255}
	barsTop.Color = color.RGBA{R: 255, A: 255}

	plot.Add(barsBase, barsTop)
	plot.Legend.Add("coupon", barsBase)
	plot.Legend.Add("div", barsTop)

	plot.NominalX(labels...)

	return nil
}
