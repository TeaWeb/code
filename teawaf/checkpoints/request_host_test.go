package checkpoints

import (
	"net/http"
	"testing"
)

func TestRequestHostCheckpoint_RequestValue(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://teaos.cn/?name=lu", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Host", "cloud.teaos.cn")

	checkpoint := new(RequestHostCheckpoint)
	t.Log(checkpoint.RequestValue(req, ""))
}
