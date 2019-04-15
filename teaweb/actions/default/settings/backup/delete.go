package backup

import (
	"github.com/TeaWeb/code/teaconst"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
)

type DeleteAction actions.Action

// 删除备份
func (this *DeleteAction) Run(params struct {
	File string
}) {
	if teaconst.DemoEnabled {
		this.Fail("演示版无法删除")
	}

	if len(params.File) == 0 {
		this.Fail("请指定要删除的备份文件")
	}

	file := files.NewFile(Tea.Root + "/backups/" + params.File)
	if file.Exists() {
		err := file.Delete()
		if err != nil {
			this.Fail("删除失败：" + err.Error())
		}
	}

	this.Success()
}
