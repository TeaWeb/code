package fastcgi

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
)

type UpdateParamAction actions.Action

func (this *UpdateParamAction) Run(params struct {
	Filename string
	Index    int
	OldName  string
	Name     string
	Value    string
	Must     *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入参数名")

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

	delete(fastcgi.Params, params.OldName)
	fastcgi.Params[params.Name] = params.Value
	proxy.WriteBack()

	global.NotifyChange()

	this.Success()
}
