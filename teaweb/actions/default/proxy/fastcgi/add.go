package fastcgi

import (
	"github.com/iwind/TeaGo/actions"
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/global"
	"github.com/TeaWeb/code/teaconfigs"
)

type AddAction actions.Action

func (this *AddAction) Run(params struct {
	Filename    string
	Index       int
	On          bool
	Pass        string
	ReadTimeout int
	Params      string
	Must        *actions.Must
}) {
	params.Must.
		Field("filename", params.Filename).
		Require("请输入配置文件名").
		Field("pass", params.Pass).
		Require("请输入Fastcgi地址")

	paramsMap := map[string]string{}
	err := ffjson.Unmarshal([]byte(params.Params), &paramsMap)
	if err != nil {
		this.Fail(err.Error())
	}

	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	location := proxy.LocationAtIndex(params.Index)
	if location == nil {
		this.Fail("找不到要修改的路径规则")
	}

	fastcgi := teaconfigs.NewFastcgiConfig()
	fastcgi.On = params.On
	fastcgi.Pass = params.Pass
	fastcgi.ReadTimeout = fmt.Sprintf("%ds", params.ReadTimeout)
	fastcgi.Params = paramsMap
	location.AddFastcgi(fastcgi)
	proxy.WriteBack()

	global.NotifyChange()

	this.Refresh().Success()
}
