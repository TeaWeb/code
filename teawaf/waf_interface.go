package teawaf

import (
	"github.com/TeaWeb/code/teawaf/actions"
	"net/http"
)

// WAF interface
type WAFInterface interface {
	MatchRequest(req *http.Request, writer http.ResponseWriter) *actions.Action
	MatchResponse(req *http.Request, resp *http.Response, writer http.ResponseWriter) *actions.Action
}
