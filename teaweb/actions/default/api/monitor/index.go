package monitor

import (
	"encoding/json"
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/teamongo"
	"github.com/TeaWeb/code/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"runtime"
)

type IndexAction actions.Action

// 监控信息
func (this *IndexAction) Run(params struct{}) {
	apiutils.ValidateUser(this)

	this.AddHeader("Content-Type", "application/json; charset=utf-8")

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
	result["mongo"] = teamongo.Test() == nil

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logs.Error(err)
		this.Error(err.Error(), http.StatusInternalServerError)
	} else {
		this.Write(data)
	}
}
