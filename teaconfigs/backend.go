package teaconfigs

import (
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/iwind/TeaGo/utils/string"
	"strings"
	"sync"
	"time"
)

// 服务后端配置
type BackendConfig struct {
	shared.HeaderList `yaml:",inline"`

	On           bool      `yaml:"on" json:"on"`                                 // 是否启用
	Id           string    `yaml:"id" json:"id"`                                 // ID
	Code         string    `yaml:"code" json:"code"`                             // 代号
	Name         []string  `yaml:"name" json:"name"`                             // 域名 TODO
	Address      string    `yaml:"address" json:"address"`                       // 地址
	Scheme       string    `yaml:"scheme" json:"scheme"`                         // 协议，http或者https
	Weight       uint      `yaml:"weight" json:"weight"`                         // 是否为备份
	IsBackup     bool      `yaml:"backup" json:"isBackup"`                       // 超时时间
	FailTimeout  string    `yaml:"failTimeout" json:"failTimeout"`               // 失败超时
	MaxFails     uint      `yaml:"maxFails" json:"maxFails"`                     // 最多失败次数
	CurrentFails uint      `yaml:"currentFails" json:"currentFails"`             // 当前已失败次数
	MaxConns     uint      `yaml:"maxConns" json:"maxConns"`                     // 最大并发连接数
	CurrentConns uint      `yaml:"currentConns" json:"currentConns"`             // 当前连接数
	IsDown       bool      `yaml:"down" json:"isDown"`                           // 是否下线
	DownTime     time.Time `yaml:"downTime,omitempty" json:"downTime,omitempty"` // 下线时间

	failTimeoutDuration time.Duration
	failsLocker         sync.Mutex
	connsLocker         sync.Mutex
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

	return nil
}

// 超时时间
func (this *BackendConfig) FailTimeoutDuration() time.Duration {
	return this.failTimeoutDuration
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
func (this *BackendConfig) IncreaseFails() uint {
	this.failsLocker.Lock()
	defer this.failsLocker.Unlock()

	this.CurrentFails ++

	return this.CurrentFails
}

// 增加连接数
func (this *BackendConfig) IncreaseConn() {
	this.connsLocker.Lock()
	defer this.connsLocker.Unlock()

	this.CurrentConns ++
}

// 减少连接数
func (this *BackendConfig) DecreaseConn() {
	this.connsLocker.Lock()
	defer this.connsLocker.Unlock()

	if this.CurrentConns > 0 {
		this.CurrentConns --
	}
}
