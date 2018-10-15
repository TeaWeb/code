package teacharts

import "github.com/TeaWeb/code/teainterfaces"

func ConvertInterface(chart teainterfaces.ChartInterface) ChartInterface {
	if chart.Type() == "progressBar" {
		return NewProgressBarFromInterface(chart.(teainterfaces.ProgressBarInterface))
	}

	return nil
}
