package actions

import "net/http"

type ActionString = string

const (
	ActionLog     = "log"     // allow and log
	ActionBlock   = "block"   // block
	ActionCaptcha = "captcha" // block and show captcha // TODO
	ActionAllow   = "allow"   // allow
)

type ActionInterface interface {
	Perform(writer http.ResponseWriter) (allow bool)
}
