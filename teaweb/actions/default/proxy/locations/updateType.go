package locations

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
)

type UpdateTypeAction actions.Action

func (this *UpdateTypeAction) Run(params struct {
	Filename string
	Index    int
	TypeId   int
	Must     *actions.Must
}) {
	params.Must.
		Field("typeId", params.TypeId).
		Gt(0, "请选择类型").
		In([]int{
			teaconfigs.LocationPatternTypePrefix,
			teaconfigs.LocationPatternTypeExact,
			teaconfigs.LocationPatternTypeRegexp,
		}, "选择的类型错误")

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	location := proxy.LocationAtIndex(params.Index)
	if location != nil {
		location.SetPattern(location.PatternString(), params.TypeId, location.IsCaseInsensitive(), location.IsReverse())
		proxy.WriteToFilename(params.Filename)
	}

	global.NotifyChange()

	this.Refresh().Success()
}
