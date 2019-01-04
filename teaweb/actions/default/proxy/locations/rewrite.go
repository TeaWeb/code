package locations

import (
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type RewriteAction actions.Action

// 重写规则
func (this *RewriteAction) Run(params struct {
	Server     string
	LocationId string
}) {
	locationutils.SetCommonInfo(this, params.Server, params.LocationId, "rewrite")

	this.Data["queryParams"] = maps.Map{
		"server":     params.Server,
		"locationId": params.LocationId,
	}

	this.Show()
}
