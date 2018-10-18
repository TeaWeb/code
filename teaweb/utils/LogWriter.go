package utils

import (
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
		dir.Mkdir()
	}

	// 先删除原来的
	logFile := files.NewFile(Tea.LogFile("teaweb.log"))
	if logFile.Exists() {
		logFile.Delete()
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
		this.fileAppender.AppendString(timeutil.Format("Y/m/d H:i:s ") + message + "\n")
	}
}

func (this *LogWriter) Close() {
	if this.fileAppender != nil {
		this.fileAppender.Close()
	}
}
