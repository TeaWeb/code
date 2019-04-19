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

	checkPoint := new(RequestArgCheckPoint)
	t.Log(checkPoint.RequestValue(req, "name"))
	t.Log(checkPoint.RequestValue(req, "name2"))
}
