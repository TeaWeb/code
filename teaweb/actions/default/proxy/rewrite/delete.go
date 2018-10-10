package rewrite

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
)

type DeleteAction actions.Action

func (this *DeleteAction) Run(params struct {
	Filename     string
	Index        int
	RewriteIndex int
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	location := proxy.LocationAtIndex(params.Index)
	if location == nil {
		this.Fail("找不到要修改的路径规则")
	}

	if params.RewriteIndex >= 0 && params.RewriteIndex < len(location.Rewrite) {
		location.Rewrite = lists.Remove(location.Rewrite, params.RewriteIndex).([]*teaconfigs.RewriteRule)
	}

	proxy.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Refresh().Success()
}
