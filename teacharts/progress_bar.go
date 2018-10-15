package teacharts

import "github.com/TeaWeb/code/teainterfaces"

type ProgressBarColor string

type ProgressBar struct {
	Chart

	Value float64 `json:"value"`
	Color Color   `json:"color"`
}

func NewProgressBar() *ProgressBar {
	p := &ProgressBar{
		Color: ColorBlue,
	}
	p.Type = "progressBar"
	return p
}

func NewProgressBarFromInterface(chart teainterfaces.ProgressBarInterface) *ProgressBar {
	p := &ProgressBar{
		Color: ColorBlue,
	}
	p.Type = "progressBar"
	p.Value = chart.Value()
	p.Color = chart.Color()
	return p
}

func (this *ProgressBar) UniqueId() string {
	return this.Id
}

func (this *ProgressBar) SetUniqueId(id string) {
	this.Id = id
}
