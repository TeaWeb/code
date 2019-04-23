package checkpoints

import (
	"net/http"
)

type SampleCheckpoint struct {
	Checkpoint
}

func (this *SampleCheckpoint) RequestValue(req *http.Request, param string) (value interface{}, err error) {
	return
}
