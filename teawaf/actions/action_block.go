package actions

import "net/http"

type BlockAction struct {
}

func (this *BlockAction) Perform(writer http.ResponseWriter) (allow bool) {
	if writer != nil {
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("The request is blocked by TeaWAF"))
	}
	return false
}
