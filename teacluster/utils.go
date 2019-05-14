package teacluster

import (
	"github.com/TeaWeb/code/teacluster/configs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"path/filepath"
	"strings"
)

func BuildSum() {
	sumList := []string{}
	RangeFiles(func(file *files.File, relativePath string) {
		sum, err := file.Md5()
		if err != nil {
			logs.Error(err)
			return
		}
		sumList = append(sumList, relativePath+"|"+sum)
	})

	sumFile := files.NewFile(Tea.ConfigFile("node.sum"))
	err := sumFile.WriteString(strings.Join(sumList, "\n"))
	if err != nil {
		logs.Error(err)
	}
}

func PushItems() {
	action := &PushAction{}

	// proxy
	RangeFiles(func(file *files.File, relativePath string) {
		data, err := file.ReadAll()
		if err != nil {
			logs.Error(err)
			return
		}

		item := configs.NewItem()
		item.Id = relativePath
		item.Data = data

		sum, err := file.Md5()
		if err != nil {
			logs.Error(err)
			return
		}
		item.Sum = sum

		stat, err := file.Stat()
		if err == nil {
			item.Flags = []int{int(stat.Mode.Perm())}
		}
		action.AddItem(item)
	})

	err := ClusterManager.Write(action)
	if err != nil {
		logs.Error(err)
	}
}

func PullItems() {
	logs.Println("pull items")
}

func RangeFiles(f func(file *files.File, relativePath string)) {
	configAbs, err := filepath.Abs(Tea.ConfigDir())
	if err != nil {
		logs.Error(err)
		return
	}
	files.NewFile(configAbs).Range(func(file *files.File) {
		// *.conf & ssl.*
		if !strings.HasSuffix(file.Name(), ".conf") && !strings.HasPrefix(file.Name(), "ssl.") {
			return
		}
		if lists.ContainsString([]string{"node.conf"}, file.Name()) {
			return
		}
		absPath, _ := file.AbsPath()
		relativePath := strings.TrimLeft(strings.TrimPrefix(absPath, configAbs), "/\\")
		f(file, relativePath)
	})
}
