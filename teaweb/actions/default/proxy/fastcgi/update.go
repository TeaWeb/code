package fastcgi

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction actions.Action

// 修改
func (this *UpdateAction) Run(params struct {
	From       string
	Server     string
	LocationId string
	FastcgiId  string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}
	fastcgiList, err := server.FindFastcgiList(params.LocationId)
	if err != nil {
		this.Fail(err.Error())
	}
	fastcgi := fastcgiList.FindFastcgi(params.FastcgiId)
	if fastcgi == nil {
		this.Fail("找不到要修改的Fastcgi")
	}

	m := maps.Map{
		"on":       fastcgi.On,
		"id":       fastcgi.Id,
		"pass":     fastcgi.Pass,
		"poolSize": fastcgi.PoolSize,
		"params":   fastcgi.Params,
	}
	if fastcgi.ReadTimeout != "0s" {
		m["readTimeoutSeconds"] = int(fastcgi.ReadTimeoutDuration().Seconds())
	} else {
		m["readTimeoutSeconds"] = 0
	}
	this.Data["fastcgi"] = m

	this.Data["from"] = params.From
	this.Data["server"] = maps.Map{
		"filename": params.Server,
	}
	this.Data["locationId"] = params.LocationId

	this.Show()
}

// 修改
func (this *UpdateAction) RunPost(params struct {
	Server      string
	LocationId  string
	On          bool
	Pass        string
	ReadTimeout int
	ParamNames  []string
	ParamValues []string
	PoolSize    int

	FastcgiId string

	Must *actions.Must
}) {
	params.Must.
		Field("pass", params.Pass).
		Require("请输入Fastcgi地址").
		Field("poolSize", params.PoolSize).
		Gte(0, "连接池尺寸不能小于0")

	paramsMap := map[string]string{}
	for index, paramName := range params.ParamNames {
		if index < len(params.ParamValues) {
			paramsMap[paramName] = params.ParamValues[index]
		}
	}

	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	fastcgiList, err := server.FindFastcgiList(params.LocationId)
	if err != nil {
		this.Fail(err.Error())
	}

	fastcgi := fastcgiList.FindFastcgi(params.FastcgiId)
	if fastcgi == nil {
		this.Fail("找不到要修改的Fastcgi")
	}

	fastcgi.On = params.On
	fastcgi.Pass = params.Pass
	fastcgi.ReadTimeout = fmt.Sprintf("%ds", params.ReadTimeout)
	fastcgi.Params = paramsMap
	fastcgi.PoolSize = params.PoolSize
	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
