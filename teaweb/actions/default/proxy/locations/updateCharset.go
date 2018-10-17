package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/iwind/TeaGo/actions"
)

type UpdateCharsetAction actions.Action

func (this *UpdateCharsetAction) Run(params struct {
	Filename string
	Index    int
	Charset  string
	Must     *actions.Must
}) {
	params.Must.
		Field("charset", params.Charset).
		Require("请选择字符集")

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	location := proxy.LocationAtIndex(params.Index)
	if location != nil {
		location.Charset = params.Charset
		proxy.WriteBack()

		global.NotifyChange()
	}

	this.Success()
}
