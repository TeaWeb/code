package agents

import (
	"github.com/iwind/TeaGo/maps"
	"testing"
)

func TestEvalParam(t *testing.T) {
	t.Log(EvalParam("This is is message, ${ITEM.name}, ${ITEM}, ${ITE}", nil, nil, maps.Map{
		"ITEM": maps.Map{
			"name": "MySQL",
		},
	}, false))
}

// 测试空格
func TestEvalParam_Spaces(t *testing.T) {
	threshold := NewThreshold()
	threshold.Param = "[${data    .	version}]"
	threshold.Operator = ThresholdOperatorEq
	threshold.Value = "1.0.25"
	err := threshold.Validate()
	if err != nil {
		t.Error(err)
	}
	t.Log(threshold.Eval(map[string]interface{}{
		"data": maps.Map{
			"version": "1.0.26",
		},
	}, nil))
	t.Log(EvalParam(threshold.Param, nil, nil, maps.Map{
		"data": maps.Map{
			"version": "1.1",
		},
	}, true))
}
