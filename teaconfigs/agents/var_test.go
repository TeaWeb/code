package agents

import (
	"github.com/iwind/TeaGo/maps"
	"testing"
)

func TestEvalParam(t *testing.T) {
	t.Log(EvalParam("This is is message, ${ITEM.name}", nil, nil, maps.Map{
		"ITEM": maps.Map{
			"name": "MySQL",
		},
	}))
}
