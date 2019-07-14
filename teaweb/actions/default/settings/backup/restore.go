package backup

import (
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/teaweb/actions/default/settings/backup/backuputils"
	"github.com/iwind/TeaGo/actions"
)

type RestoreAction actions.Action

// 从备份恢复
func (this *RestoreAction) Run(params struct {
	File string
}) {
	if teaconst.DemoEnabled {
		this.Fail("演示版无法恢复")
	}

	if len(params.File) == 0 {
		this.Fail("请指定要恢复的文件")
	}

	if !backuputils.RestoreFile(params.File, func(message string) {
		this.Fail(message)
	}) {
		return
	}

	this.Success()
}
