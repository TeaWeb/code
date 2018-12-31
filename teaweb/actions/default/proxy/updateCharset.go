package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type UpdateCharsetAction actions.Action

func (this *UpdateCharsetAction) Run(params struct {
	Filename string
	Charset  string
}) {
	if len(params.Charset) == 0 {
		this.Fail("请选择正确的字符编码")
	}

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy.Charset = params.Charset
	err = proxy.Save()
	if err != nil {
		this.Fail(err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
