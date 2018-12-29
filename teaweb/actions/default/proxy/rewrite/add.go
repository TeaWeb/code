package rewrite

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaGo/actions"
	"regexp"
)

type AddAction actions.Action

// 添加重写规则
func (this *AddAction) Run(params struct {
	Filename       string
	Index          int
	Pattern        string
	Replace        string
	ProxyId        string
	TargetType     string
	RedirectMethod string
	Must           *actions.Must
}) {
	//@TODO proxyId 支持一个Host

	params.Must.
		Field("pattern", params.Pattern).
		Require("请输入匹配规则").

		Field("targetType", params.TargetType).
		In([]string{"url", "proxy"}, "目标类型错误")

	if params.TargetType == "proxy" {
		params.Must.
			Field("proxyId", params.ProxyId).
			Require("请选择目标代理")
	}

	params.Must.
		Field("replace", params.Replace).
		Require("请输入目标URL")

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if len(params.Replace) == 0 {
		params.Replace = "/"
	} else if params.Replace[0] != '/' && !regexp.MustCompile("(?i)^(http|https|ftp)://").MatchString(params.Replace) {
		params.Replace = "/" + params.Replace
	}

	location := proxy.LocationAtIndex(params.Index)
	if location != nil {
		rewriteRule := teaconfigs.NewRewriteRule()
		rewriteRule.On = true
		rewriteRule.Pattern = params.Pattern
		if params.TargetType == "url" {
			rewriteRule.Replace = params.Replace
		} else {
			rewriteRule.Replace = "proxy://" + params.ProxyId + params.Replace
		}
		if len(params.RedirectMethod) > 0 {
			rewriteRule.AddFlag(params.RedirectMethod, nil)
		}
		location.Rewrite = append(location.Rewrite, rewriteRule)
	}

	err = proxy.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	global.NotifyChange()

	this.Refresh().Success("添加成功")
}
