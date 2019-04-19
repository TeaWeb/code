package backend

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
	"regexp"
)

type AddAction actions.Action

// 添加服务器
func (this *AddAction) Run(params struct {
	From       string
	ServerId   string
	LocationId string // 路径
	Websocket  bool   // 是否是Websocket设置
	Backup     bool
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if len(params.LocationId) > 0 {
		this.Data["selectedTab"] = "location"
	} else {
		this.Data["selectedTab"] = "backend"
	}
	this.Data["server"] = server

	this.Data["from"] = params.From
	this.Data["locationId"] = params.LocationId
	this.Data["websocket"] = types.Int(params.Websocket)
	this.Data["isBackup"] = params.Backup

	this.Show()
}

// 提交
func (this *AddAction) RunPost(params struct {
	ServerId        string
	LocationId      string // 路径
	Websocket       bool   // 是否是Websocket设置
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

	backend := teaconfigs.NewBackendConfig()
	backend.Address = params.Address
	backend.Scheme = params.Scheme
	backend.Weight = params.Weight
	backend.RequestGroupIds = params.RequestGroupIds
	backend.On = params.On
	backend.IsDown = false
	backend.Code = params.Code
	backend.FailTimeout = fmt.Sprintf("%d", params.FailTimeout) + "s"
	backend.ReadTimeout = fmt.Sprintf("%d", params.ReadTimeout) + "s"
	backend.MaxFails = params.MaxFails
	backend.MaxConns = params.MaxConns
	backend.IsBackup = params.IsBackup
	backend.RequestURI = params.RequestURI
	backend.CheckURL = params.CheckURL
	backend.CheckInterval = params.CheckInterval

	// 请求Header
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

	backendList, err := server.FindBackendList(params.LocationId, params.Websocket)
	if err != nil {
		this.Fail(err.Error())
	}

	backendList.AddBackend(backend)

	err = server.Save()
	if err != nil {
		this.Fail(err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
