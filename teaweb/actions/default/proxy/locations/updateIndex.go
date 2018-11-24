package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaGo/actions"
	"regexp"
)

type UpdateIndexAction actions.Action

func (this *UpdateIndexAction) Run(params struct {
	Filename string
	Index    int
	Indexes  string
	Must     *actions.Must
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	location := proxy.LocationAtIndex(params.Index)
	if location == nil {
		this.Fail("找不到要修改的路径规则")
	}

	if len(params.Indexes) > 0 {
		indexes := regexp.MustCompile("\\s+").Split(params.Indexes, -1)
		location.Index = indexes
	} else {
		location.Index = []string{}
	}

	proxy.Save()

	global.NotifyChange()

	this.Data["indexes"] = location.Index

	this.Success()
}
