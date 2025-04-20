package plotters

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"test2/internal/models"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
)

type integerTicks struct{ Step int }

func (tic integerTicks) Ticks(min, max float64) []plot.Tick {
	var t []plot.Tick
	for i := math.Trunc(min); i <= max; i += float64(tic.Step) {
		t = append(t, plot.Tick{Value: i, Label: fmt.Sprint(i)})
	}
	return t
}

func InitPlot(title string, yLabel string) *plot.Plot {
	p := plot.New()

	p.Title.Text = title
	p.Title.TextStyle.Font.Size = 16
	p.Title.Padding = vg.Length(10) // Space above title

	p.X.Padding = vg.Length(5) // Space below X axis
	p.X.Tick.Label.Rotation = math.Pi / 4
	p.X.Tick.Label.XAlign = text.XRight // Align to right of tick
	p.X.Tick.Label.YAlign = text.YTop   // Align above tick

	// p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02"}

	p.Y.Label.Text = yLabel
	p.Y.Tick.Marker = integerTicks{Step: 1000}

	p.Legend.Top = true              // Place at top
	p.Legend.Left = true             // Not on left side
	p.Legend.XOffs = 0               // Horizontal offset
	p.Legend.YOffs = -50             // Vertical offset
	p.Legend.Padding = vg.Length(10) // Space around legend
	p.Legend.TextStyle.Font.Size = 14

	return p
}

func RenderPlot(plot *plot.Plot) (*bytes.Buffer, error) {
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

// https://github.com/gonum/plot/issues/656
func AddHistogramCoupAndDiv(stats []models.StatsMoneyOperationSnapshoot, plot *plot.Plot) error {
	ptsCoupon := make(plotter.Values, len(stats))
	ptsDiv := make(plotter.Values, len(stats))
	labels := make([]string, len(stats))

	var labelsText []string
	var labelsPos []plotter.XY
	for i, stat := range stats {
		labels[i] = stat.Time.Month().String()
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
	plot.Legend.Add("coupon", barsCoupon)
	plot.Legend.Add("div", barsDiv)

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
