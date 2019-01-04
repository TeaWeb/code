package locations

import (
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type FastcgiAction actions.Action

// Fastcgi设置
func (this *FastcgiAction) Run(params struct {
	Server     string
	LocationId string
}) {
	locationutils.SetCommonInfo(this, params.Server, params.LocationId, "fastcgi")

	this.Data["queryParams"] = maps.Map{
		"server":     params.Server,
		"locationId": params.LocationId,
	}

	this.Show()
}
