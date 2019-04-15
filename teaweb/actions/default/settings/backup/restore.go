package backup

import (
	"archive/zip"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/utils/time"
	"os"
	"path/filepath"
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

	file := files.NewFile(Tea.Root + "/backups/" + params.File)
	if !file.Exists() {
		this.Fail("指定的备份文件不存在")
	}

	// 解压
	reader, err := zip.OpenReader(file.Path())
	if err != nil {
		this.Fail("无法读取：" + err.Error())
	}
	defer reader.Close()

	// 清除backup configs
	tmpDir := files.NewFile(Tea.Root + "/backups/configs")
	if tmpDir.Exists() {
		err := tmpDir.DeleteAll()
		if err != nil {
			this.Fail("无法清除backups/configs")
		}
	}

	for _, entry := range reader.File {
		dir := filepath.Dir(entry.Name)
		target := files.NewFile(Tea.Root + "/backups/" + dir)
		if !target.Exists() {
			err := target.MkdirAll()
			if err != nil {
				this.Fail("创建目录失败：" + dir)
			}
		}
		reader, err := entry.Open()
		if err != nil {
			this.Fail("文件读取失败：" + err.Error())
		}
		data := []byte{}
		for {
			buf := make([]byte, 1024)
			n, err := reader.Read(buf)
			if n > 0 {
				data = append(data, buf[:n]...)
			}
			if err != nil {
				break
			}
		}
		err = files.NewFile(Tea.Root + "/backups/" + entry.Name).Write(data)
		if err != nil {
			reader.Close()
			this.Fail("文件写入失败：" + err.Error())
		}
		reader.Close()
	}

	// 修改老的配置文件
	oldDir := Tea.Root + "/old.configs." + timeutil.Format("YmdHis")
	err = os.Rename(Tea.ConfigDir(), oldDir)
	if err != nil {
		this.Fail("原配置清空失败：" + err.Error())
	}

	// 创建新的目录
	err = os.Rename(Tea.Root+"/backups/configs", Tea.ConfigDir())
	if err != nil {
		// 还原
		os.Rename(oldDir, Tea.ConfigDir())

		this.Fail("新配置拷贝失败")
	}

	teaproxy.SharedManager.Reload()
	agents.NotifyAgentsChange()

	shouldRestart = true

	this.Success()
}
