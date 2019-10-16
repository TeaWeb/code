package actions

import (
	"net/http"
)

type LogAction struct {
}

func (this *LogAction) Perform(request *http.Request, writer http.ResponseWriter) (allow bool) {
	return true
}
