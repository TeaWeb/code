package agents

import (
	"github.com/iwind/TeaGo/maps"
	"testing"
	"time"
)

func TestThreshold_Test(t *testing.T) {
	threshold := NewThreshold()
	threshold.Param = "${0}"
	threshold.Operator = ThresholdOperatorGt
	threshold.Value = "12"
	threshold.Validate()
	t.Log(threshold.Test("123"))

	threshold.Param = "${1}"
	threshold.Operator = ThresholdOperatorGt
	threshold.Validate()
	t.Log(threshold.Test([]interface{}{1, 200, 3}))

	threshold.Param = "${host}"
	threshold.Operator = ThresholdOperatorPrefix
	threshold.Value = "127."
	threshold.Validate()
	t.Log(threshold.Test(map[string]interface{}{
		"host": "127.0.0.1",
	}))

	threshold.Param = "${data.version}"
	threshold.Operator = ThresholdOperatorEq
	threshold.Value = "1.0.25"
	threshold.Validate()
	t.Log(threshold.Test(map[string]interface{}{
		"data": maps.Map{
			"version": "1.0.25",
		},
	}))

	threshold.Param = "${data.hello.world.0}"
	threshold.Operator = ThresholdOperatorEq
	threshold.Value = "1"
	t.Log(threshold.Test(map[string]interface{}{
		"data": maps.Map{
			"version": "1.0.25",
			"hello": maps.Map{
				"world": []string{"1", "2", "3", "4", "5"},
			},
		},
	}))
}

func TestThreshold_Eval(t *testing.T) {
	threshold := NewThreshold()
	threshold.Param = "${data.hello.world.0} * 100 / ${data.hello.world.1}"
	t.Log(threshold.Eval(map[string]interface{}{
		"data": maps.Map{
			"version": "1.0.25",
			"hello": maps.Map{
				"world": []string{"1", "2", "3", "4", "5"},
			},
		},
	}))
}

func TestThreshold_Eval_Date(t *testing.T) {
	threshold := NewThreshold()
	threshold.Param = "new Date().getTime() / 1000 - ${timestamp}"
	t.Log(threshold.Eval(map[string]interface{}{
		"timestamp": time.Now().Unix() - 10,
	}))
}

func TestThreshold_RunActions(t *testing.T) {
	threshold := NewThreshold()
	threshold.Actions = []map[string]interface{}{
		{
			"code": "script",
			"options": map[string]interface{}{
				"scriptType": "path",
				"path":       "1",
			},
		},
	}
	t.Log(threshold.RunActions(nil))
}
