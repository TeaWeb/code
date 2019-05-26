package server

import (
	"github.com/TeaWeb/code/teaweb/actions/default/settings"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
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

	reader, err := files.NewReader(Tea.ConfigFile("server.conf"))
	if err != nil {
		this.Fail("无法读取配置文件（'configs/server.conf'），请检查文件是否存在，或者是否有权限读取")
	}
	defer reader.Close()

	server := &TeaGo.ServerConfig{}
	err = reader.ReadYAML(server)
	if err != nil {
		this.Fail("配置文件（'configs/server.conf'）格式错误")
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

	writer, err := files.NewWriter(Tea.ConfigFile("server.conf"))
	if err != nil {
		this.Fail("配置文件（'configs/server.conf'）打开失败")
	}
	defer writer.Close()

	writer.WriteYAML(server)

	settings.NotifyServerChange()

	this.Next("/settings", nil).Success("保存成功，重启服务后生效")
}
