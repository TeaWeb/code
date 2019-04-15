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
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// 文件缓存管理器
type FileManager struct {
	Manager

	Capacity float64       // 容量
	Life     time.Duration // 有效期

	looper      *timers.Looper
	dir         string
	writeLocker sync.RWMutex
}

func NewFileManager() *FileManager {
	manager := &FileManager{}

	// 删除过期
	manager.looper = timers.Loop(30*time.Minute, func(looper *timers.Looper) {
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
					data := reader.Read(12)
					if len(data) != 12 {
						continue
					}
					timestamp := types.Int64(string(data))
					reader.Close()
					if timestamp < time.Now().Unix()-100 { // 超时100秒以上的
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
	if this.Life <= 0 {
		this.Life = 1800 * time.Second
	}

	dir, found := options["dir"]
	if found {
		this.dir = types.String(dir)
	}
}

// 写入
func (this *FileManager) Write(key string, data []byte) error {
	if len(this.dir) == 0 {
		return errors.New("cache dir should not be empty")
	}

	this.writeLocker.Lock()
	defer this.writeLocker.Unlock()

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

	// 头部加入有效期
	var life = int64(this.Life.Seconds())
	if life <= 0 {
		life = 30 * 86400
	} else if life >= 365*86400 { // 最大值限制
		life = 365 * 86400
	}
	data = append([]byte(fmt.Sprintf("%012d", time.Now().Unix()+life)), data ...)
	err := newFile.Write(data)

	return err
}

// 读取
func (this *FileManager) Read(key string) (data []byte, err error) {
	md5 := stringutil.Md5(key)
	file := files.NewFile(this.dir + Tea.DS + md5[:2] + Tea.DS + md5[2:4] + Tea.DS + md5 + ".cache")

	this.writeLocker.RLock()
	defer this.writeLocker.RUnlock()

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

// 删除
func (this *FileManager) Delete(key string) error {
	if len(this.dir) == 0 {
		return errors.New("cache dir should not be empty")
	}

	this.writeLocker.Lock()
	defer this.writeLocker.Unlock()

	dirFile := files.NewFile(this.dir)
	if !dirFile.IsDir() {
		return errors.New("cache dir should be a valid dir")
	}

	md5 := stringutil.Md5(key)
	newDir := files.NewFile(this.dir + Tea.DS + md5[:2] + Tea.DS + md5[2:4])
	if !newDir.Exists() {
		return nil
	}

	newFile := files.NewFile(newDir.Path() + Tea.DS + md5 + ".cache")
	if !newFile.Exists() {
		return nil
	}
	return newFile.Delete()
}

// 统计
func (this *FileManager) Stat() (size int64, countKeys int, err error) {
	// 检查目录是否存在
	info, err := os.Stat(this.dir)
	if err != nil || !info.IsDir() {
		return
	}

	filepath.Walk(this.dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".cache") {
			return nil
		}
		size += info.Size()
		countKeys ++

		return nil
	})

	return
}

// 清理
func (this *FileManager) Clean() error {
	dirReg := regexp.MustCompile("^[0-9a-f]{2}$")
	for _, file := range files.NewFile(this.dir).List() {
		if !file.IsDir() {
			continue
		}
		if !dirReg.MatchString(file.Name()) {
			continue
		}

		err := file.DeleteAll()
		if err != nil {
			logs.Error(err)
		}
	}
	return nil
}

func (this *FileManager) Close() error {
	//logs.Println("[cache]close cache policy instance: file")
	if this.looper != nil {
		this.looper.Stop()
		this.looper = nil
	}

	// TODO 删除所有文件和目录

	return nil
}
