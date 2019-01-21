package agents

import (
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

// 数据源接口定义
type SourceInterface interface {
	// 名称
	Name() string

	// 代号
	Code() string

	// 描述
	Description() string

	// 校验
	Validate() error

	// 执行
	Execute(params map[string]string) (value interface{}, err error)

	// 获得数据格式
	DataFormatCode() SourceDataFormat

	// 获取简要信息
	Summary() maps.Map
}

// 获取所有的数据源信息
func AllDataSources() []maps.Map {
	result := []maps.Map{}
	for _, i := range []SourceInterface{NewScriptSource(), NewWebHookSource(), NewFileSource(),} {
		summary := i.Summary()
		summary["instance"] = i
		result = append(result, summary)
	}
	return result
}

// 查找单个数据源信息
func FindDataSource(code string) maps.Map {
	for _, summary := range AllDataSources() {
		if summary["code"] == code {
			return summary
		}
	}
	return nil
}

// 查找单个数据源实例
func FindDataSourceInstance(code string, options map[string]interface{}) SourceInterface {
	for _, summary := range AllDataSources() {
		if summary["code"] == code {
			instance := summary["instance"].(SourceInterface)
			if options != nil {
				err := teautils.MapToObjectJSON(options, instance)
				if err != nil {
					logs.Error(err)
				}
			}
			return instance
		}
	}
	return nil
}
