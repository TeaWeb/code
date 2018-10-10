package tealogs

import (
	"os"
	"sync"
	"github.com/TeaWeb/code/teaconfigs"
)

type AccessLogStdoutWriter struct {
	config *teaconfigs.AccessLogStdoutConfig
	locker *sync.Mutex
}

func (writer *AccessLogStdoutWriter) Init() {
	writer.locker = &sync.Mutex{}
}

func (writer *AccessLogStdoutWriter) Write(log *AccessLog) error {
	format := writer.config.Format
	if len(format) == 0 {
		format = "${remoteAddr} - [${timeLocal}] \"${request}\" ${status} ${bodyBytesSent} \"${http.Referer}\" \"${http.UserAgent}\""
	}

	logString := log.Format(format)

	writer.locker.Lock()

	var lastErr error
	_, err := os.Stdout.WriteString(logString)
	if err != nil {
		lastErr = err
	}
	_, err = os.Stdout.WriteString("\n")
	if err != nil {
		lastErr = err
	}

	writer.locker.Unlock()

	return lastErr
}

func (writer *AccessLogStdoutWriter) Close() {

}
