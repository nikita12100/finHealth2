package plotters

import (
	"bytes"
	"image/color"
	"test2/internal/stats"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func CreatePlot(data []stats.TimeValue) (*bytes.Buffer, error) {
	pts := make(plotter.XYs, len(data))
	for i, d := range data {
		pts[i].X = float64(d.Time.Unix()) // Convert time to float
		pts[i].Y = d.Value
	}

	p := plot.New()
	p.Title.Text = "Time Series"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02"}
	p.Y.Label.Text = "Values"

	line, _ := plotter.NewLine(pts)
	line.Color = color.RGBA{R: 255, A: 255}
	p.Add(line)
	p.Legend.Add("Replenishment", line)

	writer, err := p.WriterTo(8*vg.Inch, 4*vg.Inch, "png")
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