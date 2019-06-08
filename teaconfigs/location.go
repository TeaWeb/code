package teaconfigs

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teawaf"
	"github.com/iwind/TeaGo/utils/string"
	"net/http"
	"regexp"
	"strings"
)

// 路径配置
type LocationConfig struct {
	shared.HeaderList `yaml:",inline"`
	FastcgiList       `yaml:",inline"`
	RewriteList       `yaml:",inline"`
	BackendList       `yaml:",inline"`

	On      bool   `yaml:"on" json:"on"`           // 是否开启
	Id      string `yaml:"id" json:"id"`           // ID
	Name    string `yaml:"name" json:"name"`       // 名称
	Pattern string `yaml:"pattern" json:"pattern"` // 匹配规则

	Async           bool                 `yaml:"async" json:"async"`                     // 是否异步请求 @TODO
	Notify          []interface{}        `yaml:"notify" json:"notify"`                   // 转发请求，可以配置转发策略 @TODO
	LogOnly         bool                 `yaml:"logOnly" json:"logOnly"`                 // 是否只记录日志 @TODO
	Root            string               `yaml:"root" json:"root"`                       // 资源根目录
	Index           []string             `yaml:"index" json:"index"`                     // 默认文件
	Charset         string               `yaml:"charset" json:"charset"`                 // 字符集设置
	MaxBodySize     string               `yaml:"maxBodySize" json:"maxBodySize"`         // 请求body最大尺寸
	GzipLevel       int8                 `yaml:"gzipLevel" json:"gzipLevel"`             // Gzip压缩级别
	GzipMinLength   string               `yaml:"gzipMinLength" json:"gzipMinLength"`     // 需要压缩的最小内容尺寸
	AccessPolicy    *shared.AccessPolicy `yaml:"accessPolicy" json:"accessPolicy"`       // 访问控制
	RedirectToHttps bool                 `yaml:"redirectToHttps" json:"redirectToHttps"` // 是否自动跳转到Https

	// 日志
	AccessLog []*AccessLogConfig `yaml:"accessLog" json:"accessLog"` // 访问日志设置，如果为空表示继承上一级设置

	DisableAccessLog1 bool  `yaml:"disableAccessLog" json:"disableAccessLog"` // deprecated: 是否禁用访问日志
	AccessLogFields1  []int `yaml:"accessLogFields" json:"accessLogFields"`   // deprecated: 访问日志保留的字段，如果为nil，则表示没有设置
	DisableStat       bool  `yaml:"disableStat" json:"disableStat"`           // 是否禁用统计

	// 参考：http://nginx.org/en/docs/http/ngx_http_access_module.html
	Allow []string `yaml:"allow" json:"allow"` // 允许的终端地址 @TODO
	Deny  []string `yaml:"deny" json:"deny"`   // 禁止的终端地址 @TODO

	Proxy string `yaml:"proxy" json:"proxy"` //  代理配置 @TODO

	CachePolicy string `yaml:"cachePolicy" json:"cachePolicy"` // 缓存策略
	CacheOn     bool   `yaml:"cacheOn" json:"cacheOn"`         // 缓存是否打开
	cachePolicy *shared.CachePolicy

	WAFOn bool   `yaml:"wafOn" json:"wafOn"` // 是否启用
	WafId string `yaml:"wafId" json:"wafId"` // WAF ID
	waf   *teawaf.WAF                        // waf object

	// websocket设置
	Websocket *WebsocketConfig `yaml:"websocket" json:"websocket"`

	// 开启的条件
	// 语法为：cond param operator value 比如：
	// - cond ${status} gte 200
	// - cond ${arg.name} eq lily
	// - cond ${requestPath} regexp .*\.png
	Cond []*RequestCond `yaml:"cond" json:"cond"`

	// 请求分组（从server复制而来）
	requestGroups          []*RequestGroup
	defaultRequestGroup    *RequestGroup
	hasRequestGroupFilters bool

	maxBodySize   int64
	gzipMinLength int64

	patternType LocationPatternType // 规则类型：LocationPattern*
	prefix      string              // 前缀
	path        string              // 精确的路径

	reg             *regexp.Regexp // 匹配规则
	caseInsensitive bool           // 大小写不敏感
	reverse         bool           // 是否翻转规则，比如非前缀，非路径
}

// 获取新对象
func NewLocation() *LocationConfig {
	return &LocationConfig{
		On:      true,
		Id:      stringutil.Rand(16),
		CacheOn: true,
		WAFOn:   true,
	}
}

// 校验
func (this *LocationConfig) Validate() error {
	// 最大Body尺寸
	maxBodySize, _ := stringutil.ParseFileSize(this.MaxBodySize)
	this.maxBodySize = int64(maxBodySize)

	gzipMinLength, _ := stringutil.ParseFileSize(this.GzipMinLength)
	this.gzipMinLength = int64(gzipMinLength)

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

	// waf
	if len(this.WafId) > 0 && this.WAFOn {
		waf := SharedWAFList().FindWAF(this.WafId)
		if waf != nil {
			err := waf.Init()
			if err != nil {
				return err
			}
			this.waf = waf
		}
	}

	// 校验访问控制
	if this.AccessPolicy != nil {
		err := this.AccessPolicy.Validate()
		if err != nil {
			return err
		}
	}

	// 校验RewriteRule配置
	err := this.ValidateRewriteRules()
	if err != nil {
		return err
	}

	// 校验Fastcgi配置
	err = this.ValidateFastcgi()
	if err != nil {
		return err
	}

	// 校验Header
	err = this.ValidateHeaders()
	if err != nil {
		return err
	}

	//websocket
	if this.Websocket != nil {
		err = this.Websocket.Validate()
		if err != nil {
			return err
		}
	}

	// backend
	err = this.ValidateBackends()
	if err != nil {
		return err
	}

	// 校验条件
	for _, cond := range this.Cond {
		err := cond.Validate()
		if err != nil {
			return err
		}
	}

	// request groups
	for _, group := range this.requestGroups {
		group.Backends = []*BackendConfig{}
		group.Scheduling = this.Scheduling

		if group.IsDefault {
			this.defaultRequestGroup = group
		}

		for _, backend := range this.Backends {
			if len(backend.RequestGroupIds) == 0 && group.Id == "default" {
				group.AddBackend(backend)
			} else if backend.HasRequestGroupId(group.Id) {
				group.AddBackend(backend)
			}
		}

		err := group.Validate()
		if err != nil {
			return err
		}
		if group.HasFilters() {
			this.hasRequestGroupFilters = true
		}
	}

	return nil
}

// 兼容性设置
func (this *LocationConfig) Compatible(version string) {
	if len(version) == 0 {
		this.CacheOn = true
		this.WAFOn = true
	} else if stringutil.VersionCompare(version, "0.1.3") < 0 {
		this.CacheOn = true
		this.WAFOn = true
	} else if stringutil.VersionCompare(version, "0.1.5") <= 0 {
		if len(this.AccessLog) == 0 && this.DisableAccessLog1 {
			this.AccessLog = []*AccessLogConfig{
				{
					Id:      stringutil.Rand(16),
					On:      !this.DisableAccessLog1,
					Fields:  this.AccessLogFields1,
					Status1: true,
					Status2: true,
					Status3: true,
					Status4: true,
					Status5: true,
				},
			}
		}
	}
}

// 最大Body尺寸
func (this *LocationConfig) MaxBodyBytes() int64 {
	return this.maxBodySize
}

// 可压缩最小尺寸
func (this *LocationConfig) GzipMinBytes() int64 {
	return this.gzipMinLength
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
func (this *LocationConfig) Match(path string, formatter func(source string) string) (map[string]string, bool) {
	// 判断条件
	if len(this.Cond) > 0 {
		for _, cond := range this.Cond {
			if !cond.Match(formatter) {
				return nil, false
			}
		}
	}

	if this.patternType == LocationPatternTypePrefix {
		if this.reverse {
			if this.caseInsensitive {
				return nil, !strings.HasPrefix(strings.ToLower(path), strings.ToLower(this.prefix))
			} else {
				return nil, !strings.HasPrefix(path, this.prefix)
			}
		} else {
			if this.caseInsensitive {
				return nil, strings.HasPrefix(strings.ToLower(path), strings.ToLower(this.prefix))
			} else {
				return nil, strings.HasPrefix(path, this.prefix)
			}
		}
	}

	if this.patternType == LocationPatternTypeExact {
		if this.reverse {
			if this.caseInsensitive {
				return nil, strings.ToLower(path) != strings.ToLower(this.path)
			} else {
				return nil, path != this.path
			}
		} else {
			if this.caseInsensitive {
				return nil, strings.ToLower(path) == strings.ToLower(this.path)
			} else {
				return nil, path == this.path
			}
		}
	}

	if this.patternType == LocationPatternTypeRegexp {
		if this.reg != nil {
			if this.reverse {
				return nil, !this.reg.MatchString(path)
			} else {
				b := this.reg.MatchString(path)
				if b {
					result := map[string]string{}
					matches := this.reg.FindStringSubmatch(path)
					subNames := this.reg.SubexpNames()
					for index, value := range matches {
						result[fmt.Sprintf("%d", index)] = value
						subName := subNames[index]
						if len(subName) > 0 {
							result[subName] = value
						}
					}
					return result, true
				}
				return nil, b
			}
		}

		return nil, this.reverse
	}

	return nil, false
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

// 缓存策略
func (this *LocationConfig) CachePolicyObject() *shared.CachePolicy {
	return this.cachePolicy
}

// WAF
func (this *LocationConfig) WAF() *teawaf.WAF {
	return this.waf
}

// 是否在引用某个代理
func (this *LocationConfig) RefersProxy(proxyId string) bool {
	if this.Proxy == proxyId {
		return true
	}
	for _, r := range this.Rewrite {
		if r.RefersProxy(proxyId) {
			return true
		}
	}
	return false
}

// 添加过滤条件
func (this *LocationConfig) AddCond(cond *RequestCond) {
	this.Cond = append(this.Cond, cond)
}

// 添加请求分组
func (this *LocationConfig) AddRequestGroup(group *RequestGroup) {
	this.requestGroups = append(this.requestGroups, group)

	if this.Websocket != nil {
		this.Websocket.AddRequestGroup(group.Copy())
	}
}

// 使用请求匹配分组
func (this *LocationConfig) MatchRequestGroup(formatter func(source string) string) *RequestGroup {
	if !this.hasRequestGroupFilters {
		return nil
	}
	for _, group := range this.requestGroups {
		if group.HasFilters() && group.Match(formatter) {
			return group
		}
	}
	return nil
}

// 取得下一个可用的后端服务
func (this *LocationConfig) NextBackend(call *shared.RequestCall) *BackendConfig {
	if this.hasRequestGroupFilters {
		group := this.MatchRequestGroup(call.Formatter)
		if group != nil {
			// request
			if group.HasRequestHeaders() {
				for _, h := range group.RequestHeaders {
					call.Request.Header.Set(h.Name, call.Formatter(h.Value))
				}
			}

			// response
			if group.HasResponseHeaders() {
				call.AddResponseCall(func(resp http.ResponseWriter) {
					for _, h := range group.ResponseHeaders {
						resp.Header().Set(h.Name, call.Formatter(h.Value))
					}
				})
			}

			return group.BackendList.NextBackend(call)
		}
	}

	// 默认分组
	if this.defaultRequestGroup != nil {
		// request
		if this.defaultRequestGroup.HasRequestHeaders() {
			for _, h := range this.defaultRequestGroup.RequestHeaders {
				call.Request.Header.Set(h.Name, call.Formatter(h.Value))
			}
		}

		// response
		if this.defaultRequestGroup.HasResponseHeaders() {
			call.AddResponseCall(func(resp http.ResponseWriter) {
				for _, h := range this.defaultRequestGroup.ResponseHeaders {
					resp.Header().Set(h.Name, call.Formatter(h.Value))
				}
			})
		}

		return this.defaultRequestGroup.NextBackend(call)
	}

	return this.BackendList.NextBackend(call)
}

// 设置调度算法
func (this *LocationConfig) SetupScheduling(isBackup bool) {
	for _, group := range this.requestGroups {
		group.SetupScheduling(isBackup)
	}
	this.BackendList.SetupScheduling(isBackup)
}

// 装载事件
func (this *LocationConfig) OnAttach() {
	// 开启WAF
	if this.waf != nil {
		this.waf.Start()
	}
}

// 卸载事件
func (this *LocationConfig) OnDetach() {
	// 停止WAF
	if this.waf != nil {
		this.waf.Stop()
		this.waf = nil
	}
}
