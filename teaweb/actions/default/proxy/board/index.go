package board

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/board/scripts"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"io/ioutil"
)

type IndexAction actions.Action

// 面板
func (this *IndexAction) Run(params struct {
	Server string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail("找不到要查看的代理服务")
	}

	this.Data["server"] = maps.Map{
		"id":       server.Id,
		"filename": params.Server,
	}

	configFile := "board." + server.Id + ".conf"
	if !files.NewFile(Tea.ConfigFile(configFile)).Exists() {
		configFile = "board.default.conf"

		// 如果配置文件不存在，则尝试创建
		if !files.NewFile(Tea.ConfigFile(configFile)).Exists() {
			err := ioutil.WriteFile(Tea.ConfigFile(configFile), []byte(`widgets:
- id: "1545562554961080824"
  code: "proxy_status@tea"
- id: "1545562554961080825"
  code: "locations@tea"
- id: "1545562554961080826"
  code: "bandwidth_realtime@tea"
- id: "1545562554961080827"
  code: "request_realtime@tea"
- id: "1545562554961080828"
  code: "request_time@tea"
- id: "1545562554961080829"
  code: "status_stat@tea"
- id: "1545562554961080830"
  code: "latest_error_log@tea"`), 0777)
			if err != nil {
				logs.Println("failed to create '" + configFile + "'")
			}
		}
	}

	this.Data["config"] = configFile

	this.Show()
}

// 面板数据
func (this *IndexAction) RunPost(params struct {
	Server string
	Config string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail("找不到要查看的代理服务")
	}

	engine := scripts.NewEngine()
	engine.SetContext(&scripts.Context{
		Server: server,
	})
	err = engine.RunConfig(Tea.ConfigFile(params.Config), maps.Map{})
	this.Data["widgetError"] = ""
	if err != nil {
		this.Data["widgetError"] = err.Error()
	}
	this.Data["charts"] = engine.Charts()
	this.Success()
}
