package actions

import "net/http"

type AllowAction struct {
}

func (this *AllowAction) Perform(request *http.Request, writer http.ResponseWriter) (allow bool) {
	// do nothing
	return true
}
