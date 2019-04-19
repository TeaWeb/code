package checkpoints

import (
	"net/http"
	"testing"
)

func TestRequestPathCheckPoint_RequestValue(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://teaos.cn/index?name=lu", nil)
	if err != nil {
		t.Fatal(err)
	}

	checkPoint := new(RequestPathCheckPoint)
	t.Log(checkPoint.RequestValue(req, ""))
}
