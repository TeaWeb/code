package actions

import (
	"net/http"
)

type LogAction struct {
}

func (this *LogAction) Perform(writer http.ResponseWriter) (allow bool) {
	return true
}
