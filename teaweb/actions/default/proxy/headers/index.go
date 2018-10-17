package headers

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct {
	Filename string
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["selectedTab"] = "header"
	this.Data["filename"] = params.Filename
	this.Data["proxy"] = proxy

	// headers
	this.Data["headers"] = proxy.Headers
	this.Data["ignoreHeaders"] = lists.NewList(proxy.IgnoreHeaders).Map(func(k int, v interface{}) interface{} {
		return map[string]interface{}{
			"name": v,
		}
	}).Slice

	this.Show()
}
