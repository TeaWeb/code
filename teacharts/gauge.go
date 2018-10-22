package teacharts

import (
	"github.com/TeaWeb/plugin/charts"
)

func NewGaugeChart() *GaugeChart {
	p := &GaugeChart{}
	p.Type = "gauge"
	return p
}

func NewGaugeChartFromInterface(chart *charts.GaugeChart) *GaugeChart {
	p := &GaugeChart{}
	p.Type = "gauge"
	p.Name = chart.Name
	p.Detail = chart.Detail

	p.Value = chart.Value
	p.Label = chart.Label
	p.Min = chart.Min
	p.Max = chart.Max
	p.Unit = chart.Unit
	return p
}

type GaugeChart struct {
	Chart

	Value float64 `json:"value"`
	Label string  `json:"label"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Unit  string  `json:"unit"`
}
