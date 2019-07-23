package teaproxy

import (
	"github.com/TeaWeb/code/teahooks"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/logs"
	"sync/atomic"
	"time"
)

// 当前QPS
var qps = int32(0)

// 对外可读取的QPS
var QPS = int32(0)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 计算QPS
		teautils.Every(1*time.Second, func(ticker *teautils.Ticker) {
			QPS = qps
			atomic.StoreInt32(&qps, 0)
		})
	})

	teahooks.On(teahooks.EventReload, func() {
		// 重启服务
		err := SharedManager.Restart()
		if err != nil {
			logs.Error(err)
		}
	})
}
