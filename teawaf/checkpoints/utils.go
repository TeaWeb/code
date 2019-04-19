package checkpoints

// all check points list
var AllCheckPoints = []CheckPointDefinition{
	{
		Name:     "TeaWeb版本",
		Prefix:   "teaVersion",
		Instance: new(TeaVersionCheckPoint),
	},
	{
		Name:     "客户端源地址（IP）",
		Prefix:   "rawRemoteAddr",
		Instance: new(RequestRawRemoteAddrCheckPoint),
	},
	{
		Name:     "客户端地址（IP）",
		Prefix:   "remoteAddr",
		Instance: new(RequestRemoteAddrCheckPoint),
	},
	{
		Name:     "客户端端口",
		Prefix:   "remotePort",
		Instance: new(RequestRemotePortCheckPoint),
	},
	{
		Name:     "客户端用户名",
		Prefix:   "remoteUser",
		Instance: new(RequestRemoteUserCheckPoint),
	},
	{
		Name:     "请求URI",
		Prefix:   "requestURI",
		Instance: new(RequestURICheckPoint),
	},
	{
		Name:     "请求路径（不包含参数）",
		Prefix:   "requestPath",
		Instance: new(RequestPathCheckPoint),
	},
	{
		Name:     "请求内容长度",
		Prefix:   "requestLength",
		Instance: new(RequestLengthCheckPoint),
	},
	{
		Name:     "请求方法",
		Prefix:   "requestMethod",
		Instance: new(RequestMethodCheckPoint),
	},
	{
		Name:     "请求协议，http或https",
		Prefix:   "scheme",
		Instance: new(RequestSchemeCheckPoint),
	},
	{
		Name:     "包含版本的HTTP请求协议，类似于HTTP/1.0",
		Prefix:   "proto",
		Instance: new(RequestProtoCheckPoint),
	},
	{
		Name:     "主机名",
		Prefix:   "host",
		Instance: new(RequestHostCheckPoint),
	},
	{
		Name:     "请求来源URL",
		Prefix:   "referer",
		Instance: new(RequestRefererCheckPoint),
	},
	{
		Name:     "客户端信息",
		Prefix:   "userAgent",
		Instance: new(RequestUserAgentCheckPoint),
	},
	{
		Name:     "请求头部的Content-Type",
		Prefix:   "contentType",
		Instance: new(RequestContentTypeCheckPoint),
	},
	{
		Name:     "所有cookie组合字符串",
		Prefix:   "cookies",
		Instance: new(RequestCookiesCheckPoint),
	},
	{
		Name:     "单个cookie值",
		Prefix:   "cookie",
		Instance: new(RequestCookieCheckPoint),
	},
	{
		Name:     "所有请求参数组合",
		Prefix:   "args",
		Instance: new(RequestArgsCheckPoint),
	},
	{
		Name:     "单个请求参数值",
		Prefix:   "arg",
		Instance: new(RequestArgCheckPoint),
	},
	{
		Name:     "所有Header信息组合字符串",
		Prefix:   "headers",
		Instance: new(RequestHeadersCheckPoint),
	},
	{
		Name:     "单个Header值",
		Prefix:   "header",
		Instance: new(RequestHeaderCheckPoint),
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
