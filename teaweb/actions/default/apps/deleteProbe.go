package apps

import (
	"github.com/TeaWeb/jsapps/probes"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
)

type DeleteProbeAction actions.Action

// 删除探针
func (this *DeleteProbeAction) Run(params struct {
	ProbeId string
}) {
	if len(params.ProbeId) == 0 {
		this.Fail("请选择要删除的探针")
	}

	parser := probes.NewParser(Tea.ConfigFile("jsapps.js"))
	err := parser.RemoveProbe(params.ProbeId)
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	this.Success()
}
