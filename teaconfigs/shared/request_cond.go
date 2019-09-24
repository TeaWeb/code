package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"net"
	"regexp"
	"strings"
)

// 重写条件定义
type RequestCond struct {
	Id string `yaml:"id" json:"id"` // ID

	// 要测试的字符串
	// 其中可以使用跟请求相关的参数，比如：
	// ${arg.name}, ${requestPath}
	Param string `yaml:"param" json:"param"`

	// 运算符
	Operator RequestCondOperator `yaml:"operator" json:"operator"`

	// 对比
	Value string `yaml:"value" json:"value"`

	isInt   bool
	isFloat bool
	isIP    bool

	regValue   *regexp.Regexp
	floatValue float64
	ipValue    net.IP
	arrayValue []string
}

// 取得新对象
func NewRequestCond() *RequestCond {
	return &RequestCond{
		Id: stringutil.Rand(16),
	}
}

// 校验配置
func (this *RequestCond) Validate() error {
	this.isInt = RegexpDigitNumber.MatchString(this.Value)
	this.isFloat = RegexpFloatNumber.MatchString(this.Value)

	if lists.ContainsString([]string{
		RequestCondOperatorRegexp,
		RequestCondOperatorNotRegexp,
	}, this.Operator) {
		reg, err := regexp.Compile(this.Value)
		if err != nil {
			return err
		}
		this.regValue = reg
	} else if lists.ContainsString([]string{
		RequestCondOperatorEqFloat,
		RequestCondOperatorGtFloat,
		RequestCondOperatorGteFloat,
		RequestCondOperatorLtFloat,
		RequestCondOperatorLteFloat,
	}, this.Operator) {
		this.floatValue = types.Float64(this.Value)
	} else if lists.ContainsString([]string{
		RequestCondOperatorEqIP,
		RequestCondOperatorGtIP,
		RequestCondOperatorGteIP,
		RequestCondOperatorLtIP,
		RequestCondOperatorLteIP,
	}, this.Operator) {
		this.ipValue = net.ParseIP(this.Value)
		this.isIP = this.ipValue != nil

		if !this.isIP {
			return errors.New("value should be a valid ip")
		}
	} else if lists.ContainsString([]string{
		RequestCondOperatorIPInRange,
	}, this.Operator) {
		if strings.Contains(this.Value, ",") {
			ipList := strings.SplitN(this.Value, ",", 2)
			ipString1 := strings.TrimSpace(ipList[0])
			ipString2 := strings.TrimSpace(ipList[1])

			if len(ipString1) > 0 {
				ip1 := net.ParseIP(ipString1)
				if ip1 == nil {
					return errors.New("start ip is invalid")
				}
			}

			if len(ipString2) > 0 {
				ip2 := net.ParseIP(ipString2)
				if ip2 == nil {
					return errors.New("end ip is invalid")
				}
			}
		} else if strings.Contains(this.Value, "/") {
			_, _, err := net.ParseCIDR(this.Value)
			if err != nil {
				return err
			}
		} else {
			return errors.New("invalid ip range")
		}
	} else if lists.ContainsString([]string{
		RequestCondOperatorIn,
		RequestCondOperatorNotIn,
	}, this.Operator) {
		stringsValue := []string{}
		err := json.Unmarshal([]byte(this.Value), &stringsValue)
		if err != nil {
			return err
		}
		this.arrayValue = stringsValue
	}
	return nil
}

// 将此条件应用于请求，检查是否匹配
func (this *RequestCond) Match(formatter func(source string) string) bool {
	paramValue := formatter(this.Param)
	switch this.Operator {
	case RequestCondOperatorRegexp:
		if this.regValue == nil {
			return false
		}
		return this.regValue.MatchString(paramValue)
	case RequestCondOperatorNotRegexp:
		if this.regValue == nil {
			return false
		}
		return !this.regValue.MatchString(paramValue)
	case RequestCondOperatorEqInt:
		return this.isInt && paramValue == this.Value
	case RequestCondOperatorEqFloat:
		return this.isFloat && types.Float64(paramValue) == this.floatValue
	case RequestCondOperatorGtFloat:
		return this.isFloat && types.Float64(paramValue) > this.floatValue
	case RequestCondOperatorGteFloat:
		return this.isFloat && types.Float64(paramValue) >= this.floatValue
	case RequestCondOperatorLtFloat:
		return this.isFloat && types.Float64(paramValue) < this.floatValue
	case RequestCondOperatorLteFloat:
		return this.isFloat && types.Float64(paramValue) <= this.floatValue
	case RequestCondOperatorEqString:
		return paramValue == this.Value
	case RequestCondOperatorNeqString:
		return paramValue != this.Value
	case RequestCondOperatorHasPrefix:
		return strings.HasPrefix(paramValue, this.Value)
	case RequestCondOperatorHasSuffix:
		return strings.HasSuffix(paramValue, this.Value)
	case RequestCondOperatorContainsString:
		return strings.Contains(paramValue, this.Value)
	case RequestCondOperatorNotContainsString:
		return !strings.Contains(paramValue, this.Value)
	case RequestCondOperatorEqIP:
		ip := net.ParseIP(paramValue)
		if ip == nil {
			return false
		}
		return this.isIP && bytes.Compare(this.ipValue, ip) == 0
	case RequestCondOperatorGtIP:
		ip := net.ParseIP(paramValue)
		if ip == nil {
			return false
		}
		return this.isIP && bytes.Compare(ip, this.ipValue) > 0
	case RequestCondOperatorGteIP:
		ip := net.ParseIP(paramValue)
		if ip == nil {
			return false
		}
		return this.isIP && bytes.Compare(ip, this.ipValue) >= 0
	case RequestCondOperatorLtIP:
		ip := net.ParseIP(paramValue)
		if ip == nil {
			return false
		}
		return this.isIP && bytes.Compare(ip, this.ipValue) < 0
	case RequestCondOperatorLteIP:
		ip := net.ParseIP(paramValue)
		if ip == nil {
			return false
		}
		return this.isIP && bytes.Compare(ip, this.ipValue) <= 0
	case RequestCondOperatorIPInRange:
		ip := net.ParseIP(paramValue)
		if ip == nil {
			return false
		}

		// 检查IP范围格式
		if strings.Contains(this.Value, ",") {
			ipList := strings.SplitN(this.Value, ",", 2)
			ipString1 := strings.TrimSpace(ipList[0])
			ipString2 := strings.TrimSpace(ipList[1])

			if len(ipString1) > 0 {
				ip1 := net.ParseIP(ipString1)
				if ip1 == nil {
					return false
				}

				if bytes.Compare(ip, ip1) < 0 {
					return false
				}
			}

			if len(ipString2) > 0 {
				ip2 := net.ParseIP(ipString2)
				if ip2 == nil {
					return false
				}

				if bytes.Compare(ip, ip2) > 0 {
					return false
				}
			}

			return true
		} else if strings.Contains(this.Value, "/") {
			_, ipNet, err := net.ParseCIDR(this.Value)
			if err != nil {
				return false
			}
			return ipNet.Contains(ip)
		} else {
			return false
		}
	case RequestCondOperatorIn:
		return lists.ContainsString(this.arrayValue, paramValue)
	case RequestCondOperatorNotIn:
		return !lists.ContainsString(this.arrayValue, paramValue)
	}
	return false
}
