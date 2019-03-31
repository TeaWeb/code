package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type TestAction actions.Action

// 测试
func (this *TestAction) Run(params struct {
	Pattern           string
	PatternType       int
	IsReverse         bool
	IsCaseInsensitive bool
	TestingPath       string

	CondParams []string
	CondOps    []string
	CondValues []string
}) {
	location := teaconfigs.NewLocation()
	location.Pattern = params.Pattern
	location.SetPattern(params.Pattern, params.PatternType, params.IsCaseInsensitive, params.IsReverse)

	location.Cond = []*teaconfigs.RequestCond{}
	if len(params.CondParams) > 0 {
		for index, param := range params.CondParams {
			if index < len(params.CondOps) && index < len(params.CondValues) {
				cond := teaconfigs.NewRequestCond()
				cond.Param = param
				cond.Value = params.CondValues[index]
				cond.Operator = params.CondOps[index]
				err := cond.Validate()
				if err != nil {
					this.Fail("匹配条件\"" + cond.Param + " " + cond.Value + "\"校验失败：" + err.Error())
				}
				location.AddCond(cond)
			}
		}
	}

	err := location.Validate()
	if err != nil {
		this.Fail("校验失败：" + err.Error())
	}

	rawReq, err := http.NewRequest(http.MethodGet, params.TestingPath, nil)
	if err != nil {
		this.Fail("请输入正确的URL")
	}

	req := teaproxy.NewRequest(rawReq)
	req.SetURI(params.TestingPath)
	req.SetHost(rawReq.Host)

	mapping, ok := location.Match(rawReq.URL.Path, func(source string) string {
		if req == nil {
			return source
		} else {
			return req.Format(source)
		}
	})
	if ok {
		this.Data["mapping"] = mapping
		this.Success()
	} else {
		this.Data["mapping"] = nil
		this.Fail()
	}
}
