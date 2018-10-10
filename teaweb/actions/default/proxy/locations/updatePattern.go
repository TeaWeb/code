package locations

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
)

type UpdatePatternAction actions.Action

func (this *UpdatePatternAction) Run(params struct {
	Filename string
	Index    int
	Pattern  string
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	location := proxy.LocationAtIndex(params.Index)
	if location != nil {
		location.SetPattern(params.Pattern, location.PatternType(), location.IsCaseInsensitive(), location.IsReverse())
		proxy.WriteToFilename(params.Filename)
	}

	global.NotifyChange()

	this.Refresh().Success()
}
