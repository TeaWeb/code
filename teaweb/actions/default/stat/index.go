package stat

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/TeaWeb/code/teastats"
)

type IndexAction actions.Action

func (this *IndexAction) Run(params struct{}) {
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
