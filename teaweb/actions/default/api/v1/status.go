package v1

import (
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"os"
	"runtime"
)

type StatusAction actions.Action

// 状态
func (this *StatusAction) RunGet(params struct{}) {
	result := maps.Map{
		"pid":      os.Getpid(),
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
		"go":       runtime.Version(),
		"routines": runtime.NumGoroutine(),
		"version":  teaconst.TeaVersion,
	}

	stat := runtime.MemStats{}
	runtime.ReadMemStats(&stat)
	result["heap"] = stat.HeapAlloc
	result["memory"] = stat.Sys
	result["objects"] = stat.HeapObjects
	result["mongo"] = teamongo.Test() == nil

	apiutils.Success(this, result)
}
