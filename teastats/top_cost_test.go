package teastats

import (
	"testing"
	"github.com/TeaWeb/code/tealogs"
)

func TestTopCostStat_Process(t *testing.T) {
	log := &tealogs.AccessLog{
		RequestTime: 1,
		Scheme:      "http",
		Host:        "google.com",
		RequestURI:  "/hello",
	}
	stat := new(TopCostStat)
	stat.Process(log)
}
