package checkpoints

import (
	"net/http"
)

type SampleCheckPoint struct {
	CheckPoint
}

func (this *SampleCheckPoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	return
}
