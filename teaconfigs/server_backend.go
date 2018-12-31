package teaconfigs

import (
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/iwind/TeaGo/utils/string"
	"strings"
	"sync"
	"time"
)

// 服务后端配置
type ServerBackendConfig struct {
	shared.HeaderList `yaml:",inline"`

	On           bool      `yaml:"on" json:"on"`                                 // 是否启用
	Id           string    `yaml:"id" json:"id"`                                 // ID
	Code         string    `yaml:"code" json:"code"`                             // 代号
	Name         []string  `yaml:"name" json:"name"`                             // 域名 TODO
	Address      string    `yaml:"address" json:"address"`                       // 地址
	Weight       uint      `yaml:"weight" json:"weight"`                         // 是否为备份
	IsBackup     bool      `yaml:"backup" json:"isBackup"`                       // 超时时间
	FailTimeout  string    `yaml:"failTimeout" json:"failTimeout"`               // 失败超时
	MaxFails     uint      `yaml:"maxFails" json:"maxFails"`                     // 最多失败次数
	CurrentFails uint      `yaml:"currentFails" json:"currentFails"`             // 当前已失败
	MaxConns     uint      `yaml:"maxConns" json:"maxConns"`                     // 并发连接数
	IsDown       bool      `yaml:"down" json:"isDown"`                           // 是否下线
	DownTime     time.Time `yaml:"downTime,omitempty" json:"downTime,omitempty"` // 下线时间

	failTimeoutDuration time.Duration
	failsLocker         sync.Mutex

	slowStartDuration time.Duration
}

// 获取新对象
func NewServerBackendConfig() *ServerBackendConfig {
	return &ServerBackendConfig{
		On: true,
		Id: stringutil.Rand(16),
	}
}

// 校验
func (this *ServerBackendConfig) Validate() error {
	// failTimeout
	if len(this.FailTimeout) > 0 {
		this.failTimeoutDuration, _ = time.ParseDuration(this.FailTimeout)
	}

	// 是否有端口
	if strings.Index(this.Address, ":") == -1 {
		// @TODO 如果是tls，则为443
		this.Address += ":80"
	}

	// Headers
	err := this.ValidateHeaders()
	if err != nil {
		return err
	}

	return nil
}

// 超时时间
func (this *ServerBackendConfig) FailTimeoutDuration() time.Duration {
	return this.failTimeoutDuration
}

// 启动时间
func (this *ServerBackendConfig) SlowStartDuration() time.Duration {
	return this.slowStartDuration
}

// 候选对象代号
func (this *ServerBackendConfig) CandidateCodes() []string {
	codes := []string{this.Id}
	if len(this.Code) > 0 {
		codes = append(codes, this.Code)
	}
	return codes
}

// 候选对象权重
func (this *ServerBackendConfig) CandidateWeight() uint {
	return this.Weight
}

// 增加错误次数
func (this *ServerBackendConfig) IncreaseFails() uint {
	this.failsLocker.Lock()
	defer this.failsLocker.Unlock()

	this.CurrentFails ++

	return this.CurrentFails
}
