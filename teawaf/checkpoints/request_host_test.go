package checkpoints

import (
	"net/http"
	"testing"
)

func TestRequestHostCheckPoint_RequestValue(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://teaos.cn/?name=lu", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Host", "cloud.teaos.cn")

	checkPoint := new(RequestHostCheckPoint)
	t.Log(checkPoint.RequestValue(req, ""))
}
