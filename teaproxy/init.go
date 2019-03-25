package teaproxy

import (
	"net/http"
)

// 状态码筛选
var StatusCodeParser func(statusCode int, headers http.Header, respData []byte, parserScript string) (string, error) = nil
