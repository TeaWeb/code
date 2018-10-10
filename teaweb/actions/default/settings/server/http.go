package server

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo"
)

type HttpAction actions.Action

func (this *HttpAction) Run(params struct{}) {
	this.Data["error"] = ""

	reader, err := files.NewReader(Tea.ConfigFile("server.conf"))
	if err != nil {
		this.Data["error"] = "无法读取配置文件（'configs/server.conf'），请检查文件是否存在，或者是否有权限读取"
		this.Show()
		return
	}
	defer reader.Close()

	server := &TeaGo.ServerConfig{}
	err = reader.ReadYAML(server)
	if err != nil {
		this.Data["error"] = "配置文件（'configs/server.conf'）格式错误"
		this.Show()
		return
	}

	this.Data["server"] = server

	this.Show()
}
