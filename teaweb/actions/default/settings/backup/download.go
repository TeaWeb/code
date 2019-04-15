package backup

import (
	"github.com/TeaWeb/code/teaconst"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"io"
	"os"
	"regexp"
)

type DownloadAction actions.Action

// 下载
func (this *DownloadAction) Run(params struct {
	Filename string
}) {
	if teaconst.DemoEnabled {
		this.Fail("演示版无法下载")
	}

	reg := regexp.MustCompile("^\\d{8}\\.zip$")
	if !reg.MatchString(params.Filename) {
		this.WriteString("file not found")
		return
	}

	fp, err := os.Open(Tea.Root + "/backups/" + params.Filename)
	if err != nil {
		this.WriteString("file not found")
		return
	}
	defer fp.Close()

	this.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename=\""+params.Filename+"\"")
	_, err = io.Copy(this.ResponseWriter, fp)
	if err != nil {
		logs.Error(err)
		return
	}
}
