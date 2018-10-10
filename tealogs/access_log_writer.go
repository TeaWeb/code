package tealogs

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teautils"
	"errors"
)

type AccessLogWriter interface {
	Init()
	Write(log *AccessLog) error
	Close()
}

func NewAccessLogWriter(config *teaconfigs.AccessLogConfig) (AccessLogWriter, error) {
	if config.Target == "file" {
		fileConfig := &teaconfigs.AccessLogFileConfig{}
		teautils.MapToObjectYAML(config.Config, fileConfig)

		writer := &AccessLogFileWriter{
			config: fileConfig,
		}

		writer.Init()
		return writer, nil
	} else if config.Target == "stdout" {
		stdoutConfig := &teaconfigs.AccessLogStdoutConfig{}
		teautils.MapToObjectYAML(config.Config, stdoutConfig)
		writer := &AccessLogStdoutWriter{
			config: stdoutConfig,
		}

		writer.Init()
		return writer, nil
	}

	return nil, errors.New("writer '" + config.Target + "' not found")
}
