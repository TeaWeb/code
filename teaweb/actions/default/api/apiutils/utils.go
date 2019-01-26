package apiutils

import (
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
)

func ValidateUser(actionPtr actions.ActionWrapper) {
	action := actionPtr.Object()
	key, found := action.Param("TeaKey")
	if !found || len(key) == 0 {
		action.Fail("Authenticate Failed 001")
	}

	user := configs.SharedAdminConfig().FindUserWithKey(key)
	if user == nil {
		action.Fail("Authenticate Failed 002")
	}
}
