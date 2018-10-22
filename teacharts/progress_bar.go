package teacharts

import (
	"github.com/TeaWeb/plugin/charts"
)

type ProgressBarColor string

func NewProgressBar() *ProgressBar {
	p := &ProgressBar{
		Color: ColorBlue,
	}
	p.Type = "progressBar"
	return p
}

func NewProgressBarFromInterface(chart *charts.ProgressBar) *ProgressBar {
	p := &ProgressBar{
		Color: ColorBlue,
	}
	p.Type = "progressBar"
	p.Name = chart.Name
	p.Detail = chart.Detail

	p.Value = chart.Value
	p.Color = chart.Color
	return p
}

type ProgressBar struct {
	Chart

	Value float64 `json:"value"`
	Color Color   `json:"color"`
}
