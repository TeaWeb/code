package proxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
)

type UpdateDescriptionAction actions.Action

func (this *UpdateDescriptionAction) Run(params struct {
	Filename    string
	Description string
	Must        *actions.Must
}) {
	params.Must.
		Field("description", params.Description).
		Require("代理说明不能为空")

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy.Description = params.Description
	proxy.WriteToFile(Tea.ConfigFile(params.Filename))

	this.Refresh().Success("保存成功")
}
