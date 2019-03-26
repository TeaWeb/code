package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type DetailAction actions.Action

// 代理详情
func (this *DetailAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if server.Index == nil {
		server.Index = []string{}
	}

	this.Data["selectedTab"] = "basic"
	this.Data["server"] = server

	this.Data["errs"] = teaproxy.SharedManager.FindServerErrors(params.ServerId)

	this.Data["accessLogFields"] = lists.Map(tealogs.AccessLogFields, func(k int, v interface{}) interface{} {
		m := v.(maps.Map)
		m["isChecked"] = len(server.AccessLogFields) == 0 || lists.ContainsInt(server.AccessLogFields, types.Int(m["code"]))
		return m
	})

	this.Show()
}
