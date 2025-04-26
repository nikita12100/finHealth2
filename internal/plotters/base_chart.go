package plotters

import (
	"image/color"
	"log/slog"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type ColumnsData struct {
	XYs            []plotter.XY  // координаты
	XYs2           []plotter.XY  // координаты 2 (сверху)
	Legend         []string      // легенда
	Colors         []color.Color // цвет столбцов
	Labels         []string      // подписи оси х
	ColumnLabels   []string      // подписи столбика
	ColumnLabelPos []plotter.XY  // позиция подписи столбика
	BarWidth       vg.Length     // ширина столбика
}

func addBarChart(data ColumnsData, plot *plot.Plot) error {
	values := make(plotter.Values, len(data.XYs))
	for i, currDataXY := range data.XYs {
		values[i] = currDataXY.Y
	}

	bars, err := plotter.NewBarChart(values, data.BarWidth)
	if err != nil {
		slog.Error("Failed to build chart", "error", err)
		return err
	}

	if len(data.XYs2) == len(data.XYs) {
		values2 := make(plotter.Values, len(data.XYs2))
		for i, currDataXY2 := range data.XYs2 {
			values2[i] = currDataXY2.Y
		}

		barsTop, err := plotter.NewBarChart(values2, data.BarWidth)
		if err != nil {
			slog.Error("Failed to build chart", "error", err)
			return err
		}

		barsTop.StackOn(bars)
		if len(data.Colors) != 2 {
			data.Colors = []color.Color{
				color.RGBA{B: 255, A: 255},
				color.RGBA{R: 255, A: 255},
			}
		}

		bars.Color = data.Colors[0]
		barsTop.Color = data.Colors[1]

		plot.Add(bars, barsTop)

		if len(data.Legend) == 2 {
			plot.Legend.Add(data.Legend[0], bars)
			plot.Legend.Add(data.Legend[1], barsTop)
		}
	} else {
		plot.Add(bars)
	}

	plot.NominalX(data.Labels...)

	if len(data.ColumnLabels) == len(data.XYs) {
		columnLabels, err := plotter.NewLabels(plotter.XYLabels{
			XYs:    data.ColumnLabelPos,
			Labels: data.ColumnLabels,
		},
		)
		if err != nil {
			slog.Error("Could not creates labels column", "error", err)
			return err
		}
		for i := range columnLabels.TextStyle {
			columnLabels.TextStyle[i].Color = color.Black
			columnLabels.TextStyle[i].Font.Size = 12
		}

		plot.Add(columnLabels)
	}

	return nil
}

func addHistogramChart(data ColumnsData, plot *plot.Plot) error {
	var xys plotter.XYs
	xys = data.XYs
	bars, err := plotter.NewHistogram(xys, len(data.XYs)*2)
	if err != nil {
		slog.Error("Failed to build chart", "error", err)
		return err
	}

	plot.Add(bars)
	plot.NominalX(data.Labels...)

	if len(data.ColumnLabels) == len(data.XYs) {
		columnLabels, err := plotter.NewLabels(plotter.XYLabels{
			XYs:    data.ColumnLabelPos,
			Labels: data.ColumnLabels,
		},
		)
		if err != nil {
			slog.Error("Could not creates labels column", "error", err)
			return err
		}
		for i := range columnLabels.TextStyle {
			columnLabels.TextStyle[i].Color = color.Black
			columnLabels.TextStyle[i].Font.Size = 12
		}

		plot.Add(columnLabels)
	}

	return nil
}
