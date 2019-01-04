package locations

import (
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type HeadersAction actions.Action

// 自定义Http Header
func (this *HeadersAction) Run(params struct {
	Server     string // 必填
	LocationId string
}) {
	locationutils.SetCommonInfo(this, params.Server, params.LocationId, "headers")

	this.Data["headerQuery"] = maps.Map{
		"server":     params.Server,
		"locationId": params.LocationId,
	}

	this.Show()
}
