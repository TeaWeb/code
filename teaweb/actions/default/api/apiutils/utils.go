package apiutils

import (
	"github.com/TeaWeb/code/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/pquerna/ffjson/ffjson"
)

// 校验用户
func ValidateUser(actionPtr actions.ActionWrapper) {
	action := actionPtr.Object()
	action.AddHeader("Content-Type", "application/json; charset=utf-8")

	key, found := action.Param("TeaKey")
	if !found || len(key) == 0 {
		Fail(actionPtr, "Authenticate Failed 001")
	}

	user := configs.SharedAdminConfig().FindUserWithKey(key)
	if user == nil {
		Fail(actionPtr, "Authenticate Failed 002")
	}
}

// 错误提示
func Fail(actionPtr actions.ActionWrapper, message string) {
	action := actionPtr.Object()
	action.ResponseWriter.WriteHeader(400)
	action.Fail(message)
}

// 成功并返回数据
func Success(actionPtr actions.ActionWrapper, data interface{}) {
	dataBytes, err := ffjson.Marshal(data)
	if err != nil {
		Fail(actionPtr, err.Error())
		return
	}
	actionPtr.Object().Write(dataBytes)
}

// 成功
func SuccessOK(actionPtr actions.ActionWrapper) {
	actionPtr.Object().WriteJSON(maps.Map{
		"ok": 1,
	})
}

//从panic中恢复信息
func Recover(actionPtr actions.ActionWrapper, containsData bool) {
	i := recover()
	if i != nil {
		action, ok := i.(*actions.ActionObject)
		if ok {
			w, ok := action.ResponseWriter.(*actions.TestingResponseWriter)
			if !ok {
				return
			}
			m := maps.Map{}
			err := ffjson.Unmarshal(w.Data, &m)
			if err != nil {
				logs.Error(err)
				return
			}
			if m.GetInt("code") != 200 {
				Fail(actionPtr, m.GetString("message"))
			} else {
				if containsData {
					Success(actionPtr, m.Get("data"))
				} else {
					SuccessOK(actionPtr)
				}
			}
		}
	}
}
