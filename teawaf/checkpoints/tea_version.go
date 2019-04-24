package checkpoints

import (
	"github.com/TeaWeb/code/teaconst"
	"github.com/TeaWeb/code/teawaf/requests"
	"net/http"
)

type TeaVersionCheckpoint struct {
	Checkpoint
}

func (this *TeaVersionCheckpoint) RequestValue(requests *requests.Request, param string) (value interface{}, sysErr error, userErr error) {
	value = teaconst.TeaVersion
	return
}

func (this *TeaVersionCheckpoint) ResponseValue(requests *requests.Request, resp *http.Response, param string) (value interface{}, sysErr error, userErr error) {
	value = teaconst.TeaVersion
	return
}
