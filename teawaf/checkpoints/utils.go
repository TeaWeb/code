package checkpoints

import "github.com/TeaWeb/code/teaconst"

// all check points list
var AllCheckPoints = []*CheckPointDefinition{
	{
		Name:        "客户端地址（IP）",
		Prefix:      "remoteAddr",
		Description: "试图通过分析X-Real-IP等Header获取的客户端地址，比如192.168.1.100",
		HasParams:   false,
		Instance:    new(RequestRemoteAddrCheckPoint),
	},
	{
		Name:        "客户端源地址（IP）",
		Prefix:      "rawRemoteAddr",
		Description: "直接连接的客户端地址，比如192.168.1.100",
		HasParams:   false,
		Instance:    new(RequestRawRemoteAddrCheckPoint),
	},
	{
		Name:        "客户端端口",
		Prefix:      "remotePort",
		Description: "直接连接的客户端地址端口",
		HasParams:   false,
		Instance:    new(RequestRemotePortCheckPoint),
	},
	{
		Name:        "客户端用户名",
		Prefix:      "remoteUser",
		Description: "通过BasicAuth登录的客户端用户名",
		HasParams:   false,
		Instance:    new(RequestRemoteUserCheckPoint),
	},
	{
		Name:        "请求URI",
		Prefix:      "requestURI",
		Description: "包含URL参数的请求URI，比如/hello/world?lang=go",
		HasParams:   false,
		Instance:    new(RequestURICheckPoint),
	},
	{
		Name:        "请求路径",
		Prefix:      "requestPath",
		Description: "不包含URL参数的请求路径，比如/hello/world",
		HasParams:   false,
		Instance:    new(RequestPathCheckPoint),
	},
	{
		Name:        "请求内容长度",
		Prefix:      "requestLength",
		Description: "请求Header中的Content-Length",
		HasParams:   false,
		Instance:    new(RequestLengthCheckPoint),
	},
	{
		Name:        "请求方法",
		Prefix:      "requestMethod",
		Description: "比如GET、POST",
		HasParams:   false,
		Instance:    new(RequestMethodCheckPoint),
	},
	{
		Name:        "请求协议",
		Prefix:      "scheme",
		Description: "比如http或https",
		HasParams:   false,
		Instance:    new(RequestSchemeCheckPoint),
	},
	{
		Name:        "HTTP协议版本",
		Prefix:      "proto",
		Description: "比如HTTP/1.1",
		HasParams:   false,
		Instance:    new(RequestProtoCheckPoint),
	},
	{
		Name:        "主机名",
		Prefix:      "host",
		Description: "比如teaos.cn",
		HasParams:   false,
		Instance:    new(RequestHostCheckPoint),
	},
	{
		Name:        "请求来源URL",
		Prefix:      "referer",
		Description: "请求Header中的Referer值",
		HasParams:   false,
		Instance:    new(RequestRefererCheckPoint),
	},
	{
		Name:        "客户端信息",
		Prefix:      "userAgent",
		Description: "比如Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103",
		HasParams:   false,
		Instance:    new(RequestUserAgentCheckPoint),
	},
	{
		Name:        "内容类型",
		Prefix:      "contentType",
		Description: "请求Header的Content-Type",
		HasParams:   false,
		Instance:    new(RequestContentTypeCheckPoint),
	},
	{
		Name:        "所有cookie组合字符串",
		Prefix:      "cookies",
		Description: "比如sid=IxZVPFhE&city=beijing&uid=18237",
		HasParams:   false,
		Instance:    new(RequestCookiesCheckPoint),
	},
	{
		Name:        "单个cookie值",
		Prefix:      "cookie",
		Description: "单个cookie值",
		HasParams:   true,
		Instance:    new(RequestCookieCheckPoint),
	},
	{
		Name:        "所有URL参数组合",
		Prefix:      "args",
		Description: "比如name=lu&age=20",
		HasParams:   false,
		Instance:    new(RequestArgsCheckPoint),
	},
	{
		Name:        "单个URL参数值",
		Prefix:      "arg",
		Description: "单个URL参数值",
		HasParams:   true,
		Instance:    new(RequestArgCheckPoint),
	},
	{
		Name:        "所有Header信息",
		Prefix:      "headers",
		Description: "使用\n隔开的Header信息字符串",
		HasParams:   false,
		Instance:    new(RequestHeadersCheckPoint),
	},
	{
		Name:        "单个Header值",
		Prefix:      "header",
		Description: "单个Header值",
		HasParams:   true,
		Instance:    new(RequestHeaderCheckPoint),
	},
	{
		Name:        "TeaWeb版本",
		Prefix:      "teaVersion",
		Description: "比如" + teaconst.TeaVersion,
		HasParams:   false,
		Instance:    new(TeaVersionCheckPoint),
	},
}

// find a check point
func FindCheckPoint(prefix string) CheckPointInterface {
	for _, def := range AllCheckPoints {
		if def.Prefix == prefix {
			return def.Instance
		}
	}
	return nil
}
