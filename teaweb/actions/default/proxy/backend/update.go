package backend

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"regexp"
)

type UpdateAction actions.Action

// 修改后端服务器
func (this *UpdateAction) Run(params struct {
	ServerId   string
	LocationId string
	Websocket  bool
	Backend    string
	From       string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["server"] = server
	if len(params.LocationId) > 0 {
		this.Data["selectedTab"] = "location"
	} else {
		this.Data["selectedTab"] = "backend"
	}
	this.Data["locationId"] = params.LocationId
	this.Data["websocket"] = types.Int(params.Websocket)
	this.Data["from"] = params.From

	backendList, err := server.FindBackendList(params.LocationId, params.Websocket)
	if err != nil {
		this.Fail(err.Error())
	}
	backend := backendList.FindBackend(params.Backend)
	if backend == nil {
		this.Fail("找不到要修改的后端服务器")
	}

	backend.Validate()

	if len(backend.RequestGroupIds) == 0 {
		backend.AddRequestGroupId("default")
	}

	if len(backend.RequestURI) == 0 {
		backend.RequestURI = "${requestURI}"
	}

	this.Data["backend"] = maps.Map{
		"id":              backend.Id,
		"address":         backend.Address,
		"scheme":          backend.Scheme,
		"code":            backend.Code,
		"weight":          backend.Weight,
		"failTimeout":     int(backend.FailTimeoutDuration().Seconds()),
		"readTimeout":     int(backend.ReadTimeoutDuration().Seconds()),
		"on":              backend.On,
		"maxConns":        backend.MaxConns,
		"maxFails":        backend.MaxFails,
		"isDown":          backend.IsDown,
		"isBackup":        backend.IsBackup,
		"requestGroupIds": backend.RequestGroupIds,
		"requestURI":      backend.RequestURI,
		"checkURL":        backend.CheckURL,
		"checkInterval":   backend.CheckInterval,
		"requestHeaders":  backend.RequestHeaders,
		"responseHeaders": backend.ResponseHeaders,
		"host":            backend.Host,
	}

	this.Show()
}

// 提交
func (this *UpdateAction) RunPost(params struct {
	ServerId        string
	LocationId      string
	Websocket       bool
	BackendId       string
	Address         string
	Scheme          string
	Weight          uint
	On              bool
	Code            string
	FailTimeout     uint
	ReadTimeout     uint
	MaxFails        int32
	MaxConns        int32
	IsBackup        bool
	RequestGroupIds []string
	RequestURI      string
	CheckURL        string
	CheckInterval   int

	RequestHeaderNames  []string
	RequestHeaderValues []string

	ResponseHeaderNames  []string
	ResponseHeaderValues []string

	Host string

	Must *actions.Must
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	params.Must.
		Field("address", params.Address).
		Require("请输入后端服务器地址")

	if len(params.CheckURL) > 0 {
		if !regexp.MustCompile("(?i)(http://|https://)").MatchString(params.CheckURL) {
			this.FailField("checkURL", "健康检查URL必须以http://或https://开头")
		}
	}

	backendList, err := server.FindBackendList(params.LocationId, params.Websocket)
	if err != nil {
		this.Fail(err.Error())
	}

	backend := backendList.FindBackend(params.BackendId)
	if backend == nil {
		this.Fail("找不到要修改的后端服务器")
	}

	backend.Address = params.Address
	backend.Scheme = params.Scheme
	backend.Weight = params.Weight
	backend.On = params.On
	backend.IsDown = false
	backend.Code = params.Code
	backend.FailTimeout = fmt.Sprintf("%d", params.FailTimeout) + "s"
	backend.ReadTimeout = fmt.Sprintf("%d", params.ReadTimeout) + "s"
	backend.MaxFails = params.MaxFails
	backend.MaxConns = params.MaxConns
	backend.IsBackup = params.IsBackup
	backend.RequestGroupIds = params.RequestGroupIds
	backend.RequestURI = params.RequestURI
	backend.CheckURL = params.CheckURL
	backend.CheckInterval = params.CheckInterval

	// 请求Header
	backend.RequestHeaders = []*shared.HeaderConfig{}
	if len(params.RequestHeaderNames) > 0 {
		for index, headerName := range params.RequestHeaderNames {
			if index < len(params.RequestHeaderValues) {
				header := shared.NewHeaderConfig()
				header.Name = headerName
				header.Value = params.RequestHeaderValues[index]
				backend.AddRequestHeader(header)
			}
		}
	}

	// 响应Header
	backend.ResponseHeaders = []*shared.HeaderConfig{}
	if len(params.ResponseHeaderNames) > 0 {
		for index, headerName := range params.ResponseHeaderNames {
			if index < len(params.ResponseHeaderValues) {
				header := shared.NewHeaderConfig()
				header.Name = headerName
				header.Value = params.ResponseHeaderValues[index]
				backend.AddResponseHeader(header)
			}
		}
	}

	backend.Host = params.Host

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
