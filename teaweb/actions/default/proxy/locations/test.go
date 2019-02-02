package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type TestAction actions.Action

// 测试
func (this *TestAction) Run(params struct {
	Pattern           string
	PatternType       int
	IsReverse         bool
	IsCaseInsensitive bool
	TestingPath       string
}) {
	location := teaconfigs.NewLocation()
	location.Pattern = params.Pattern
	location.SetPattern(params.Pattern, params.PatternType, params.IsCaseInsensitive, params.IsReverse)
	err := location.Validate()
	if err != nil {
		this.Fail("校验失败：" + err.Error())
	}
	mapping, ok := location.Match(params.TestingPath)
	if ok {
		this.Data["mapping"] = mapping
		this.Success()
	} else {
		this.Data["mapping"] = nil
		this.Fail()
	}
}
