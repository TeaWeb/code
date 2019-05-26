package tealogs

import (
	"errors"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"os"
	"path/filepath"
	"sync"
)

type AccessLogFileWriter struct {
	config *teaconfigs.AccessLogFileConfig
	file   *os.File
	locker *sync.Mutex
}

func (writer *AccessLogFileWriter) Init() {
	if !filepath.IsAbs(writer.config.Path) {
		writer.config.Path = Tea.Root + Tea.DS + writer.config.Path
	}

	file, err := os.OpenFile(writer.config.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logs.Errorf("AccessLogFileConfig.Init(): %s", err)
		return
	}

	writer.file = file
	writer.locker = &sync.Mutex{}
}

func (writer *AccessLogFileWriter) Write(log *AccessLog) error {
	if len(writer.config.Path) == 0 {
		return errors.New("access log 'path' invalid")
	}

	format := writer.config.Format
	if len(format) == 0 {
		format = "${remoteAddr} - [${timeLocal}] \"${request}\" ${status} ${bodyBytesSent} \"${http.Referer}\" \"${http.UserAgent}\""
	}

	logString := log.Format(format)

	var lastErr error = nil
	if writer.file != nil {
		writer.locker.Lock()

		_, err := writer.file.WriteString(logString)
		if err != nil {
			lastErr = err
		}
		_, err = writer.file.WriteString("\n")
		if err != nil {
			lastErr = err
		}

		writer.locker.Unlock()
	}

	return lastErr
}

func (writer *AccessLogFileWriter) Close() {
	if writer.file != nil {
		writer.file.Close()
	}
}
