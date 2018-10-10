package fastcgi

import (
	"github.com/iwind/TeaGo/actions"
	"fmt"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/TeaWeb/code/teaconfigs"
)

type UpdateTimeoutAction actions.Action

func (this *UpdateTimeoutAction) Run(params struct {
	Filename string
	Index    int
	Timeout  int
}) {
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

	fastcgi.ReadTimeout = fmt.Sprintf("%ds", params.Timeout)
	proxy.WriteToFilename(params.Filename)

	global.NotifyChange()

	this.Success()
}
