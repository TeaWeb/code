package headers

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

func (this *DeleteAction) Run(params struct {
	Filename      string
	LocationIndex int
	Index         int
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	location := proxy.LocationAtIndex(params.LocationIndex)
	if location == nil {
		this.Fail("找不到要修改的路径规则")
	}

	location.DeleteHeaderAtIndex(params.Index)
	proxy.WriteBack()

	global.NotifyChange()

	this.Success()
}
