package rules

import (
	"errors"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teawaf/checkpoints"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

var singleParamRegexp = regexp.MustCompile("^\\${[\\w.-]+}$")

// rule
type Rule struct {
	Param    string       `yaml:"param" json:"param"`       // such as ${arg.name} or ${args}, can be composite as ${arg.firstName}${arg.lastName}
	Operator RuleOperator `yaml:"operator" json:"operator"` // such as contains, gt,  ...
	Value    string       `yaml:"value" json:"value"`       // compared value

	singleParam      string                          // real param after prefix
	singleCheckPoint checkpoints.CheckPointInterface // if is single check point

	multipleCheckPoints map[string]checkpoints.CheckPointInterface

	floatValue float64
	reg        *regexp.Regexp
}

func NewRule() *Rule {
	return &Rule{}
}

func (this *Rule) Init() error {
	// operator
	switch this.Operator {
	case RuleOperatorGt:
		this.floatValue = types.Float64(this.Value)
	case RuleOperatorGte:
		this.floatValue = types.Float64(this.Value)
	case RuleOperatorLt:
		this.floatValue = types.Float64(this.Value)
	case RuleOperatorLte:
		this.floatValue = types.Float64(this.Value)
	case RuleOperatorEq:
		this.floatValue = types.Float64(this.Value)
	case RuleOperatorNeq:
		this.floatValue = types.Float64(this.Value)
	case RuleOperatorMatch:
		reg, err := regexp.Compile(this.Value)
		if err != nil {
			return err
		}
		this.reg = reg
	case RuleOperatorNotMatch:
		reg, err := regexp.Compile(this.Value)
		if err != nil {
			return err
		}
		this.reg = reg
	}

	if singleParamRegexp.MatchString(this.Param) {
		param := this.Param[2 : len(this.Param)-1]
		pieces := strings.SplitN(param, ".", 2)
		prefix := pieces[0]
		if len(pieces) == 1 {
			this.singleParam = ""
		} else {
			this.singleParam = pieces[1]
		}

		point := checkpoints.FindCheckPoint(prefix)
		if point == nil {
			return errors.New("no check point '" + prefix + "' found")
		}
		this.singleCheckPoint = point

		return nil
	}

	this.multipleCheckPoints = map[string]checkpoints.CheckPointInterface{}
	var err error = nil
	teautils.ParseVariables(this.Param, func(varName string) (value string) {
		pieces := strings.SplitN(varName, ".", 2)
		prefix := pieces[0]
		checkPoint := checkpoints.FindCheckPoint(prefix)
		if checkPoint == nil {
			err = errors.New("no check point '" + prefix + "' found")
		} else {
			this.multipleCheckPoints[prefix] = checkPoint
		}
		return ""
	})

	return err
}

func (this *Rule) MatchRequest(req *http.Request) (b bool, err error) {
	if this.singleCheckPoint != nil {
		value, err := this.singleCheckPoint.RequestValue(req, this.singleParam)
		if err != nil {
			return false, err
		}
		return this.Test(value), nil
	}

	value := teautils.ParseVariables(this.Param, func(varName string) (value string) {
		pieces := strings.SplitN(varName, ".", 2)
		prefix := pieces[0]
		point, ok := this.multipleCheckPoints[prefix]
		if !ok {
			return ""
		}

		if len(pieces) == 1 {
			value1, err1 := point.RequestValue(req, "")
			if err1 != nil {
				err = err1
			}
			return types.String(value1)
		}

		value1, err1 := point.RequestValue(req, pieces[1])
		if err1 != nil {
			err = err1
		}
		return types.String(value1)
	})

	if err != nil {
		return false, err
	}

	return this.Test(value), nil
}

func (this *Rule) MatchResponse(req *http.Request, resp *http.Response) (b bool, err error) {
	if this.singleCheckPoint != nil {
		// if is request param
		if this.singleCheckPoint.IsRequest() {
			value, err := this.singleCheckPoint.RequestValue(req, this.singleParam)
			if err != nil {
				return false, err
			}
			return this.Test(value), nil
		}

		// response param
		value, err := this.singleCheckPoint.ResponseValue(req, resp, this.singleParam)
		if err != nil {
			return false, err
		}
		return this.Test(value), nil
	}

	value := teautils.ParseVariables(this.Param, func(varName string) (value string) {
		pieces := strings.SplitN(varName, ".", 2)
		prefix := pieces[0]
		point, ok := this.multipleCheckPoints[prefix]
		if !ok {
			return ""
		}

		if len(pieces) == 1 {
			if point.IsRequest() {
				value1, err1 := point.RequestValue(req, "")
				if err1 != nil {
					err = err1
				}
				return types.String(value1)
			} else {
				value1, err1 := point.ResponseValue(req, resp, "")
				if err1 != nil {
					err = err1
				}
				return types.String(value1)
			}
		}

		if point.IsRequest() {
			value1, err1 := point.RequestValue(req, pieces[1])
			if err1 != nil {
				err = err1
			}
			return types.String(value1)
		} else {
			value1, err1 := point.ResponseValue(req, resp, pieces[1])
			if err1 != nil {
				err = err1
			}
			return types.String(value1)
		}
	})

	if err != nil {
		return false, err
	}

	return this.Test(value), nil
}

func (this *Rule) Test(value interface{}) bool {
	// operator
	switch this.Operator {
	case RuleOperatorGt:
		return types.Float64(value) > this.floatValue
	case RuleOperatorGte:
		return types.Float64(value) >= this.floatValue
	case RuleOperatorLt:
		return types.Float64(value) < this.floatValue
	case RuleOperatorLte:
		return types.Float64(value) <= this.floatValue
	case RuleOperatorEq:
		return types.Float64(value) == this.floatValue
	case RuleOperatorNeq:
		return types.Float64(value) != this.floatValue
	case RuleOperatorEqString:
		return types.String(value) == this.Value
	case RuleOperatorNeqString:
		return types.String(value) != this.Value
	case RuleOperatorMatch:
		return this.reg.MatchString(types.String(value))
	case RuleOperatorNotMatch:
		return !this.reg.MatchString(types.String(value))
	case RuleOperatorContains:
		if types.IsSlice(value) {
			ok := false
			lists.Each(value, func(k int, v interface{}) {
				if types.String(v) == this.Value {
					ok = true
				}
			})
			return ok
		}
		if types.IsMap(value) {
			for _, v := range maps.NewMap(value) {
				if types.String(v) == this.Value {
					return true
				}
			}
			return false
		}
		return strings.Contains(types.String(value), this.Value)
	case RuleOperatorNotContains:
		return !strings.Contains(types.String(value), this.Value)
	case RuleOperatorPrefix:
		return strings.HasPrefix(types.String(value), this.Value)
	case RuleOperatorSuffix:
		return strings.HasSuffix(types.String(value), this.Value)
	case RuleOperatorHasKey:
		if types.IsSlice(value) {
			index := types.Int(this.Value)
			if index < 0 {
				return false
			}
			return reflect.ValueOf(value).Len() > index
		} else if types.IsMap(value) {
			m := maps.NewMap(value)
			return m.Has(this.Value)
		} else {
			return false
		}

	case RuleOperatorVersionGt:
		return stringutil.VersionCompare(this.Value, types.String(value)) > 0
	case RuleOperatorVersionLt:
		return stringutil.VersionCompare(this.Value, types.String(value)) < 0
	}
	return false
}

func (this *Rule) IsSingleCheckPoint() bool {
	return this.singleCheckPoint != nil
}
