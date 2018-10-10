package teastats

import (
	"testing"
	"time"
	"github.com/iwind/TeaGo/assert"
)

func TestStat_UniqueId(t *testing.T) {
	op := &IncrementOperation{
		filter: map[string]interface{}{
			"name": "lu",
			"age":  20,
			"page": 1024,
			"bool": true,
		},
		field: "count",
	}

	t.Log("unique id:", op.uniqueId())
}

func TestStat_Increase(t *testing.T) {
	stat := new(Stat)

	a := assert.NewAssertion(t)

	for i := 0; i < 1000; i ++ {
		stat.Increase(findCollection("stats.test", nil), map[string]interface{}{
			"serverId": "123",
		}, map[string]interface{}{
			"serverId": "123",
		}, "count")
	}
	stat.Increase(findCollection("stats.test", nil), map[string]interface{}{
		"serverId": "1234",
	}, map[string]interface{}{
		"serverId": "1234",
	}, "count")

	a.Cost()

	time.Sleep(3 * time.Second)
	t.Log("OK")
}
