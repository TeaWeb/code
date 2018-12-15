package teacache

import (
	"errors"
	"fmt"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"sync"
	"time"
)

// 文件缓存管理器
type FileManager struct {
	Capacity float64       // 容量
	Life     time.Duration // 有效期

	dir          string
	writingFiles map[string]bool // file path => true
	writeLocker  sync.Mutex
}

func NewFileManager() *FileManager {
	manager := &FileManager{}
	manager.writingFiles = map[string]bool{}

	// 删除过期
	timers.Loop(10*time.Minute, func(looper *timers.Looper) {
		if len(manager.dir) == 0 {
			return
		}
		dirFile := files.NewFile(manager.dir)
		if !dirFile.IsDir() {
			return
		}
		for _, dirFile1 := range dirFile.List() {
			if !dirFile1.IsDir() {
				continue
			}
			for _, dirFile2 := range dirFile1.List() {
				for _, file := range dirFile2.List() {
					if file.Ext() != ".cache" {
						continue
					}
					reader, err := file.Reader()
					if err != nil {
						logs.Error(err)
						continue
					}
					timestamp := types.Int64(string(reader.Read(12)))
					reader.Close()
					if timestamp < time.Now().Unix() {
						err := file.Delete()
						if err != nil {
							logs.Error(err)
						}
					}
				}

				time.Sleep(500 * time.Millisecond)
			}
		}
	})

	return manager
}

func (this *FileManager) SetOptions(options map[string]interface{}) {
	dir, found := options["dir"]
	if found {
		this.dir = types.String(dir)
	}
}

func (this *FileManager) Write(key string, data []byte) error {
	if len(this.dir) == 0 {
		return errors.New("cache dir should not be empty")
	}

	dirFile := files.NewFile(this.dir)
	if !dirFile.IsDir() {
		return errors.New("cache dir should be a valid dir")
	}

	md5 := stringutil.Md5(key)
	newDir := files.NewFile(this.dir + Tea.DS + md5[:2] + Tea.DS + md5[2:4])
	if !newDir.Exists() {
		err := newDir.MkdirAll()
		if err != nil {
			return err
		}
	}

	newFile := files.NewFile(newDir.Path() + Tea.DS + md5 + ".cache")
	if this.isLocking(newFile) {
		return errors.New("file is locking")
	}

	this.writeLocker.Lock()
	this.writingFiles[newFile.Path()] = true
	this.writeLocker.Unlock()

	// 头部加入有效期
	var life = int64(this.Life.Seconds())
	if life <= 0 {
		life = 30 * 86400
	}
	data = append([]byte(fmt.Sprintf("%012d", time.Now().Unix()+life)), data ...)
	err := newFile.Write(data)

	// 解除锁定
	this.writeLocker.Lock()
	defer this.writeLocker.Unlock()
	delete(this.writingFiles, newFile.Path())
	return err
}

func (this *FileManager) Read(key string) (data []byte, err error) {
	md5 := stringutil.Md5(key)
	file := files.NewFile(this.dir + Tea.DS + md5[:2] + Tea.DS + md5[2:4] + Tea.DS + md5 + ".cache")
	if this.isLocking(file) {
		return nil, errors.New("file is locking")
	}
	if !file.Exists() {
		return nil, ErrNotFound
	}
	data, err = file.ReadAll()
	if err != nil || len(data) < 12 { // 12是时间戳长度
		return nil, err
	}
	timestamp := types.Int64(string(data[:12]))
	if timestamp < time.Now().Unix() {
		return nil, ErrNotFound
	}
	return data[12:], nil
}

func (this *FileManager) isLocking(file *files.File) bool {
	this.writeLocker.Lock()
	defer this.writeLocker.Unlock()

	path := file.Path()
	_, found := this.writingFiles[path]
	return found
}
