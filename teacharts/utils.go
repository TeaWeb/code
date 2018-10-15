package teacharts

import "github.com/TeaWeb/code/teainterfaces"

func ConvertInterface(chart teainterfaces.ChartInterface) ChartInterface {
	switch chart.Type() {
	case "progressBar":
		c, ok := chart.(teainterfaces.ProgressBarInterface)
		if !ok {
			return nil
		}
		return NewProgressBarFromInterface(c)
	case "pie":
		c, ok := chart.(teainterfaces.PieChartInterface)
		if !ok {
			return nil
		}
		return NewPieChartFromInterface(c)
	case "line":
		c, ok := chart.(teainterfaces.LineChartInterface)
		if !ok {
			return nil
		}
		return NewLineChartFromInterface(c)
	case "gauge":
		c, ok := chart.(teainterfaces.GaugeChartInterface)
		if !ok {
			return nil
		}
		return NewGaugeChartFromInterface(c)
	case "table":
		c, ok := chart.(teainterfaces.TableInterface)
		if !ok {
			return nil
		}
		return NewTableFromInterface(c)
	}

	return nil
}
