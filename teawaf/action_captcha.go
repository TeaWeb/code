package teawaf

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/TeaWeb/code/teawaf/requests"
	"github.com/dchest/captcha"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"net/http"
	"time"
)

var captchaSalt = stringutil.Rand(32)

const (
	CaptchaSeconds = 600 // 10 minutes
)

type CaptchaAction struct {
}

func (this *CaptchaAction) Perform(waf *WAF, request *requests.Request, writer http.ResponseWriter) (allow bool) {
	// TEAWEB_CAPTCHA:
	cookie, err := request.Cookie("TEAWEB_WAF_CAPTCHA")
	if err == nil && cookie != nil && len(cookie.Value) > 32 {
		m := cookie.Value[:32]
		timestamp := cookie.Value[32:]
		if stringutil.Md5(captchaSalt+timestamp) == m && time.Now().Unix() < types.Int64(timestamp) { // verify md5
			return true
		}
	}

	// verify
	if request.Method == http.MethodPost {
		captchaId := request.FormValue("TEAWEB_WAF_CAPTCHA_ID")
		if len(captchaId) > 0 {
			captchaCode := request.FormValue("TEAWEB_WAF_CAPTCHA_CODE")
			if captcha.VerifyString(captchaId, captchaCode) {
				// set cookie
				timestamp := fmt.Sprintf("%d", time.Now().Unix()+CaptchaSeconds)
				m := stringutil.Md5(captchaSalt + timestamp)
				http.SetCookie(writer, &http.Cookie{
					Name:   "TEAWEB_WAF_CAPTCHA",
					Value:  m + timestamp,
					MaxAge: CaptchaSeconds,
					Path:   "/", // all of dirs
				})

				http.Redirect(writer, request.Raw(), request.URL.String(), http.StatusTemporaryRedirect)

				return false
			}
		}
	}

	// show captcha
	captchaId := captcha.NewLen(6)
	buf := bytes.NewBuffer([]byte{})
	err = captcha.WriteImage(buf, captchaId, 200, 100)
	if err != nil {
		logs.Error(err)
		return true
	}

	_, _ = writer.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<title>Verify Yourself</title>
</head>
<body>
<form method="POST">
	<input type="hidden" name="TEAWEB_WAF_CAPTCHA_ID" value="` + captchaId + `"/>
	<img src="data:image/png;base64, ` + base64.StdEncoding.EncodeToString(buf.Bytes()) + `"/>
	<div>
		<p>Input verify code above:</p>
		<input type="text" name="TEAWEB_WAF_CAPTCHA_CODE" maxlength="6" size="18" autocomplete="off" z-index="1" style="font-size:16px;line-height:24px; letter-spacing: 15px; padding-left: 4px"/>
	</div>
	<div>
		<button type="submit" style="line-height:24px;margin-top:10px">Verify Yourself</button>
	</div>
</form>
</body>
</html>`))
	return false
}
