package rewrite

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"regexp"
)

type TestAction actions.Action

// 匹配测试
func (this *TestAction) Run(params struct {
	Pattern      string
	Replace      string
	ProxyId      string
	TargetType   string
	RedirectMode string
	CondParams   []string
	CondOps      []string
	CondValues   []string
	TestingPath  string
	Must         *actions.Must
}) {
	params.Must.
		Field("pattern", params.Pattern).
		Require("请输入匹配规则").
		Expect(func() (message string, success bool) {
			_, err := regexp.Compile(params.Pattern)
			if err != nil {
				return "匹配规则错误：" + err.Error(), false
			}
			return "", true
		})

	rewriteRule := teaconfigs.NewRewriteRule()
	rewriteRule.On = true
	rewriteRule.Pattern = params.Pattern
	if params.TargetType == "url" {
		rewriteRule.Replace = params.Replace
	} else {
		rewriteRule.Replace = "proxy://" + params.ProxyId + params.Replace
	}
	if len(params.RedirectMode) > 0 {
		rewriteRule.AddFlag(params.RedirectMode, nil)
	}

	if len(params.CondParams) > 0 {
		for index, param := range params.CondParams {
			if index < len(params.CondOps) && index < len(params.CondValues) {
				cond := teaconfigs.NewRewriteCond()
				cond.Param = param
				cond.Value = params.CondValues[index]
				cond.Operator = params.CondOps[index]
				err := cond.Validate()
				if err != nil {
					this.Fail("过滤条件\"" + cond.Param + " " + cond.Value + "\"校验失败：" + err.Error())
				}
				rewriteRule.AddCond(cond)
			}
		}
	}

	err := rewriteRule.Validate()
	if err != nil {
		this.Fail("校验失败：" + err.Error())
	}

	replace, mapping, ok := rewriteRule.Match(params.TestingPath, func(source string) string {
		return source
	})
	if ok {
		this.Data["replace"] = replace
		this.Data["mapping"] = mapping
		this.Success()
	} else {
		this.Fail()
	}
}
