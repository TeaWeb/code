package proxy

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
)

// 添加新的服务
type AddAction actions.Action

func (this *AddAction) Run(params struct {
}) {
	this.Show()
}

func (this *AddAction) RunPost(params struct {
	Description string
	Must        *actions.Must
}) {
	if len(params.Description) == 0 {
		params.Description = "新服务"
	}

	server := teaconfigs.NewServerConfig()
	server.Http = true
	server.Description = params.Description
	server.Charset = "utf-8"
	server.Index = []string{"index.html", "index.htm", "index.php"}

	filename := stringutil.Rand(16) + ".proxy.conf"
	configPath := Tea.ConfigFile(filename)
	err := server.WriteToFile(configPath)
	if err != nil {
		this.Fail(err.Error())
	}

	global.NotifyChange()

	this.Next("/proxy/detail", map[string]interface{}{
		"filename": filename,
	}, "").Success("添加成功，现在去配置详细信息")
}
