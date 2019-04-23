package checkpoints

import (
	"net/http"
	"testing"
)

func TestRequestPathCheckpoint_RequestValue(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://teaos.cn/index?name=lu", nil)
	if err != nil {
		t.Fatal(err)
	}

	checkpoint := new(RequestPathCheckpoint)
	t.Log(checkpoint.RequestValue(req, ""))
}
