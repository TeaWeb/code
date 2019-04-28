package rules

import (
	"errors"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teawaf/checkpoints"
	"github.com/TeaWeb/code/teawaf/requests"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"reflect"
	"regexp"
	"strings"
)

var singleParamRegexp = regexp.MustCompile("^\\${[\\w.-]+}$")

// rule
type Rule struct {
	Param             string            `yaml:"param" json:"param"`       // such as ${arg.name} or ${args}, can be composite as ${arg.firstName}${arg.lastName}
	Operator          RuleOperator      `yaml:"operator" json:"operator"` // such as contains, gt,  ...
	Value             string            `yaml:"value" json:"value"`       // compared value
	IsCaseInsensitive bool              `yaml:"isCaseInsensitive" json:"isCaseInsensitive"`
	CheckpointOptions map[string]string `yaml:"checkpointOptions" json:"checkpointOptions"`

	checkpointFinder func(prefix string) checkpoints.CheckpointInterface

	singleParam      string                          // real param after prefix
	singleCheckpoint checkpoints.CheckpointInterface // if is single check point

	multipleCheckpoints map[string]checkpoints.CheckpointInterface

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
		v := this.Value
		if this.IsCaseInsensitive && !strings.HasPrefix(v, "(?i)") {
			v = "(?i)" + v
		}

		v = this.unescape(v)

		reg, err := regexp.Compile(v)
		if err != nil {
			return err
		}
		this.reg = reg
	case RuleOperatorNotMatch:
		v := this.Value
		if this.IsCaseInsensitive && !strings.HasPrefix(v, "(?i)") {
			v = "(?i)" + v
		}

		v = this.unescape(v)

		reg, err := regexp.Compile(v)
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

		if this.checkpointFinder != nil {
			checkpoint := this.checkpointFinder(prefix)
			if checkpoint == nil {
				return errors.New("no check point '" + prefix + "' found")
			}
			this.singleCheckpoint = checkpoint
		} else {
			checkpoint := checkpoints.FindCheckpoint(prefix)
			if checkpoint == nil {
				return errors.New("no check point '" + prefix + "' found")
			}
			checkpoint.Init()
			this.singleCheckpoint = checkpoint
		}

		return nil
	}

	this.multipleCheckpoints = map[string]checkpoints.CheckpointInterface{}
	var err error = nil
	teautils.ParseVariables(this.Param, func(varName string) (value string) {
		pieces := strings.SplitN(varName, ".", 2)
		prefix := pieces[0]
		if this.checkpointFinder != nil {
			checkpoint := this.checkpointFinder(prefix)
			if checkpoint == nil {
				err = errors.New("no check point '" + prefix + "' found")
			} else {
				this.multipleCheckpoints[prefix] = checkpoint
			}
		} else {
			checkpoint := checkpoints.FindCheckpoint(prefix)
			if checkpoint == nil {
				err = errors.New("no check point '" + prefix + "' found")
			} else {
				checkpoint.Init()
				this.multipleCheckpoints[prefix] = checkpoint
			}
		}
		return ""
	})

	return err
}

func (this *Rule) MatchRequest(req *requests.Request) (b bool, err error) {
	if this.singleCheckpoint != nil {
		value, err, _ := this.singleCheckpoint.RequestValue(req, this.singleParam, this.CheckpointOptions)
		if err != nil {
			return false, err
		}
		return this.Test(value), nil
	}

	value := teautils.ParseVariables(this.Param, func(varName string) (value string) {
		pieces := strings.SplitN(varName, ".", 2)
		prefix := pieces[0]
		point, ok := this.multipleCheckpoints[prefix]
		if !ok {
			return ""
		}

		if len(pieces) == 1 {
			value1, err1, _ := point.RequestValue(req, "", this.CheckpointOptions)
			if err1 != nil {
				err = err1
			}
			return types.String(value1)
		}

		value1, err1, _ := point.RequestValue(req, pieces[1], this.CheckpointOptions)
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

func (this *Rule) MatchResponse(req *requests.Request, resp *requests.Response) (b bool, err error) {
	if this.singleCheckpoint != nil {
		// if is request param
		if this.singleCheckpoint.IsRequest() {
			value, err, _ := this.singleCheckpoint.RequestValue(req, this.singleParam, this.CheckpointOptions)
			if err != nil {
				return false, err
			}
			return this.Test(value), nil
		}

		// response param
		value, err, _ := this.singleCheckpoint.ResponseValue(req, resp, this.singleParam, this.CheckpointOptions)
		if err != nil {
			return false, err
		}
		return this.Test(value), nil
	}

	value := teautils.ParseVariables(this.Param, func(varName string) (value string) {
		pieces := strings.SplitN(varName, ".", 2)
		prefix := pieces[0]
		point, ok := this.multipleCheckpoints[prefix]
		if !ok {
			return ""
		}

		if len(pieces) == 1 {
			if point.IsRequest() {
				value1, err1, _ := point.RequestValue(req, "", this.CheckpointOptions)
				if err1 != nil {
					err = err1
				}
				return types.String(value1)
			} else {
				value1, err1, _ := point.ResponseValue(req, resp, "", this.CheckpointOptions)
				if err1 != nil {
					err = err1
				}
				return types.String(value1)
			}
		}

		if point.IsRequest() {
			value1, err1, _ := point.RequestValue(req, pieces[1], this.CheckpointOptions)
			if err1 != nil {
				err = err1
			}
			return types.String(value1)
		} else {
			value1, err1, _ := point.ResponseValue(req, resp, pieces[1], this.CheckpointOptions)
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
		if this.IsCaseInsensitive {
			return strings.ToLower(types.String(value)) == strings.ToLower(this.Value)
		} else {
			return types.String(value) == this.Value
		}
	case RuleOperatorNeqString:
		if this.IsCaseInsensitive {
			return strings.ToLower(types.String(value)) != strings.ToLower(this.Value)
		} else {
			return types.String(value) != this.Value
		}
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
			lowerValue := ""
			if this.IsCaseInsensitive {
				lowerValue = strings.ToLower(this.Value)
			}
			for _, v := range maps.NewMap(value) {
				if this.IsCaseInsensitive {
					if strings.ToLower(types.String(v)) == lowerValue {
						return true
					}
				} else {
					if types.String(v) == this.Value {
						return true
					}
				}
			}
			return false
		}

		if this.IsCaseInsensitive {
			return strings.Contains(strings.ToLower(types.String(value)), strings.ToLower(this.Value))
		} else {
			return strings.Contains(types.String(value), this.Value)
		}
	case RuleOperatorNotContains:
		if this.IsCaseInsensitive {
			return !strings.Contains(strings.ToLower(types.String(value)), strings.ToLower(this.Value))
		} else {
			return !strings.Contains(types.String(value), this.Value)
		}
	case RuleOperatorPrefix:
		if this.IsCaseInsensitive {
			return strings.HasPrefix(strings.ToLower(types.String(value)), strings.ToLower(this.Value))
		} else {
			return strings.HasPrefix(types.String(value), this.Value)
		}
	case RuleOperatorSuffix:
		if this.IsCaseInsensitive {
			return strings.HasSuffix(strings.ToLower(types.String(value)), strings.ToLower(this.Value))
		} else {
			return strings.HasSuffix(types.String(value), this.Value)
		}
	case RuleOperatorHasKey:
		if types.IsSlice(value) {
			index := types.Int(this.Value)
			if index < 0 {
				return false
			}
			return reflect.ValueOf(value).Len() > index
		} else if types.IsMap(value) {
			m := maps.NewMap(value)
			if this.IsCaseInsensitive {
				lowerValue := strings.ToLower(this.Value)
				for k := range m {
					if strings.ToLower(k) == lowerValue {
						return true
					}
				}
			} else {
				return m.Has(this.Value)
			}
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

func (this *Rule) IsSingleCheckpoint() bool {
	return this.singleCheckpoint != nil
}

func (this *Rule) SetCheckpointFinder(finder func(prefix string) checkpoints.CheckpointInterface) {
	this.checkpointFinder = finder
}

func (this *Rule) unescape(v string) string {
	//replace urlencoded characters
	v = strings.Replace(v, `\s`, `(\s|%09|%0A|\+)`, -1)
	v = strings.Replace(v, `\(`, `(\(|%28)`, -1)
	v = strings.Replace(v, `=`, `(=|%3D)`, -1)
	v = strings.Replace(v, `<`, `(<|%3C)`, -1)
	v = strings.Replace(v, `\*`, `(\*|%2A)`, -1)
	v = strings.Replace(v, `\\`, `(\\|%2F)`, -1)
	v = strings.Replace(v, `!`, `(!|%21)`, -1)
	v = strings.Replace(v, `/`, `(/|%2F)`, -1)
	v = strings.Replace(v, `;`, `(;|%3B)`, -1)
	return v
}
