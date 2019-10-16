package actions

import "net/http"

type ActionString = string

const (
	ActionLog     = "log"     // allow and log
	ActionBlock   = "block"   // block
	ActionCaptcha = "captcha" // block and show captcha
	ActionAllow   = "allow"   // allow
)

type ActionInterface interface {
	Perform(request *http.Request, writer http.ResponseWriter) (allow bool)
}
