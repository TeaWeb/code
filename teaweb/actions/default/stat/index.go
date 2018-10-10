package stat

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teastats"
	"github.com/TeaWeb/code/teamongo"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct{}) {
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

	// 访问量排行
	{
		stat := new(teastats.TopRequestStat)
		this.Data["topRequests"] = stat.List(10)
	}

	// 请求耗时排行
	{
		stat := new(teastats.TopCostStat)
		this.Data["topCostRequests"] = stat.List(10)
	}

	// 操作系统排行
	{
		stat := new(teastats.TopOSStat)
		this.Data["topOS"] = stat.List(10)
	}

	// 浏览器排行
	{
		stat := new(teastats.TopBrowserStat)
		this.Data["topBrowsers"] = stat.List(10)
	}

	// 地区排行
	{
		stat := new(teastats.TopRegionStat)
		this.Data["topRegions"] = stat.List(10)
	}

	// 省份排行
	{
		stat := new(teastats.TopStateStat)
		this.Data["topStates"] = stat.List(10)
	}

	this.Show()
}
