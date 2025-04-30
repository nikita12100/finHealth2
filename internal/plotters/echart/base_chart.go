package echart

import (
	"log/slog"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type ColumnsData struct {
	Data   []opts.BarData // координаты
	Data2  []opts.BarData // координаты 2 (сверху)
	Legend []string       // легенда
	Labels []string       // подписи оси х
}

type ChartOptions struct {
	Title      string
	Subtitle   string
	EnableZoom bool
}

func validateData(data ColumnsData) bool {
	if len(data.Data) != len(data.Labels) {
		slog.Error("Len labeles not equal data len")
		return false
	}
	if len(data.Data) > 0 && len(data.Legend) == 0 {
		slog.Error("Empty Legend for data1")
		return false
	}
	if len(data.Data2) > 0 && len(data.Legend) < 1 {
		slog.Error("Empty Legend for data2")
		return false
	}
	if len(data.Data2) > 0 && len(data.Data) != len(data.Data2) {
		slog.Error("Data and data2 len not equal")
		return false
	}

	return true
}

func setGlobalOptions(bar *charts.Bar, options ChartOptions) {
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    options.Title,
			Subtitle: options.Subtitle,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
			AxisLabel: &opts.AxisLabel{
				Rotate: 45,
			},
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Right: "50",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Type:  "jpg",
					Title: "Скачать график",
				},
				DataView: &opts.ToolBoxFeatureDataView{
					Title: "Изменить данные",
					Lang:  []string{"Данные", "отмена", "обновить"},
				},
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			AxisLabel: &opts.AxisLabel{Formatter: "{value} руб."},
		}),
	)

	if options.EnableZoom {
		bar.SetGlobalOptions(
			charts.WithDataZoomOpts(opts.DataZoom{
				Type:  "slider",
				Start: 50,
				End:   100,
			}),
		)
	}
}

func getChart(data ColumnsData, options ChartOptions) *components.Page {
	if isCorrect := validateData(data); !isCorrect {
		return nil
	}

	bar := charts.NewBar()
	setGlobalOptions(bar, options)

	bar.SetXAxis(data.Labels)
	bar.AddSeries(data.Legend[0], data.Data).
		SetSeriesOptions(charts.WithBarChartOpts(opts.BarChart{
			Stack: "stackA",
		}))

	if len(data.Data2) == len(data.Data) {
		bar.AddSeries(data.Legend[1], data.Data2).
			SetSeriesOptions(charts.WithBarChartOpts(opts.BarChart{
				Stack: "stackA",
			}))
	}

	page := components.NewPage()
	page.AddCharts(bar)

	return page
}
