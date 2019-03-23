package teastats

// 数值增长型的统计值
type CounterValue struct {
	Timestamp int64                  `json:"timestamp"`
	Params    map[string]string      `json:"params"`
	Value     map[string]interface{} `json:"value"`
}
