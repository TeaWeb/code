package apps

import (
	"github.com/TeaWeb/jsapps/probes"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
)

type CopyProbeAction actions.Action

// 拷贝探针到本地
func (this *CopyProbeAction) Run(params struct {
	ProbeId string
}) {
	if len(params.ProbeId) == 0 {
		this.Fail("找不到要使用的探针")
	}

	var probeFunc = ""
	{
		parser := probes.NewParser(Tea.Root + Tea.DS + "plugins" + Tea.DS + "jsapps.js")
		s, err := parser.FindProbeFunction(params.ProbeId)
		if err != nil {
			this.Fail("发生错误：" + err.Error())
		}
		probeFunc = s
	}

	{
		parser := probes.NewParser(Tea.ConfigFile("jsapps.js"))
		err := parser.AddProbeFunc(probeFunc)
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}
	}

	this.Success()
}
