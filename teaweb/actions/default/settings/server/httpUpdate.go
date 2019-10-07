package server

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/settings"
	"github.com/iwind/TeaGo/actions"
	"net"
	"strings"
)

type HttpUpdateAction actions.Action

func (this *HttpUpdateAction) Run(params struct {
	On        bool
	Addresses string
	Must      *actions.Must
}) {
	params.Must.
		Field("addresses", params.Addresses).
		Require("请输入绑定地址")

	server, err := teaconfigs.LoadWebConfig()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	server.Http.On = params.On

	listen := []string{}
	for _, addr := range strings.Split(params.Addresses, "\n") {
		addr = strings.TrimSpace(addr)
		if len(addr) == 0 {
			continue
		}
		if _, _, err := net.SplitHostPort(addr); err != nil {
			addr += ":80"
		}
		listen = append(listen, addr)
	}
	server.Http.Listen = listen

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	settings.NotifyServerChange()

	this.Next("/settings", nil).
		Success("保存成功，重启服务后生效")
}
