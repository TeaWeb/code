package monitor

import (
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/teadb"
	"github.com/TeaWeb/code/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"runtime"
)

type IndexAction actions.Action

// 监控信息
func (this *IndexAction) Run(params struct{}) {
	apiutils.ValidateUser(this)

	result := maps.Map{
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
		"routines": runtime.NumGoroutine(),
		"version":  teaconst.TeaVersion,
	}

	stat := runtime.MemStats{}
	runtime.ReadMemStats(&stat)
	result["heap"] = stat.HeapAlloc
	result["memory"] = stat.Sys
	result["mongo"] = teadb.SharedDB().Test() == nil

	apiutils.Success(this, result)
}
