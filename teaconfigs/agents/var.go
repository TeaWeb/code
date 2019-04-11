package agents

import (
	"errors"
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/robertkrimen/otto"
	"strings"
)

// 使用某个参数执行数值运算，使用Javascript语法
func EvalParam(param string, value interface{}, old interface{}, varMapping maps.Map, supportsMath bool) (string, error) {
	if old == nil {
		old = value
	}
	var resultErr error = nil
	paramValue := teaconfigs.RegexpNamedVariable.ReplaceAllStringFunc(param, func(s string) string {
		varName := s[2 : len(s)-1]

		// 从varMapping中查找
		if varMapping != nil {
			index := strings.Index(varName, ".")
			var firstKey = varName
			if index > 0 {
				firstKey = varName[:index]
			}
			if varMapping.Has(firstKey) {
				result := teautils.Get(varMapping, strings.Split(varName, "."))
				if result == nil {
					return ""
				}
				return types.String(result)
			}
		}

		if value == nil {
			return ""
		}

		// 支持${OLD}和${OLD.xxx}
		if varName == "OLD" {
			result, err := EvalParam("${0}", old, nil, nil, supportsMath)
			if err != nil {
				resultErr = err
			}
			return result
		} else if strings.HasPrefix(varName, "OLD.") {
			result, err := EvalParam("${"+varName[4:]+"}", old, nil, nil, supportsMath)
			if err != nil {
				resultErr = err
			}
			return result
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
		default:
			if types.IsSlice(value) || types.IsMap(value) {
				result := teautils.Get(v, strings.Split(varName, "."))
				if result == nil {
					return ""
				}
				return types.String(result)
			}
		}
		return s
	})

	// 支持加、减、乘、除、余
	if len(paramValue) > 0 {
		if supportsMath && strings.ContainsAny(param, "+-*/%") {
			vm := otto.New()
			v, err := vm.Run(paramValue)
			if err != nil {
				return "", errors.New("\"" + param + "\": eval \"" + paramValue + "\":" + err.Error())
			} else {
				paramValue = v.String()
			}
		}

		// javascript
		if strings.HasPrefix(paramValue, "javascript:") {
			vm := otto.New()
			v, err := vm.Run(paramValue[len("javascript:"):])
			if err != nil {
				return "", errors.New("\"" + param + "\": eval \"" + paramValue + "\":" + err.Error())
			} else {
				paramValue = v.String()
			}
		}
	}

	return paramValue, resultErr
}
