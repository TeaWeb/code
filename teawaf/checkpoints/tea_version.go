package checkpoints

import (
	"github.com/TeaWeb/code/teaconst"
	"net/http"
)

type TeaVersionCheckpoint struct {
	Checkpoint
}

func (this *TeaVersionCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = teaconst.TeaVersion
	return
}

func (this *TeaVersionCheckpoint) ResponseValue(req *http.Request, resp *http.Response, param string) (value interface{}, err error) {
	value = teaconst.TeaVersion
	return
}
