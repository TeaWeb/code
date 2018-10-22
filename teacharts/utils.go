package teacharts

import (
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/plugin/charts"
)

func ConvertInterface(chart interface{}) ChartInterface {
	switch c := chart.(type) {
	case *charts.ProgressBar:
		return NewProgressBarFromInterface(c)
	case *charts.PieChart:
		return NewPieChartFromInterface(c)
	case *charts.LineChart:
		return NewLineChartFromInterface(c)
	case *charts.GaugeChart:
		return NewGaugeChartFromInterface(c)
	case *charts.Table:
		return NewTableFromInterface(c)
	case map[string]interface{}:
		chartType, found := c["ChartType"]
		if found {
			switch chartType {
			case "progressBar":
				c2 := new(ProgressBar)
				teautils.MapToObjectJSON(c, c2)
				c2.Type = chartType.(string)
				return c2
			case "pie":
				c2 := new(PieChart)
				teautils.MapToObjectJSON(c, c2)
				c2.Type = chartType.(string)
				return c2
			case "line":
				c2 := new(LineChart)
				teautils.MapToObjectJSON(c, c2)
				c2.Type = chartType.(string)
				return c2
			case "gauge":
				c2 := new(GaugeChart)
				teautils.MapToObjectJSON(c, c2)
				c2.Type = chartType.(string)
				return c2
			case "table":
				c2 := new(Table)
				teautils.MapToObjectJSON(c, c2)
				c2.Type = chartType.(string)
				return c2
			}
		}
	}

	return nil
}
