package locations

import (
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/iwind/TeaGo/actions"
)

type AccessAction actions.Action

// 访问控制
func (this *AccessAction) Run(params struct {
	Server     string
	LocationId string
}) {
	_, location := locationutils.SetCommonInfo(this, params.Server, params.LocationId, "access")

	this.Data["policy"] = location.AccessPolicy

	this.Show()
}
