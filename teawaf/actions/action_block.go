package actions

import "net/http"

type BlockAction struct {
}

func (this *BlockAction) Perform(writer http.ResponseWriter) (allow bool) {
	if writer != nil {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("The quest is blocked by TeaWAF"))
	}
	return false
}
