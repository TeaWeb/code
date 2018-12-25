package stat

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teastats"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// 统计
func (this *IndexAction) Run(params struct {
	Server string
}) {
	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail("找不到要查看的代理服务：" + err.Error())
	}
	serverId := server.Id

	// 检查MongoDB连接
	this.Data["mongoError"] = ""
	err = teamongo.Test()
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

	this.Data["serverId"] = serverId

	// 访问量排行
	{
		stat := new(teastats.TopRequestStat)
		this.Data["topRequests"] = stat.List(serverId, 10)
	}

	// 请求耗时排行
	{
		stat := new(teastats.TopCostStat)
		this.Data["topCostRequests"] = stat.List(serverId, 10)
	}

	// 操作系统排行
	{
		stat := new(teastats.TopOSStat)
		this.Data["topOS"] = stat.List(serverId, 10)
	}

	// 浏览器排行
	{
		stat := new(teastats.TopBrowserStat)
		this.Data["topBrowsers"] = stat.List(serverId, 10)
	}

	// 地区排行
	{
		stat := new(teastats.TopRegionStat)
		this.Data["topRegions"] = stat.List(serverId, 10)
	}

	// 省份排行
	{
		stat := new(teastats.TopStateStat)
		this.Data["topStates"] = stat.List(serverId, 10)
	}

	this.Show()
}
