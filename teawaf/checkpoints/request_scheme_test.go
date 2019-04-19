package checkpoints

import (
	"net/http"
	"testing"
)

func TestRequestSchemeCheckPoint_RequestValue(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://teaos.cn/?name=lu", nil)
	if err != nil {
		t.Fatal(err)
	}

	checkPoint := new(RequestSchemeCheckPoint)
	t.Log(checkPoint.RequestValue(req, ""))
}
