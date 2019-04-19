package actions

import "net/http"

type AllowAction struct {
}

func (this *AllowAction) Perform(writer http.ResponseWriter) (allow bool) {
	// do nothing
	return true
}
