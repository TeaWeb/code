package shared

import "regexp"

// 常用的正则表达式
var (
	RegexpDigitNumber   = regexp.MustCompile(`^\d+$`)                    // 整数
	RegexpFloatNumber   = regexp.MustCompile(`^\d+(\.\d+)?$`)            // 浮点数，只支持xxx.xxx
	RegexpExternalURL   = regexp.MustCompile("(?i)^(http|https|ftp)://") // URL
	RegexpNamedVariable = regexp.MustCompile("\\${[\\w.-]+}")            // 命名变量
)
