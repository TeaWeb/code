package teaproxy

import (
	"github.com/TeaWeb/code/teahooks"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"net/http"
	"sync/atomic"
	"time"
)

// 状态码筛选
var StatusCodeParser func(statusCode int, headers http.Header, respData []byte, parserScript string) (string, error) = nil

// 当前QPS
var qps = int32(0)

// 对外可读取的QPS
var QPS = int32(0)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 计算QPS
		timers.Every(1*time.Second, func(ticker *time.Ticker) {
			QPS = qps
			atomic.StoreInt32(&qps, 0)
		})
	})

	teahooks.On(teahooks.EventReload, func() {
		err := SharedManager.Restart()
		if err != nil {
			logs.Error(err)
		}
	})
}
