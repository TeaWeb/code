package teautils

import (
	"github.com/fsnotify/fsnotify"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"os"
	"runtime"
	"sync"
	"time"
)

// 文件统计Map
const (
	MaxWatchingFiles = 10240
)

var (
	fileStatMap     = map[string]*FileInfo{} // path => os.FileInfo
	fileStatLocker  = &sync.RWMutex{}
	fileWatcher, _  = fsnotify.NewWatcher()
	fileStatEnabled = runtime.GOOS == "linux" || runtime.GOOS == "darwin"
)

type FileInfo struct {
	info       os.FileInfo
	accessTime int64
}

// 初始化
func init() {
	if fileStatEnabled && fileWatcher != nil {
		go func() {
			for event := range fileWatcher.Events {
				logs.Println("[watcher]changed:", event.Op, event.Name)

				fileStatLocker.Lock()
				delete(fileStatMap, event.Name)
				fileStatLocker.Unlock()
			}
		}()
	}
}

// 临时文件
func TmpFile(path string) string {
	return Tea.Root + Tea.DS + "web" + Tea.DS + "tmp" + Tea.DS + path
}

// 文件统计信息
func StatFile(path string) (os.FileInfo, error) {
	if fileWatcher == nil || !fileStatEnabled {
		return os.Stat(path)
	}

	fileStatLocker.RLock()
	fileInfo, ok := fileStatMap[path]
	fileStatLocker.RUnlock()
	if ok {
		fileInfo.accessTime = time.Now().Unix()
		return fileInfo.info, nil
	}

	stat, err := os.Stat(path)
	if err != nil {
		return stat, err
	}

	fileStatLocker.Lock()
	if len(fileStatMap) >= MaxWatchingFiles {
		timestamp := time.Now().Unix()
		hasDeleted := false
		for k, v := range fileStatMap {
			if timestamp-v.accessTime >= 1800 { // 1800 seconds
				delete(fileStatMap, k)
				_ = fileWatcher.Remove(k)
				hasDeleted = true
			}
		}

		// no space left
		if !hasDeleted {
			fileStatLocker.Unlock()
			return stat, err
		}
	}
	fileStatLocker.Unlock()

	// watch file
	err = fileWatcher.Add(path)

	if err != nil {
		logs.Error(err)
		_ = fileWatcher.Remove(path)
	} else {
		fileStatLocker.Lock()
		fileStatMap[path] = &FileInfo{
			info:       stat,
			accessTime: time.Now().Unix(),
		}
		fileStatLocker.Unlock()
	}

	return stat, nil
}
