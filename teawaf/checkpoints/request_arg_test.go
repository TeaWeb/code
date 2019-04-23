package checkpoints

import (
	"net/http"
	"testing"
)

func TestArgParam_RequestValue(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://teaos.cn/?name=lu", nil)
	if err != nil {
		t.Fatal(err)
	}

	checkpoint := new(RequestArgCheckpoint)
	t.Log(checkpoint.RequestValue(req, "name"))
	t.Log(checkpoint.RequestValue(req, "name2"))
}
