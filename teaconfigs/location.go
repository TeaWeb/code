package teaconfigs

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/iwind/TeaGo/utils/string"
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
	Pattern string `yaml:"pattern" json:"pattern"` // 匹配规则

	Async         bool                 `yaml:"async" json:"async"`                 // 是否异步请求 @TODO
	Notify        []interface{}        `yaml:"notify" json:"notify"`               // 转发请求，可以配置转发策略 @TODO
	LogOnly       bool                 `yaml:"logOnly" json:"logOnly"`             // 是否只记录日志 @TODO
	Root          string               `yaml:"root" json:"root"`                   // 资源根目录
	Index         []string             `yaml:"index" json:"index"`                 // 默认文件
	Charset       string               `yaml:"charset" json:"charset"`             // 字符集设置
	MaxBodySize   string               `yaml:"maxBodySize" json:"maxBodySize"`     // 请求body最大尺寸
	GzipLevel     int8                 `yaml:"gzipLevel" json:"gzipLevel"`         // Gzip压缩级别
	GzipMinLength string               `yaml:"gzipMinLength" json:"gzipMinLength"` // 需要压缩的最小内容尺寸
	AccessPolicy  *shared.AccessPolicy `yaml:"accessPolicy" json:"accessPolicy"`   // 访问控制

	// 日志
	DisableAccessLog bool               `yaml:"disableAccessLog" json:"disableAccessLog"` // 是否禁用访问日志
	AccessLog        []*AccessLogConfig `yaml:"accessLog" json:"accessLog"`               // 访问日志设置 TODO

	// 参考：http://nginx.org/en/docs/http/ngx_http_access_module.html
	Allow []string `yaml:"allow" json:"allow"` // 允许的终端地址 @TODO
	Deny  []string `yaml:"deny" json:"deny"`   // 禁止的终端地址 @TODO

	Proxy string `yaml:proxy" json:"proxy"` //  代理配置 @TODO

	CachePolicy string `yaml:"cachePolicy" json:"cachePolicy"` // 缓存策略
	CacheOn     bool   `yaml:"cacheOn" json:"cacheOn"`         // 缓存是否打开 TODO
	cachePolicy *shared.CachePolicy

	// websocket设置
	Websocket *WebsocketConfig `yaml:"websocket" json:"websocket"`

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
		On: true,
		Id: stringutil.Rand(16),
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

	return nil
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
func (this *LocationConfig) Match(path string) (map[string]string, bool) {
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
