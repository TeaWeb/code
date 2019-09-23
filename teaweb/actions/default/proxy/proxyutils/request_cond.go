package proxyutils

import (
	"github.com/TeaWeb/code/teaconfigs"
	"net/http"
	"strings"
)

// 从请求中分析请求匹配条件
func ParseRequestConds(req *http.Request, prefix string) (conds []*teaconfigs.RequestCond, breakCond *teaconfigs.RequestCond, err error) {
	conds = []*teaconfigs.RequestCond{}

	params, ok := req.Form[prefix+"_condParams"]
	if !ok || len(params) == 0 {
		return
	}

	operators, _ := req.Form[prefix+"_condOperators"]
	values, _ := req.Form[prefix+"_condValues"]
	for index, param := range params {
		cond := teaconfigs.NewRequestCond()
		cond.Param = strings.TrimSpace(param)

		if index < len(operators) {
			cond.Operator = strings.TrimSpace(operators[index])
		} else {
			break
		}

		if index < len(values) {
			cond.Value = strings.TrimSpace(values[index])
		}

		err = cond.Validate()
		if err != nil {
			breakCond = cond
			return
		}

		conds = append(conds, cond)
	}
	return
}
