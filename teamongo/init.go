package teamongo

import (
	"context"
	"github.com/TeaWeb/code/teahooks"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/processes"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/utils/time"
	"regexp"
	"runtime"
	"time"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		cleanAccessLogs()
	})

	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		startInstalledMongo()
	})

	teahooks.On(teahooks.EventReload, func() {
		RestartClient()
	})
}

// 清理访问日志任务
func cleanAccessLogs() {
	reg := regexp.MustCompile("^logs\\.\\d{8}$")
	timers.Loop(1*time.Minute, func(looper *timers.Looper) {
		config, _ := configs.LoadMongoConfig()
		now := time.Now()
		if config != nil && config.AccessLog != nil &&
			config.AccessLog.CleanHour == now.Hour() &&
			now.Minute() == 0 &&
			config.AccessLog.KeepDays >= 1 {
			compareDay := "logs." + timeutil.Format("Ymd", time.Now().Add(-time.Duration(config.AccessLog.KeepDays * 24)*time.Hour))
			logs.Println("[mongo]clean access logs before '" + compareDay + "'")

			db := SharedClient().Database(DatabaseName)
			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			cursor, err := db.ListCollections(ctx, maps.Map{})
			if err != nil {
				logs.Error(err)
				return
			}
			defer cursor.Close(context.Background())

			for cursor.Next(context.Background()) {
				m := maps.Map{}
				err := cursor.Decode(&m)
				if err != nil {
					logs.Error(err)
					return
				}
				name := m.GetString("name")
				if len(name) == 0 {
					continue
				}
				if !reg.MatchString(name) {
					continue
				}

				if name < compareDay {
					logs.Println("[mongo]clean collection '" + name + "'")
					ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
					err := db.Collection(name).Drop(ctx)
					if err != nil {
						logs.Error(err)
					}
				}
			}
		}
	})
}

// 启动本机安装的Mongo
func startInstalledMongo() {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		return
	}

	config := configs.SharedMongoConfig()
	if config.Host != "127.0.0.1" && config.Host != "localhost" {
		return
	}

	err := Test()

	if err != nil {
		mongodbDir := Tea.Root + "/mongodb"

		// 是否已安装
		if !files.NewFile(mongodbDir + "/bin/mongod").Exists() {
			return
		}

		// 启动
		p := processes.NewProcess(mongodbDir+"/bin/mongod", "--dbpath="+mongodbDir+"/data", "--fork", "--logpath="+mongodbDir+"/data/fork.log")
		p.SetPwd(mongodbDir)

		logs.Println("[mongo]start mongo: ", mongodbDir+"/bin/mongod", "--dbpath="+mongodbDir+"/data", "--fork", "--logpath="+mongodbDir+"/data/fork.log")

		err := p.StartBackground()
		if err != nil {
			logs.Println("[mongo]start error: " + err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}
