package groups

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type UpdateAction actions.Action

// 修改分组
func (this *UpdateAction) Run(params struct {
	ServerId   string
	LocationId string
	Websocket  int
	GroupId    string
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
	this.Data["locationId"] = params.LocationId
	this.Data["websocket"] = params.Websocket

	group := server.FindRequestGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到Group")
	}

	this.Data["group"] = group
	this.Data["operators"] = teaconfigs.AllRequestOperators()

	// 请求变量
	this.Data["variables"] = proxyutils.DefaultRequestVariables()

	this.Show()
}

// 提交保存
func (this *UpdateAction) RunPost(params struct {
	ServerId string
	GroupId  string

	Name string

	CondParams []string
	CondOps    []string
	CondValues []string

	IPRangeTypeList     []string `alias:"ipRangeTypeList"`
	IPRangeFromList     []string `alias:"ipRangeFromList"`
	IPRangeToList       []string `alias:"ipRangeToList"`
	IPRangeCIDRIPList   []string `alias:"ipRangeCIDRIPList"`
	IPRangeCIDRBitsList []string `alias:"ipRangeCIDRBitsList"`
	IPRangeVarList      []string `alias:"ipRangeVarList"`

	RequestHeaderNames  []string
	RequestHeaderValues []string

	ResponseHeaderNames  []string
	ResponseHeaderValues []string

	Must *actions.Must
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到要修改的Server")
	}

	group := server.FindRequestGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到Group")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入分组名")

	group.Name = params.Name
	group.Cond = []*teaconfigs.RequestCond{}
	group.IPRanges = []*teaconfigs.IPRangeConfig{}
	group.RequestHeaders = []*shared.HeaderConfig{}
	group.ResponseHeaders = []*shared.HeaderConfig{}

	// 匹配条件
	if len(params.CondParams) > 0 {
		for index, param := range params.CondParams {
			if index < len(params.CondOps) && index < len(params.CondValues) {
				cond := teaconfigs.NewRequestCond()
				cond.Param = param
				cond.Value = params.CondValues[index]
				cond.Operator = params.CondOps[index]
				err := cond.Validate()
				if err != nil {
					this.Fail("匹配条件\"" + cond.Param + " " + cond.Value + "\"校验失败：" + err.Error())
				}
				group.AddCond(cond)
			}
		}
	}

	// IP范围
	if len(params.IPRangeTypeList) > 0 {
		for index, ipRangeType := range params.IPRangeTypeList {
			if index < len(params.IPRangeFromList) && index < len(params.IPRangeToList) && index < len(params.IPRangeCIDRIPList) && index < len(params.IPRangeCIDRBitsList) {
				if ipRangeType == "range" {
					config := teaconfigs.NewIPRangeConfig()
					config.Type = teaconfigs.IPRangeTypeRange
					config.IPFrom = params.IPRangeFromList[index]
					config.IPTo = params.IPRangeToList[index]
					config.Param = params.IPRangeVarList[index]
					err := config.Validate()
					if err != nil {
						this.Fail("校验失败：" + err.Error())
					}
					group.AddIPRange(config)
				} else if ipRangeType == "cidr" {
					config := teaconfigs.NewIPRangeConfig()
					config.Type = teaconfigs.IPRangeTypeCIDR
					config.CIDR = params.IPRangeCIDRIPList[index] + "/" + params.IPRangeCIDRBitsList[index]
					config.Param = params.IPRangeVarList[index]
					err := config.Validate()
					if err != nil {
						this.Fail("校验失败：" + err.Error())
					}
					group.AddIPRange(config)
				}
			}
		}
	}

	// 请求Header
	if len(params.RequestHeaderNames) > 0 {
		for index, headerName := range params.RequestHeaderNames {
			if index < len(params.RequestHeaderValues) {
				header := shared.NewHeaderConfig()
				header.Name = headerName
				header.Value = params.RequestHeaderValues[index]
				group.AddRequestHeader(header)
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
				group.AddResponseHeader(header)
			}
		}
	}

	// 保存
	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知改变
	proxyutils.NotifyChange()

	this.Success()
}
