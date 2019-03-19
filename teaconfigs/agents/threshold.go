package agents

import (
	"errors"
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/robertkrimen/otto"
	"regexp"
	"strings"
)

// 阈值定义
type Threshold struct {
	Id            string                   `yaml:"id" json:"id"`                       // ID
	Param         string                   `yaml:"param" json:"param"`                 // 参数
	Operator      ThresholdOperator        `yaml:"operator" json:"operator"`           // 运算符
	Value         string                   `yaml:"value" json:"value"`                 // 对比值
	NoticeLevel   notices.NoticeLevel      `yaml:"noticeLevel" json:"noticeLevel"`     // 通知级别
	NoticeMessage string                   `yaml:"noticeMessage" json:"noticeMessage"` // 通知消息
	Actions       []map[string]interface{} `yaml:"actions" json:"actions"`             // 动作配置
	MaxFails      int                      `yaml:"maxFails" json:"maxFails"`           // 连续失败次数

	regValue   *regexp.Regexp
	floatValue float64
}

// 新阈值对象
func NewThreshold() *Threshold {
	return &Threshold{
		Id: stringutil.Rand(16),
	}
}

// 校验
func (this *Threshold) Validate() error {
	if this.Operator == ThresholdOperatorRegexp || this.Operator == ThresholdOperatorNotRegexp {
		reg, err := regexp.Compile(this.Value)
		if err != nil {
			return err
		}
		this.regValue = reg
	} else if this.Operator == ThresholdOperatorGt || this.Operator == ThresholdOperatorGte || this.Operator == ThresholdOperatorLt || this.Operator == ThresholdOperatorLte {
		this.floatValue = types.Float64(this.Value)
	}

	return nil
}

// 将此条件应用于阈值，检查是否匹配
func (this *Threshold) Test(value interface{}, oldValue interface{}) bool {
	paramValue := this.Eval(value, oldValue)

	switch this.Operator {
	case ThresholdOperatorRegexp:
		if this.regValue == nil {
			return false
		}
		return this.regValue.MatchString(types.String(paramValue))
	case ThresholdOperatorNotRegexp:
		if this.regValue == nil {
			return false
		}
		return !this.regValue.MatchString(types.String(paramValue))
	case ThresholdOperatorGt:
		return types.Float64(paramValue) > this.floatValue
	case ThresholdOperatorGte:
		return types.Float64(paramValue) >= this.floatValue
	case ThresholdOperatorLt:
		return types.Float64(paramValue) < this.floatValue
	case ThresholdOperatorLte:
		return types.Float64(paramValue) <= this.floatValue
	case ThresholdOperatorEq:
		return paramValue == this.Value
	case ThresholdOperatorNot:
		return paramValue != this.Value
	case ThresholdOperatorPrefix:
		return strings.HasPrefix(types.String(paramValue), this.Value)
	case ThresholdOperatorSuffix:
		return strings.HasSuffix(types.String(paramValue), this.Value)
	case ThresholdOperatorContains:
		return strings.Contains(types.String(paramValue), this.Value)
	case ThresholdOperatorNotContains:
		return !strings.Contains(types.String(paramValue), this.Value)
	}
	return false
}

// 执行数值运算，使用Javascript语法
func (this *Threshold) Eval(value interface{}, old interface{}) string {
	return this.EvalParam(this.Param, value, old)
}

// 使用某个参数执行数值运算，使用Javascript语法
func (this *Threshold) EvalParam(param string, value interface{}, old interface{}) string {
	if old == nil {
		old = value
	}
	paramValue := teaconfigs.RegexpNamedVariable.ReplaceAllStringFunc(param, func(s string) string {
		if value == nil {
			return ""
		}

		varName := s[2 : len(s)-1]

		// 支持${OLD}和${OLD.xxx}
		if varName == "OLD" {
			return this.EvalParam("${0}", old, nil)
		} else if strings.HasPrefix(varName, "OLD.") {
			return this.EvalParam("${"+varName[4:]+"}", old, nil)
		}

		switch v := value.(type) {
		case string:
			if varName == "0" {
				return v
			}
			return ""
		case int8, int16, int, int32, int64, uint8, uint16, uint, uint32, uint64:
			if varName == "0" {
				return fmt.Sprintf("%d", v)
			}
			return "0"
		case float32, float64:
			if varName == "0" {
				return fmt.Sprintf("%f", v)
			}
			return "0"
		case bool:
			if varName == "0" {
				if v {
					return "1"
				}
				return "0"
			}
			return "0"
		case []interface{}:
			index := types.Int(varName)
			if index >= 0 && index < len(v) {
				return types.String(v[index])
			}
			return ""
		case map[string]interface{}:
			result, found := v[varName]
			if found {
				return types.String(result)
			}
			result = teautils.Get(v, strings.Split(varName, "."))
			if result == nil {
				return ""
			}
			return types.String(result)
		case maps.Map:
			result, found := v[varName]
			if found {
				return types.String(result)
			}
			result = teautils.Get(v, strings.Split(varName, "."))
			if result == nil {
				return ""
			}
			return types.String(result)
		}
		return s
	})

	// 支持加、减、乘、除、余
	if len(paramValue) > 0 {
		if strings.ContainsAny(paramValue, "+-*/%") {
			vm := otto.New()
			v, err := vm.Run(paramValue)
			if err != nil {
				logs.Error(errors.New("\"" + this.Expression() + "\": eval \"" + paramValue + "\":" + err.Error()))
			} else {
				paramValue = v.String()
			}
		}
	}

	return paramValue
}

// 执行动作
func (this *Threshold) RunActions(params map[string]string) error {
	if len(this.Actions) == 0 {
		return nil
	}

	for _, a := range this.Actions {
		code, found := a["code"]
		if !found {
			return errors.New("action 'code' not found")
		}

		options, found := a["options"]
		if !found {
			return errors.New("action 'options' not found")
		}
		optionsMap, ok := options.(map[string]interface{})
		if !ok {
			return errors.New("action 'options' should be a valid map")
		}

		action := FindAction(types.String(code))
		if action == nil {
			return errors.New("action for '" + types.String(code) + "' not found")
		}

		instance := action["instance"]
		err := teautils.MapToObjectJSON(optionsMap, &instance)
		if err != nil {
			return err
		}

		output, err := instance.(ActionInterface).Run(params)
		if err != nil {
			return err
		}
		if len(output) > 0 {
			logs.Println("[threshold]run actions:", output)
		}
	}

	return nil
}

// 取得描述文本
func (this *Threshold) Expression() string {
	return this.Param + " " + this.Operator + " " + this.Value
}
