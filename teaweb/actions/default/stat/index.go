package stat

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teastats"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct {
	ServerId string
}) {
	// 检查MongoDB连接
	this.Data["mongoError"] = ""
	err := teamongo.Test()
	if err != nil {
		this.Data["mongoError"] = "此功能需要连接MongoDB"
		this.Data["topRequests"] = []bool{}
		this.Data["topCostRequests"] = []bool{}
		this.Data["topOS"] = []bool{}
		this.Data["topBrowsers"] = []bool{}
		this.Data["topRegions"] = []bool{}
		this.Data["topStates"] = []bool{}

		this.Show()
		return
	}

	// 代理列表
	servers := []maps.Map{}
	for index, config := range teaconfigs.LoadServerConfigsFromDir(Tea.ConfigDir()) {
		if len(params.ServerId) == 0 {
			params.ServerId = config.Id
		}

		servers = append(servers, maps.Map{
			"name":    config.Description,
			"subName": "",
			"active":  params.ServerId == config.Id || (len(params.ServerId) == 0 && index == 0),
			"url":     "/stat?serverId=" + config.Id,
		})
	}
	if len(servers) > 0 {
		this.Data["teaTabbar"] = servers
	}

	this.Data["serverId"] = params.ServerId

	// 访问量排行
	{
		stat := new(teastats.TopRequestStat)
		this.Data["topRequests"] = stat.List(params.ServerId, 10)
	}

	// 请求耗时排行
	{
		stat := new(teastats.TopCostStat)
		this.Data["topCostRequests"] = stat.List(params.ServerId, 10)
	}

	// 操作系统排行
	{
		stat := new(teastats.TopOSStat)
		this.Data["topOS"] = stat.List(params.ServerId, 10)
	}

	// 浏览器排行
	{
		stat := new(teastats.TopBrowserStat)
		this.Data["topBrowsers"] = stat.List(params.ServerId, 10)
	}

	// 地区排行
	{
		stat := new(teastats.TopRegionStat)
		this.Data["topRegions"] = stat.List(params.ServerId, 10)
	}

	// 省份排行
	{
		stat := new(teastats.TopStateStat)
		this.Data["topStates"] = stat.List(params.ServerId, 10)
	}

	this.Show()
}
