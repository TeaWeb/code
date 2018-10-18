package fastcgi

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
)

type UpdatePassAction actions.Action

func (this *UpdatePassAction) Run(params struct {
	Filename string
	Index    int
	Pass     string
	Must     *actions.Must
}) {
	params.Must.
		Field("filename", params.Filename).
		Require("请输入配置文件名").
		Field("pass", params.Pass).
		Require("请输入Fastcgi地址")

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	location := proxy.LocationAtIndex(params.Index)
	if location == nil {
		this.Fail("找不到要修改的路径规则")
	}

	fastcgi := location.FastcgiAtIndex(0)
	if fastcgi == nil {
		this.Fail("没有fastcgi配置，请刷新后重试")
	}
	fastcgi.Pass = params.Pass
	proxy.WriteBack()

	global.NotifyChange()

	this.Refresh().Success()
}
