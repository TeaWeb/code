package checkpoints

import (
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
	"testing"
)

func TestArgParam_RequestValue(t *testing.T) {
	rawReq, err := http.NewRequest(http.MethodGet, "http://teaos.cn/?name=lu", nil)
	if err != nil {
		t.Fatal(err)
	}

	req := requests.NewRequest(rawReq)

	checkpoint := new(RequestArgCheckpoint)
	t.Log(checkpoint.RequestValue(req, "name", nil))
	t.Log(checkpoint.ResponseValue(req, nil, "name", nil))
	t.Log(checkpoint.RequestValue(req, "name2", nil))
}
