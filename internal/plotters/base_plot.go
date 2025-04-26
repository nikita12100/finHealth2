package plotters

import (
	"bytes"
	"fmt"
	"log/slog"
	"math"

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

func renderPlot(plot *plot.Plot) *bytes.Buffer {
	writer, err := plot.WriterTo(36*vg.Centimeter, 27*vg.Centimeter, "png")
	if err != nil {
		slog.Error("Error while rendering plot")
		return nil
	}

	var buf bytes.Buffer
	_, err = writer.WriteTo(&buf)
	if err != nil {
		slog.Error("Error while saving rendered plot")
		return nil
	}

	return &buf
}

func GetPlot[T any](
	title string,
	yLabel string,
	ticks int,
	data T,
	addData func(T, *plot.Plot) error,
) *tele.Photo {
	plot := initPlot(title, yLabel, ticks)
	err := addData(data, plot)
	if err != nil {
		slog.Error("Error while adding data on plot", "title", title)
		return nil
	}

	plotBuffer := renderPlot(plot)
	photoFile := tele.FromReader(bytes.NewReader(plotBuffer.Bytes()))

	return &tele.Photo{File: photoFile}
}
