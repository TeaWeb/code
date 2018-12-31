package teaconfigs

import "github.com/iwind/TeaGo/maps"

// 运算符定义
type RewriteOperator = string

const (
	RewriteOperatorRegexp = "regexp"
	RewriteOperatorGt     = "gt"
	RewriteOperatorGte    = "gte"
	RewriteOperatorLt     = "lt"
	RewriteOperatorLte    = "lte"
	RewriteOperatorEq     = "eq"
	RewriteOperatorNot    = "not"
)

// 所有的运算符
func AllRewriteOperators() []maps.Map {
	return []maps.Map{
		{
			"name":        "正则表达式匹配",
			"op":          RewriteOperatorRegexp,
			"description": "使用正则表达式匹配",
		},
		{
			"name":        "大于",
			"op":          RewriteOperatorGt,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "大于等于",
			"op":          RewriteOperatorGte,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "小于",
			"op":          RewriteOperatorLt,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "小于等于",
			"op":          RewriteOperatorLte,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "等于",
			"op":          RewriteOperatorEq,
			"description": "使用字符串对比参数值是否相等于某个值",
		},
		{
			"name":        "不等于",
			"op":          RewriteOperatorNot,
			"description": "使用字符串对比参数值是否不相等于某个值",
		},
	}
}

// 查找某个运算符信息
func FindRewriteOperator(op string) maps.Map {
	for _, o := range AllRewriteOperators() {
		if o["op"] == op {
			return o
		}
	}
	return nil
}
