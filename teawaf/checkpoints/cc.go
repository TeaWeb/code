package checkpoints

import (
	"github.com/TeaWeb/code/teamemory"
	"github.com/TeaWeb/code/teawaf/requests"
	"github.com/iwind/TeaGo/types"
	"net"
	"regexp"
	"strings"
	"sync"
)

// ${cc.arg}
// TODO implement more traffic rules
type CCCheckpoint struct {
	Checkpoint

	grid *teamemory.Grid
	once sync.Once
}

func (this *CCCheckpoint) Init() {

}

func (this *CCCheckpoint) Start() {
	if this.grid != nil {
		this.grid.Destroy()
	}
	this.grid = teamemory.NewGrid(100)
}

func (this *CCCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	value = 0

	if this.grid == nil {
		this.once.Do(func() {
			this.Start()
		})
		if this.grid == nil {
			return
		}
	}

	periodString, ok := options["period"]
	if !ok {
		return
	}
	period := types.Int64(periodString)
	if period < 1 {
		return
	}

	if param == "requests" { // requests
		key := this.ip(req)
		value = this.grid.IncreaseInt64([]byte(key), 1, period)
	}

	return
}

func (this *CCCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}

func (this *CCCheckpoint) ParamOptions() *ParamOptions {
	option := NewParamOptions()
	option.AddParam("请求数", "requests")
	return option
}

func (this *CCCheckpoint) Options() []*Option {
	options := []*Option{}

	// period
	{
		option := NewOption("统计周期", "period")
		option.Value = "60"
		option.RightLabel = "秒"
		option.Size = 8
		option.MaxLength = 8
		option.Validate = func(value string) (ok bool, message string) {
			if regexp.MustCompile("^\\d+$").MatchString(value) {
				ok = true
				return
			}
			message = "周期需要是一个整数数字"
			return
		}
		options = append(options, option)
	}

	return options
}

func (this *CCCheckpoint) Stop() {
	if this.grid != nil {
		this.grid.Destroy()
		this.grid = nil
	}
}

func (this *CCCheckpoint) ip(req *requests.Request) string {
	// X-Forwarded-For
	forwardedFor := req.Header.Get("X-Forwarded-For")
	if len(forwardedFor) > 0 {
		commaIndex := strings.Index(forwardedFor, ",")
		if commaIndex > 0 {
			return forwardedFor[:commaIndex]
		}
		return forwardedFor
	}

	// Real-IP
	{
		realIP, ok := req.Header["X-Real-IP"]
		if ok && len(realIP) > 0 {
			return realIP[0]
		}
	}

	// Real-Ip
	{
		realIP, ok := req.Header["X-Real-Ip"]
		if ok && len(realIP) > 0 {
			return realIP[0]
		}
	}

	// Remote-Addr
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		return host
	}
	return req.RemoteAddr
}
