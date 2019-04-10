package teamongo

import (
	"context"
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/utils/time"
	"regexp"
	"time"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		cleanAccessLogs()
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
