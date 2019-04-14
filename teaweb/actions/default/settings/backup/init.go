package backup

import (
	"archive/zip"
	"errors"
	"github.com/TeaWeb/code/teaweb/actions/default/settings"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/TeaWeb/code/teaweb/helpers"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/utils/time"
	"os"
	"strings"
	"time"
)

func init() {
	// 路由设置
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(&helpers.UserMustAuth{
				Grant: configs.AdminGrantAll,
			}).
			Helper(new(settings.Helper)).
			Prefix("/settings/backup").
			Get("", new(IndexAction)).
			Post("/backup", new(BackupAction)).
			Post("/delete", new(DeleteAction)).
			Post("/restore", new(RestoreAction)).
			Get("/download", new(DownloadAction)).
			EndAll()
	})

	// 自动备份
	backup()
}

// 自动备份
func backup() {
	timers.Every(24*time.Hour, func(ticker *time.Ticker) {
		err := backupTask()
		if err != nil {
			logs.Error(err)
		}
	})
}

func backupTask() error {
	dir := files.NewFile(Tea.Root + "/backups/")
	if !dir.Exists() {
		err := dir.Mkdir()
		if err != nil {
			return err
		}
		if !dir.Exists() {
			return errors.New("'backups/' not exists")
		}
	}

	logFile := Tea.Root + "/backups/" + timeutil.Format("Ymd") + ".zip"
	fp, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer fp.Close()

	z := zip.NewWriter(fp)
	defer z.Close()

	configsDir := files.NewFile(Tea.Root + "/configs")
	configsAbs, _ := configsDir.AbsPath()
	configsAbs += Tea.DS
	configsDir.Range(func(file *files.File) {
		if !file.IsFile() {
			return
		}

		// 脚本不保存，运行时会自动生成
		if strings.HasSuffix(file.Name(), ".script") {
			return
		}

		modified, err := file.LastModified()
		if err != nil {
			modified = time.Now()
		}
		path, _ := file.AbsPath()
		h := &zip.FileHeader{
			Name:     "configs" + Tea.DS + strings.TrimPrefix(path, configsAbs),
			Modified: modified,
		}
		w, err := z.CreateHeader(h)
		if err != nil {
			logs.Error(err)
			return
		}
		data, err := file.ReadAll()
		if err != nil {
			logs.Error(err)
			return
		}
		_, err = w.Write(data)
		if err != nil {
			logs.Error(err)
			return
		}
	})
	return nil
}
