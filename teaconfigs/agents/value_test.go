package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestValue_AllFlatKeys(t *testing.T) {
	{
		value := NewValue()
		logs.PrintAsJSON(value.AllFlatKeys(), t)
	}

	{
		value := NewValue()
		value.Value = map[string]interface{}{
			"a": "1",
			"b": 2,
			"c": map[string]interface{}{
				"d": 3,
				"f": []string{"g", "h", "i"},
			},
			"e": []int{1, 2, 3},
		}
		logs.PrintAsJSON(value.AllFlatKeys(), t)
	}

	{
		value := NewValue()
		value.Value = []interface{}{
			map[string]interface{}{
				"a": "1",
				"b": 2,
				"c": map[string]interface{}{
					"d": 3,
					"f": []string{"g", "h", "i"},
				},
				"e": []int{1, 2, 3},
			},
			4,
		}
		logs.PrintAsJSON(value.AllFlatKeys(), t)
	}

	{

		value := NewValue()
		value.Value = nil
		logs.PrintAsJSON(value.AllFlatKeys(), t)
	}

	{

		value := NewValue()
		value.Value = map[string]interface{}{
			"a": nil,
			"b": 1,
			"c": []interface{}{2, nil},
		}
		logs.PrintAsJSON(value.AllFlatKeys(), t)
	}
}
