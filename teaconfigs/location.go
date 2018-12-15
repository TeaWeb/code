package teaconfigs

import (
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/utils/string"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

const (
	LocationPatternTypePrefix = 1
	LocationPatternTypeExact  = 2
	LocationPatternTypeRegexp = 3
)

// 路径配置
// @TODO 匹配的时候去除路径中多余的斜杠（/）
type LocationConfig struct {
	On      bool   `yaml:"on" json:"on"`           // 是否开启 @TODO
	Id      string `yaml:"id" json:"id"`           // @TODO
	Pattern string `yaml:"pattern" json:"pattern"` // 匹配规则  @TODO

	patternType int // 规则类型：LocationPattern*

	prefix string // 前缀
	path   string // 精确的路径

	reg *regexp.Regexp // 匹配规则

	caseInsensitive bool // 大小写不敏感
	reverse         bool // 是否翻转规则，比如非前缀，非路径

	Async   bool     `yaml:"async" json:"async"`     // 是否异步请求 @TODO
	Notify  []string `yaml:"notify" json:"notify"`   // 转发请求 @TODO
	LogOnly bool     `yaml:"logOnly" json:"logOnly"` // 是否只记录日志 @TODO
	Root    string   `yaml:"root" json:"root"`       // 资源根目录
	Index   []string `yaml:"index" json:"index"`     // 默认文件
	Charset string   `yaml:"charset" json:"charset"` // 字符集设置

	// 日志
	AccessLog []*AccessLogConfig `yaml:"accessLog" json:"accessLog"` // @TODO

	Headers       []*shared.HeaderConfig `yaml:"headers" json:"headers"`             // 添加的头信息
	IgnoreHeaders []string               `yaml:"ignoreHeaders" json:"ignoreHeaders"` // 忽略的Header

	// 参考：http://nginx.org/en/docs/http/ngx_http_access_module.html
	Allow []string `yaml:"allow" json:"allow"` // 允许的终端地址 @TODO
	Deny  []string `yaml:"deny" json:"deny"`   // 禁止的终端地址 @TODO

	Rewrite  []*RewriteRule         `yaml:"rewrite" json:"rewrite"`   // 重写规则
	Fastcgi  []*FastcgiConfig       `yaml:"fastcgi" json:"fastcgi"`   // Fastcgi配置 @TODO
	Proxy    string                 `yaml:proxy" json:"proxy"`        //  代理配置 @TODO
	Backends []*ServerBackendConfig `yaml:"backends" json:"backends"` // 后端服务器配置 @TODO

	CachePolicy string `yaml:"cachePolicy" json:"cachePolicy"` // 缓存策略
	CacheOn     bool   `yaml:"cacheOn" json:"cacheOn"`         // 缓存是否打开 TODO
	cachePolicy *shared.CachePolicy
}

func NewLocationConfig() *LocationConfig {
	return &LocationConfig{
		On: true,
		Id: stringutil.Rand(16),
	}
}

func (this *LocationConfig) Validate() error {
	// 分析pattern
	this.reverse = false
	this.caseInsensitive = false
	if len(this.Pattern) > 0 {
		spaceIndex := strings.Index(this.Pattern, " ")
		if spaceIndex < 0 {
			this.patternType = LocationPatternTypePrefix
			this.prefix = this.Pattern
		} else {
			cmd := this.Pattern[:spaceIndex]
			pattern := strings.TrimSpace(this.Pattern[spaceIndex+1:])
			if cmd == "*" { // 大小写非敏感
				this.patternType = LocationPatternTypePrefix
				this.prefix = pattern
				this.caseInsensitive = true
			} else if cmd == "!*" { // 大小写非敏感，翻转
				this.patternType = LocationPatternTypePrefix
				this.prefix = pattern
				this.caseInsensitive = true
				this.reverse = true
			} else if cmd == "!" {
				this.patternType = LocationPatternTypePrefix
				this.prefix = pattern
				this.reverse = true
			} else if cmd == "=" {
				this.patternType = LocationPatternTypeExact
				this.path = pattern
			} else if cmd == "=*" {
				this.patternType = LocationPatternTypeExact
				this.path = pattern
				this.caseInsensitive = true
			} else if cmd == "!=" {
				this.patternType = LocationPatternTypeExact
				this.path = pattern
				this.reverse = true
			} else if cmd == "!=*" {
				this.patternType = LocationPatternTypeExact
				this.path = pattern
				this.reverse = true
				this.caseInsensitive = true
			} else if cmd == "~" { // 正则
				this.patternType = LocationPatternTypeRegexp
				reg, err := regexp.Compile(pattern)
				if err != nil {
					return err
				}
				this.reg = reg
				this.path = pattern
			} else if cmd == "!~" {
				this.patternType = LocationPatternTypeRegexp
				reg, err := regexp.Compile(pattern)
				if err != nil {
					return err
				}
				this.reg = reg
				this.reverse = true
				this.path = pattern
			} else if cmd == "~*" { // 大小写非敏感小写
				this.patternType = LocationPatternTypeRegexp
				reg, err := regexp.Compile("(?i)" + pattern)
				if err != nil {
					return err
				}
				this.reg = reg
				this.caseInsensitive = true
				this.path = pattern
			} else if cmd == "!~*" {
				this.patternType = LocationPatternTypeRegexp
				reg, err := regexp.Compile("(?i)" + pattern)
				if err != nil {
					return err
				}
				this.reg = reg
				this.reverse = true
				this.caseInsensitive = true
				this.path = pattern
			} else {
				this.patternType = LocationPatternTypePrefix
				this.prefix = pattern
			}
		}
	} else {
		this.patternType = LocationPatternTypePrefix
		this.prefix = this.Pattern
	}

	// 校验缓存配置
	if len(this.CachePolicy) > 0 {
		policy := shared.NewCachePolicyFromFile(this.CachePolicy)
		if policy != nil {
			err := policy.Validate()
			if err != nil {
				return err
			}
			this.cachePolicy = policy
		}
	}

	// 校验RewriteRule配置
	for _, rewriteRule := range this.Rewrite {
		err := rewriteRule.Validate()
		if err != nil {
			return err
		}
	}

	// 校验Fastcgi配置
	for _, fastcgi := range this.Fastcgi {
		err := fastcgi.Validate()
		if err != nil {
			return err
		}
	}

	// 校验Header
	for _, header := range this.Headers {
		err := header.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// 模式类型
func (this *LocationConfig) PatternType() int {
	return this.patternType
}

// 模式字符串
// 去掉了模式字符
func (this *LocationConfig) PatternString() string {
	if this.patternType == LocationPatternTypePrefix {
		return this.prefix
	}
	return this.path
}

// 是否翻转
func (this *LocationConfig) IsReverse() bool {
	return this.reverse
}

// 是否大小写非敏感
func (this *LocationConfig) IsCaseInsensitive() bool {
	return this.caseInsensitive
}

// 判断是否匹配路径
func (this *LocationConfig) Match(path string) ([]string, bool) {
	if this.patternType == LocationPatternTypePrefix {
		if this.reverse {
			if this.caseInsensitive {
				return []string{}, !strings.HasPrefix(strings.ToLower(path), strings.ToLower(this.prefix))
			} else {
				return []string{}, !strings.HasPrefix(path, this.prefix)
			}
		} else {
			if this.caseInsensitive {
				return []string{}, strings.HasPrefix(strings.ToLower(path), strings.ToLower(this.prefix))
			} else {
				return []string{}, strings.HasPrefix(path, this.prefix)
			}
		}
	}

	if this.patternType == LocationPatternTypeExact {
		if this.reverse {
			if this.caseInsensitive {
				return []string{}, strings.ToLower(path) != strings.ToLower(this.path)
			} else {
				return []string{}, path != this.path
			}
		} else {
			if this.caseInsensitive {
				return []string{}, strings.ToLower(path) == strings.ToLower(this.path)
			} else {
				return []string{}, path == this.path
			}
		}
	}

	if this.patternType == LocationPatternTypeRegexp {
		if this.reg != nil {
			if this.reverse {
				return []string{}, !this.reg.MatchString(path)
			} else {
				b := this.reg.MatchString(path)
				if b {
					return this.reg.FindStringSubmatch(path), true
				}
				return []string{}, b
			}
		}

		return []string{}, this.reverse
	}

	return []string{}, false
}

// 组合参数为一个字符串
func (this *LocationConfig) SetPattern(pattern string, patternType int, caseInsensitive bool, reverse bool) {
	op := ""
	if patternType == LocationPatternTypePrefix {
		if caseInsensitive {
			op = "*"
			if reverse {
				op = "!*"
			}
		} else {
			if reverse {
				op = "!"
			}
		}
	} else if patternType == LocationPatternTypeExact {
		op = "="
		if caseInsensitive {
			op += "*"
		}
		if reverse {
			op = "!" + op
		}
	} else if patternType == LocationPatternTypeRegexp {
		op = "~"
		if caseInsensitive {
			op += "*"
		}
		if reverse {
			op = "!" + op
		}
	}
	if len(op) > 0 {
		pattern = op + " " + pattern
	}
	this.Pattern = pattern
}

// 取得下一个可用的后端服务
// @TODO 实现backend中的各种参数
func (this *LocationConfig) NextBackend() *ServerBackendConfig {
	if len(this.Backends) == 0 {
		return nil
	}

	availableBackends := []*ServerBackendConfig{}
	for _, backend := range this.Backends {
		if backend.On && !backend.IsDown {
			availableBackends = append(availableBackends, backend)
		}
	}

	countBackends := len(availableBackends)
	if countBackends == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Int() % countBackends
	return availableBackends[index]
}

// 取得下一个可用的fastcgi
// @TODO 实现fastcgi中的各种参数
func (this *LocationConfig) NextFastcgi() *FastcgiConfig {
	if len(this.Fastcgi) == 0 {
		return nil
	}

	availableServers := []*FastcgiConfig{}
	for _, f := range this.Fastcgi {
		if !f.On {
			continue
		}
		availableServers = append(availableServers, f)
	}
	if len(availableServers) == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Int() % len(availableServers)
	return this.Fastcgi[index]
}

// 添加fastcgi配置
func (this *LocationConfig) AddFastcgi(fastcgi *FastcgiConfig) {
	this.Fastcgi = append(this.Fastcgi, fastcgi)
}

// 取得在某个的fastcgi配置
func (this *LocationConfig) FastcgiAtIndex(index int) *FastcgiConfig {
	if index < 0 || index >= len(this.Fastcgi) {
		return nil
	}
	return this.Fastcgi[index]
}

// 移除某个fastcgi配置
func (this *LocationConfig) RemoveFastcgiAt(index int) {
	if index < 0 || index >= len(this.Fastcgi) {
		return
	}
	if index == 0 {
		this.Fastcgi = this.Fastcgi[1:]
	} else if index == len(this.Fastcgi)-1 {
		this.Fastcgi = this.Fastcgi[:index]
	} else {
		this.Fastcgi = append(this.Fastcgi[:index], this.Fastcgi[index+1:] ...)
	}
}

// 设置Header
func (this *LocationConfig) SetHeader(name string, value string) {
	found := false
	upperName := strings.ToUpper(name)
	for _, header := range this.Headers {
		if strings.ToUpper(header.Name) == upperName {
			found = true
			header.Value = value
		}
	}
	if found {
		return
	}

	header := shared.NewHeaderConfig()
	header.Name = name
	header.Value = value
	this.Headers = append(this.Headers, header)
}

// 删除指定位置上的Header
func (this *LocationConfig) DeleteHeaderAtIndex(index int) {
	if index >= 0 && index < len(this.Headers) {
		this.Headers = lists.Remove(this.Headers, index).([]*shared.HeaderConfig)
	}
}

// 取得指定位置上的Header
func (this *LocationConfig) HeaderAtIndex(index int) *shared.HeaderConfig {
	if index >= 0 && index < len(this.Headers) {
		return this.Headers[index]
	}
	return nil
}

// 格式化Header
func (this *LocationConfig) FormatHeaders(formatter func(source string) string) []*shared.HeaderConfig {
	result := []*shared.HeaderConfig{}
	for _, header := range this.Headers {
		result = append(result, &shared.HeaderConfig{
			Name:   header.Name,
			Value:  formatter(header.Value),
			Always: header.Always,
			Status: header.Status,
		})
	}
	return result
}

// 屏蔽一个Header
func (this *LocationConfig) AddIgnoreHeader(name string) {
	this.IgnoreHeaders = append(this.IgnoreHeaders, name)
}

// 移除对Header的屏蔽
func (this *LocationConfig) DeleteIgnoreHeaderAtIndex(index int) {
	if index >= 0 && index < len(this.IgnoreHeaders) {
		this.IgnoreHeaders = lists.Remove(this.IgnoreHeaders, index).([]string)
	}
}

// 更改Header的屏蔽
func (this *LocationConfig) UpdateIgnoreHeaderAtIndex(index int, name string) {
	if index >= 0 && index < len(this.IgnoreHeaders) {
		this.IgnoreHeaders[index] = name
	}
}

// 添加重写规则
func (this *LocationConfig) AddRewriteRule(rewriteRule *RewriteRule) {
	this.Rewrite = append(this.Rewrite, rewriteRule)
}

// 缓存策略
func (this *LocationConfig) CachePolicyObject() *shared.CachePolicy {
	return this.cachePolicy
}
