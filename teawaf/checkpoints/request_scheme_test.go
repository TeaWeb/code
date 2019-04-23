package checkpoints

import (
	"net/http"
	"testing"
)

func TestRequestSchemeCheckpoint_RequestValue(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://teaos.cn/?name=lu", nil)
	if err != nil {
		t.Fatal(err)
	}

	checkpoint := new(RequestSchemeCheckpoint)
	t.Log(checkpoint.RequestValue(req, ""))
}
