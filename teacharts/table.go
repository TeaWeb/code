package teacharts

import (
	"github.com/TeaWeb/plugin/charts"
	"sync"
)

type Row struct {
	Columns []*Column `json:"columns"`
}

type Column struct {
	Text  string  `json:"text"`
	Width float64 `json:"width"` // 百分比，比如 30 表示 30%
}

type Table struct {
	Chart

	Rows   []*Row `json:"rows"`
	locker sync.Mutex

	width []float64
}

func NewTable() *Table {
	p := &Table{
		Rows: []*Row{},
	}
	p.Type = "table"
	return p
}

func NewTableFromInterface(chart *charts.Table) *Table {
	p := &Table{
		Rows: []*Row{},
	}
	p.Type = "table"
	p.Name = chart.Name
	p.Detail = chart.Detail

	for _, r := range chart.Rows {
		texts := []string{}
		widths := []float64{}
		for _, c := range r.Columns {
			texts = append(texts, c.Text)
			widths = append(widths, c.Width)
		}

		p.AddRow(texts ...)
		p.SetWidth(widths ...)
	}

	return p
}

func (this *Table) ResetRows() {
	this.locker.Lock()
	defer this.locker.Unlock()

	this.Rows = []*Row{}
}

func (this *Table) AddRow(text ... string) {
	this.locker.Lock()
	defer this.locker.Unlock()

	columns := []*Column{}
	for index, t := range text {
		if index < len(this.width) {
			columns = append(columns, &Column{
				Text:  t,
				Width: this.width[index],
			})
		} else {
			columns = append(columns, &Column{
				Text: t,
			})
		}
	}
	this.Rows = append(this.Rows, &Row{
		Columns: columns,
	})
}

func (this *Table) SetWidth(wide ... float64) {
	this.width = wide

	for _, row := range this.Rows {
		for index, column := range row.Columns {
			if index < len(this.width) {
				column.Width = this.width[index]
			}
		}
	}
}
