package teawaf

import (
	"github.com/TeaWeb/code/teawaf/actions"
	"github.com/TeaWeb/code/teawaf/rules"
)

// 感谢以下规则来源：
// - Janusec: https://www.janusec.com/
func Template() *WAF {
	waf := NewWAF()
	waf.Id = "template"
	waf.On = true

	// black list
	{
		group := rules.NewRuleGroup()
		group.On = false
		group.IsInbound = true
		group.Name = "白名单"
		group.Code = "whiteList"
		group.Description = "在此名单中的IP地址可以直接跳过防火墙设置"

		{

			set := rules.NewRuleSet()
			set.On = true
			set.Name = "IP白名单"
			set.Code = "9001"
			set.Connector = rules.RuleConnectorOr
			set.Action = actions.ActionAllow
			set.AddRule(&rules.Rule{
				Param:             "${remoteAddr}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `127\.0\.0\.1|0\.0\.0\.0`,
				IsCaseInsensitive: false,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// black list
	{
		group := rules.NewRuleGroup()
		group.On = false
		group.IsInbound = true
		group.Name = "黑名单"
		group.Code = "blackList"
		group.Description = "在此名单中的IP地址直接阻止"

		{

			set := rules.NewRuleSet()
			set.On = true
			set.Name = "IP黑名单"
			set.Code = "10001"
			set.Connector = rules.RuleConnectorOr
			set.Action = actions.ActionBlock
			set.AddRule(&rules.Rule{
				Param:             "${remoteAddr}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `1\.1\.1\.1|2\.2\.2\.2`,
				IsCaseInsensitive: false,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// xss
	{
		group := rules.NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "XSS"
		group.Code = "xss"
		group.Description = "防跨站脚本攻击（Cross Site Scripting）"

		{
			set := rules.NewRuleSet()
			set.On = true
			set.Name = "Javascript事件"
			set.Code = "1001"
			set.Connector = rules.RuleConnectorOr
			set.Action = actions.ActionBlock
			set.AddRule(&rules.Rule{
				Param:             "${requestURI}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `(onmouseover|onmousemove|onmousedown|onmouseup|onerror|onload|onclick|ondblclick|onkeydown|onkeyup|onkeypress)\s*=`, // TODO more keywords here
				IsCaseInsensitive: true,
			})
			group.AddRuleSet(set)
		}

		{
			set := rules.NewRuleSet()
			set.On = true
			set.Name = "Javascript函数"
			set.Code = "1002"
			set.Connector = rules.RuleConnectorOr
			set.Action = actions.ActionBlock
			set.AddRule(&rules.Rule{
				Param:             "${requestURI}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `(alert|eval|prompt|confirm)\s*\(`, // TODO more keywords here
				IsCaseInsensitive: true,
			})
			group.AddRuleSet(set)
		}

		{
			set := rules.NewRuleSet()
			set.On = true
			set.Name = "HTML标签"
			set.Code = "1003"
			set.Connector = rules.RuleConnectorOr
			set.Action = actions.ActionBlock
			set.AddRule(&rules.Rule{
				Param:             "${requestURI}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `<(script|iframe|link)`, // TODO more keywords here
				IsCaseInsensitive: true,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// upload
	{
		group := rules.NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "文件上传"
		group.Code = "upload"
		group.Description = "防止上传可执行脚本文件到服务器"

		{
			set := rules.NewRuleSet()
			set.On = true
			set.Name = "上传文件扩展名"
			set.Code = "2001"
			set.Connector = rules.RuleConnectorOr
			set.Action = actions.ActionBlock
			set.AddRule(&rules.Rule{
				Param:             "${requestUpload.ext}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `\.(php|jsp|aspx|asp|exe|asa|rb|py)\b`, // TODO more keywords here
				IsCaseInsensitive: true,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// web shell
	{
		group := rules.NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "Web Shell"
		group.Code = "webShell"
		group.Description = "防止远程执行服务器命令"

		{
			set := rules.NewRuleSet()
			set.On = true
			set.Name = "Web Shell"
			set.Code = "3001"
			set.Connector = rules.RuleConnectorOr
			set.Action = actions.ActionBlock
			set.AddRule(&rules.Rule{
				Param:             "${requestAll}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `\b(eval|system|exec|execute|passthru|shell_exec|phpinfo)\s*\(`, // TODO more keywords here
				IsCaseInsensitive: true,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// command injection
	{
		group := rules.NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "命令注入"
		group.Code = "commandInjection"

		{
			set := rules.NewRuleSet()
			set.On = true
			set.Name = "命令注入"
			set.Code = "4001"
			set.Connector = rules.RuleConnectorOr
			set.Action = actions.ActionBlock
			set.AddRule(&rules.Rule{
				Param:             "${requestAll}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `\b(pwd|ls|ll|whoami|id|net\s+user)\b$`, // TODO more keywords here
				IsCaseInsensitive: false,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// path traversal
	{
		group := rules.NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "路径穿越"
		group.Code = "pathTraversal"
		group.Description = "防止读取网站目录之外的其他系统文件"

		{
			set := rules.NewRuleSet()
			set.On = true
			set.Name = "路径穿越"
			set.Code = "5001"
			set.Connector = rules.RuleConnectorOr
			set.Action = actions.ActionBlock
			set.AddRule(&rules.Rule{
				Param:             "${requestURI}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `((\.+)(/+)){2,}`, // TODO more keywords here
				IsCaseInsensitive: false,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// special dirs
	{
		group := rules.NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "特殊目录"
		group.Code = "denyDirs"
		group.Description = "防止通过Web访问到一些特殊目录"

		{
			set := rules.NewRuleSet()
			set.On = true
			set.Name = "特殊目录"
			set.Code = "6001"
			set.Connector = rules.RuleConnectorOr
			set.Action = actions.ActionBlock
			set.AddRule(&rules.Rule{
				Param:             "${requestPath}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `/\.(git|svn|htaccess|idea)\b`, // TODO more keywords here
				IsCaseInsensitive: true,
			})
			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// sql injection
	{
		group := rules.NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "SQL注入"
		group.Code = "sqlInjection"
		group.Description = "防止SQL注入漏洞"

		{
			set := rules.NewRuleSet()
			set.On = true
			set.Name = "Union SQL Injection"
			set.Code = "7001"
			set.Connector = rules.RuleConnectorOr
			set.Action = actions.ActionBlock

			set.AddRule(&rules.Rule{
				Param:             "${requestAll}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `union[\s/\*]+select`,
				IsCaseInsensitive: true,
			})

			group.AddRuleSet(set)
		}

		{
			set := rules.NewRuleSet()
			set.On = true
			set.Name = "SQL注释"
			set.Code = "7002"
			set.Connector = rules.RuleConnectorOr
			set.Action = actions.ActionBlock

			set.AddRule(&rules.Rule{
				Param:             "${requestAll}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `/\*(!|\x00)`,
				IsCaseInsensitive: true,
			})

			group.AddRuleSet(set)
		}

		{
			set := rules.NewRuleSet()
			set.On = true
			set.Name = "SQL条件"
			set.Code = "7003"
			set.Connector = rules.RuleConnectorOr
			set.Action = actions.ActionBlock

			set.AddRule(&rules.Rule{
				Param:             "${requestAll}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `\s(and|or|rlike)\s+(if|updatexml)\s*\(`,
				IsCaseInsensitive: true,
			})
			set.AddRule(&rules.Rule{
				Param:             "${requestAll}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `\s+(and|or|rlike)\s+(select|case)\s+`,
				IsCaseInsensitive: true,
			})
			set.AddRule(&rules.Rule{
				Param:             "${requestAll}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `\s+(and|or|procedure)\s+[\w\p{L}]+\s*=\s*[\w\p{L}]+(\s|$|--|#)`,
				IsCaseInsensitive: true,
			})
			set.AddRule(&rules.Rule{
				Param:             "${requestAll}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `\(\s*case\s+when\s+[\w\p{L}]+\s*=\s*[\w\p{L}]+\s+then\s+`,
				IsCaseInsensitive: true,
			})

			group.AddRuleSet(set)
		}

		{
			set := rules.NewRuleSet()
			set.On = true
			set.Name = "SQL函数"
			set.Code = "7004"
			set.Connector = rules.RuleConnectorOr
			set.Action = actions.ActionBlock

			set.AddRule(&rules.Rule{
				Param:             "${requestAll}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `(updatexml|extractvalue|ascii|ord|char|chr|count|concat|rand|floor|substr|length|len|user|database|benchmark|analyse)\s*\(`,
				IsCaseInsensitive: true,
			})

			group.AddRuleSet(set)
		}

		{
			set := rules.NewRuleSet()
			set.On = true
			set.Name = "SQL附加语句"
			set.Code = "7005"
			set.Connector = rules.RuleConnectorOr
			set.Action = actions.ActionBlock

			set.AddRule(&rules.Rule{
				Param:             "${requestAll}",
				Operator:          rules.RuleOperatorMatch,
				Value:             `;\s*(declare|use|drop|create|exec|delete|update|insert)\s`,
				IsCaseInsensitive: true,
			})

			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// cc
	{
		group := rules.NewRuleGroup()
		group.On = false
		group.IsInbound = true
		group.Name = "CC攻击"
		group.Description = "Challenge Collapsar，防止短时间大量请求涌入，请谨慎开启和设置"
		group.Code = "cc"

		{
			set := rules.NewRuleSet()
			set.On = true
			set.Name = "CC请求数"
			set.Code = "8001"
			set.Connector = rules.RuleConnectorAnd
			set.Action = actions.ActionBlock
			set.AddRule(&rules.Rule{
				Param:    "${cc.requests}",
				Operator: rules.RuleOperatorGt,
				Value:    "1000",
				CheckpointOptions: map[string]string{
					"period": "60",
				},
				IsCaseInsensitive: false,
			})
			set.AddRule(&rules.Rule{
				Param:             "${remoteAddr}",
				Operator:          rules.RuleOperatorNotMatch,
				Value:             `127\.0\.0\.1|192\.168\.1\.100`,
				IsCaseInsensitive: false,
			})

			group.AddRuleSet(set)
		}

		waf.AddRuleGroup(group)
	}

	// custom
	{
		group := rules.NewRuleGroup()
		group.On = true
		group.IsInbound = true
		group.Name = "自定义规则分组"
		group.Description = "我的自定义规则分组，可以将自定义的规则放在这个分组下"
		group.Code = "custom"
		waf.AddRuleGroup(group)
	}

	return waf
}
