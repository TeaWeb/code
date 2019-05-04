package teaproxy

import (
	"github.com/TeaWeb/code/teaconfigs"
	"net/http"
	"strings"
)

// 调用Rewrite
func (this *Request) callRewrite(writer *ResponseWriter) error {
	query := this.requestQueryString()
	target := this.rewriteReplace
	if len(query) > 0 {
		if strings.Index(target, "?") > 0 {
			target += "&" + query
		} else {
			target += "?" + query
		}
	}

	if this.rewriteRedirectMode == teaconfigs.RewriteFlagRedirect {
		// 跳转
		http.Redirect(writer, this.raw, target, http.StatusTemporaryRedirect)
		return nil
	}

	if this.rewriteRedirectMode == teaconfigs.RewriteFlagProxy {
		return this.callURL(writer, this.raw.Method, target)
	}

	return nil
}
