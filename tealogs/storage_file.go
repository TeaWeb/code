package tealogs

import (
	"errors"
	"github.com/iwind/TeaGo/logs"
	"os"
	"sync"
)

// 文件存储策略
type FileStorage struct {
	Storage `yaml:", inline"`

	Path string `yaml:"path" json:"path"` // 文件路径，支持变量：${year|month|week|day|hour|minute|second}

	writeLocker sync.Mutex

	files       map[string]*os.File // path => *File
	filesLocker sync.Mutex
}

// 开启
func (this *FileStorage) Start() error {
	if len(this.Path) == 0 {
		return errors.New("'path' should not be empty")
	}

	this.files = map[string]*os.File{}

	return nil
}

// 写入日志
func (this *FileStorage) Write(accessLogs []*AccessLog) error {
	if len(accessLogs) == 0 {
		return nil
	}

	fp := this.fp()
	if fp == nil {
		return errors.New("file pointer should not be nil")
	}
	this.writeLocker.Lock()
	defer this.writeLocker.Unlock()

	for _, accessLog := range accessLogs {
		data, err := this.FormatAccessLogBytes(accessLog)
		if err != nil {
			logs.Error(err)
			continue
		}
		_, err = fp.Write(data)
		if err != nil {
			this.Close()
			break
		}
		fp.WriteString("\n")
	}
	return nil
}

// 关闭
func (this *FileStorage) Close() error {
	this.filesLocker.Lock()
	defer this.filesLocker.Unlock()

	var resultErr error
	for _, f := range this.files {
		err := f.Close()
		if err != nil {
			resultErr = err
		}
	}
	return resultErr
}

func (this *FileStorage) fp() *os.File {
	path := this.FormatVariables(this.Path)

	this.filesLocker.Lock()
	defer this.filesLocker.Unlock()
	fp, ok := this.files[path]
	if ok {
		return fp
	}

	// 关闭其他的文件
	for _, f := range this.files {
		f.Close()
	}

	// 打开新文件
	fp, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logs.Error(err)
		return nil
	}
	this.files[path] = fp

	return fp
}
