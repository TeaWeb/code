package agents

import "testing"

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
}
