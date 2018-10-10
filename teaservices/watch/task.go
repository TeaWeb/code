package watch

import "time"

type Task struct {
	Id       string                 `json:"id"`
	Name     string                 `json:"name"`
	App      string                 `json:"app"`
	Info     map[string]interface{} `json:"info"`
	API      DataAPIInterface       `json:"api"`
	Timeout  float64                `json:"timeout"`
	Interval time.Duration          `json:"interval"`
}
