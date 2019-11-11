package utils

import (
	"github.com/TeaWeb/code/teaconst"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/utils/time"
	"log"
)

type LogWriter struct {
	fileAppender *files.Appender
}

func (this *LogWriter) Init() {
	// 创建目录
	dir := files.NewFile(Tea.LogDir())
	if !dir.Exists() {
		err := dir.Mkdir()
		if err != nil {
			log.Println("[error]" + err.Error())
		}
	}

	// 先删除原来的
	logFile := files.NewFile(Tea.LogFile(teaconst.TeaProcessName + ".log"))
	if logFile.Exists() {
		err := logFile.Delete()
		if err != nil {
			log.Println("[error]" + err.Error())
		}
	}

	// 打开要写入的日志文件
	appender, err := logFile.Appender()
	if err != nil {
		logs.Error(err)
	} else {
		this.fileAppender = appender
	}
}

func (this *LogWriter) Write(message string) {
	log.Println(message)

	if this.fileAppender != nil {
		_, err := this.fileAppender.AppendString(timeutil.Format("Y/m/d H:i:s ") + message + "\n")
		if err != nil {
			log.Println("[error]" + err.Error())
		}
	}
}

func (this *LogWriter) Close() {
	if this.fileAppender != nil {
		_ = this.fileAppender.Close()
	}
}
