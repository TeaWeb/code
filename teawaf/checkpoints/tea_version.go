package checkpoints

import (
	"github.com/TeaWeb/code/teaconst"
	"net/http"
)

type TeaVersionCheckPoint struct {
	CheckPoint
}

func (this *TeaVersionCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	value = teaconst.TeaVersion
	return
}

func (this *TeaVersionCheckPoint) ResponseValue(req *http.Request, resp *http.Response, param string) (value interface{}, err error) {
	value = teaconst.TeaVersion
	return
}
