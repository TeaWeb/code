package helpers

import (
	"github.com/iwind/TeaGo/actions"
	"net/http"
	"time"
)

type UserShouldAuth struct {
	action *actions.ActionObject
}

func (auth *UserShouldAuth) BeforeAction(actionPtr actions.ActionWrapper, paramName string) (goNext bool) {
	auth.action = actionPtr.Object()
	return true
}

func (auth *UserShouldAuth) StoreUsername(username string) {
	// 修改sid的时间
	cookie := &http.Cookie{
		Name:    "sid",
		Value:   auth.action.Session().Sid,
		Path:    "/",
		Expires: time.Now().Add(30 * 86400 * time.Second),
	}
	auth.action.AddCookie(cookie)
	auth.action.Session().Write("username", username)
}

func (auth *UserShouldAuth) Logout() {
	auth.action.Session().Delete()
}
