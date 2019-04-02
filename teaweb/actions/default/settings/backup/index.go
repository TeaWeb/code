package backup

import (
	"fmt"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"regexp"
	"time"
)

type IndexAction actions.Action

// 是否需要重启
var shouldRestart = false

// 备份列表
func (this *IndexAction) Run(params struct{}) {
	// 已备份
	result := []maps.Map{}

	reg := regexp.MustCompile("^\\d{8}\\.zip$")
	dir := files.NewFile(Tea.Root + "/backups/")
	if dir.Exists() {
		for _, f := range dir.List() {
			if !reg.MatchString(f.Name()) {
				continue
			}
			modifiedTime, _ := f.LastModified()
			size, _ := f.Size()
			result = append(result, maps.Map{
				"name":        f.Name(),
				"time":        timeutil.Format("Y-m-d H:i:s", modifiedTime),
				"size":        fmt.Sprintf("%.2f", float64(size)/1024/1024), // M
				"isToday":     timeutil.Format("Ymd")+".zip" == f.Name(),
				"isYesterday": timeutil.Format("Ymd", time.Now().Add(-24*time.Hour))+".zip" == f.Name(),
			})
		}
	}

	lists.Sort(result, func(i int, j int) bool {
		return result[i].GetString("name") > result[j].GetString("name")
	})

	this.Data["files"] = result
	this.Data["shouldRestart"] = shouldRestart

	this.Show()
}
