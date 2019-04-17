package teaconfigs

import (
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/utils/string"
	"strings"
	"sync/atomic"
	"time"
)

// 服务后端配置
type BackendConfig struct {
	shared.HeaderList `yaml:",inline"`

	On              bool      `yaml:"on" json:"on"`                                 // 是否启用
	Id              string    `yaml:"id" json:"id"`                                 // ID
	Code            string    `yaml:"code" json:"code"`                             // 代号
	Address         string    `yaml:"address" json:"address"`                       // 地址
	Scheme          string    `yaml:"scheme" json:"scheme"`                         // 协议，http或者https
	Weight          uint      `yaml:"weight" json:"weight"`                         // 权重
	IsBackup        bool      `yaml:"backup" json:"isBackup"`                       // 是否为备份
	FailTimeout     string    `yaml:"failTimeout" json:"failTimeout"`               // 连接失败超时
	ReadTimeout     string    `yaml:"readTimeout" json:"readTimeout"`               // 读取超时时间
	MaxFails        int32     `yaml:"maxFails" json:"maxFails"`                     // 最多失败次数
	CurrentFails    int32     `yaml:"currentFails" json:"currentFails"`             // 当前已失败次数
	MaxConns        int32     `yaml:"maxConns" json:"maxConns"`                     // 最大并发连接数
	CurrentConns    int32     `yaml:"currentConns" json:"currentConns"`             // 当前连接数
	IsDown          bool      `yaml:"down" json:"isDown"`                           // 是否下线
	DownTime        time.Time `yaml:"downTime,omitempty" json:"downTime,omitempty"` // 下线时间
	RequestGroupIds []string  `yaml:"requestGroupIds" json:"requestGroupIds"`       // 所属请求分组
	RequestURI      string    `yaml:"requestURI" json:"requestURI"`                 // 转发后的请求URI

	failTimeoutDuration time.Duration
	readTimeoutDuration time.Duration
	hasRequestURI       bool
	requestPath         string
	requestArgs         string
}

// 获取新对象
func NewBackendConfig() *BackendConfig {
	return &BackendConfig{
		On: true,
		Id: stringutil.Rand(16),
	}
}

// 校验
func (this *BackendConfig) Validate() error {
	// failTimeout
	if len(this.FailTimeout) > 0 {
		this.failTimeoutDuration, _ = time.ParseDuration(this.FailTimeout)
	}

	// readTimeout
	if len(this.ReadTimeout) > 0 {
		this.readTimeoutDuration, _ = time.ParseDuration(this.ReadTimeout)
	}

	// 是否有端口
	if strings.Index(this.Address, ":") == -1 {
		if this.Scheme == "https" {
			this.Address += ":443"
		} else {
			this.Address += ":80"
		}
	}

	// Headers
	err := this.ValidateHeaders()
	if err != nil {
		return err
	}

	// request uri
	if len(this.RequestURI) == 0 || this.RequestURI == "${requestURI}" {
		this.hasRequestURI = false
	} else {
		this.hasRequestURI = true

		if strings.Contains(this.RequestURI, "?") {
			pieces := strings.SplitN(this.RequestURI, "?", -1)
			this.requestPath = pieces[0]
			this.requestArgs = pieces[1]
		} else {
			this.requestPath = this.RequestURI
		}
	}

	return nil
}

// 连接超时时间
func (this *BackendConfig) FailTimeoutDuration() time.Duration {
	return this.failTimeoutDuration
}

// 读取超时时间
func (this *BackendConfig) ReadTimeoutDuration() time.Duration {
	return this.readTimeoutDuration
}

// 候选对象代号
func (this *BackendConfig) CandidateCodes() []string {
	codes := []string{this.Id}
	if len(this.Code) > 0 {
		codes = append(codes, this.Code)
	}
	return codes
}

// 候选对象权重
func (this *BackendConfig) CandidateWeight() uint {
	return this.Weight
}

// 增加错误次数
func (this *BackendConfig) IncreaseFails() int32 {
	atomic.AddInt32(&this.CurrentFails, 1)
	return this.CurrentFails
}

// 增加连接数
func (this *BackendConfig) IncreaseConn() {
	atomic.AddInt32(&this.CurrentConns, 1)
}

// 减少连接数
func (this *BackendConfig) DecreaseConn() {
	atomic.AddInt32(&this.CurrentConns, -1)
}

// 添加请求分组
func (this *BackendConfig) AddRequestGroupId(requestGroupId string) {
	this.RequestGroupIds = append(this.RequestGroupIds, requestGroupId)
}

// 删除某个请求分组
func (this *BackendConfig) RemoveRequestGroupId(requestGroupId string) {
	result := []string{}
	for _, groupId := range this.RequestGroupIds {
		if groupId == requestGroupId {
			continue
		}
		result = append(result, groupId)
	}
	this.RequestGroupIds = result
}

// 判断是否有某个情趣分组ID
func (this *BackendConfig) HasRequestGroupId(requestGroupId string) bool {
	if requestGroupId == "default" && len(this.RequestGroupIds) == 0 {
		return true
	}
	return lists.ContainsString(this.RequestGroupIds, requestGroupId)
}

// 判断是否设置RequestURI
func (this *BackendConfig) HasRequestURI() bool {
	return this.hasRequestURI
}

// 获取转发后的Path
func (this *BackendConfig) RequestPath() string {
	return this.requestPath
}

// 获取转发后的附加参数
func (this *BackendConfig) RequestArgs() string {
	return this.requestArgs
}
