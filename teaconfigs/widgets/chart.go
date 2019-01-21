package widgets

import (
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

// Chart接口
type ChartInterface interface {
	AsJavascript(options map[string]interface{}) (code string, err error)
}

// Chart定义
type Chart struct {
	Id      string                 `yaml:"id" json:"id"`
	On      bool                   `yaml:"on" json:"on"`
	Name    string                 `yaml:"name" json:"name"`
	Columns uint8                  `yaml:"columns" json:"columns"` // 列
	Type    string                 `yaml:"type" json:"type"`       // 类型
	Options map[string]interface{} `yaml:"options" json:"options"`
}

// 获取新对象
func NewChart() *Chart {
	return &Chart{
		On: true,
		Id: stringutil.Rand(16),
	}
}

// 校验
func (this *Chart) Validate() error {
	return nil
}

// 转换为具体对象
func (this *Chart) AsObject() (ChartInterface, error) {
	for _, chart := range AllChartTypes {
		if chart["code"] != this.Type {
			continue
		}
		instance, ok := chart["instance"].(ChartInterface)
		if ok {
			err := teautils.MapToObjectJSON(this.Options, instance)
			return instance, err
		} else {
			return nil, errors.New("chart instance should implement ChartInterface: '" + this.Type + "'")
		}
	}

	return nil, errors.New("invalid chart type '" + this.Type + "'")
}
