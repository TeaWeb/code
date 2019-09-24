package shared

import "github.com/iwind/TeaGo/maps"

// 运算符定义
type RequestCondOperator = string

const (
	RequestCondOperatorRegexp            RequestCondOperator = "regexp"
	RequestCondOperatorNotRegexp         RequestCondOperator = "not regexp"
	RequestCondOperatorEqInt             RequestCondOperator = "eq int"   // 整数等于
	RequestCondOperatorEqFloat           RequestCondOperator = "eq float" // 浮点数等于
	RequestCondOperatorGtFloat           RequestCondOperator = "gt"
	RequestCondOperatorGteFloat          RequestCondOperator = "gte"
	RequestCondOperatorLtFloat           RequestCondOperator = "lt"
	RequestCondOperatorLteFloat          RequestCondOperator = "lte"
	RequestCondOperatorEqString          RequestCondOperator = "eq"
	RequestCondOperatorNeqString         RequestCondOperator = "not"
	RequestCondOperatorHasPrefix         RequestCondOperator = "prefix"
	RequestCondOperatorHasSuffix         RequestCondOperator = "suffix"
	RequestCondOperatorContainsString    RequestCondOperator = "contains"
	RequestCondOperatorNotContainsString RequestCondOperator = "not contains"
	RequestCondOperatorIn                RequestCondOperator = "in"
	RequestCondOperatorNotIn             RequestCondOperator = "not in"
	RequestCondOperatorEqIP              RequestCondOperator = "eq ip"
	RequestCondOperatorGtIP              RequestCondOperator = "gt ip"
	RequestCondOperatorGteIP             RequestCondOperator = "gte ip"
	RequestCondOperatorLtIP              RequestCondOperator = "lt ip"
	RequestCondOperatorLteIP             RequestCondOperator = "lte ip"
	RequestCondOperatorIPInRange         RequestCondOperator = "ip range"
)

// 所有的运算符
func AllRequestOperators() []maps.Map {
	return []maps.Map{
		{
			"name":        "正则表达式匹配",
			"op":          RequestCondOperatorRegexp,
			"description": "判断是否正则表达式匹配",
		},
		{
			"name":        "正则表达式不匹配",
			"op":          RequestCondOperatorNotRegexp,
			"description": "判断是否正则表达式不匹配",
		},
		{
			"name":        "字符串等于",
			"op":          RequestCondOperatorEqString,
			"description": "使用字符串对比参数值是否相等于某个值",
		},
		{
			"name":        "字符串前缀",
			"op":          RequestCondOperatorHasPrefix,
			"description": "参数值包含某个前缀",
		},
		{
			"name":        "字符串后缀",
			"op":          RequestCondOperatorHasSuffix,
			"description": "参数值包含某个后缀",
		},
		{
			"name":        "字符串包含",
			"op":          RequestCondOperatorContainsString,
			"description": "参数值包含另外一个字符串",
		},
		{
			"name":        "字符串不包含",
			"op":          RequestCondOperatorNotContainsString,
			"description": "参数值不包含另外一个字符串",
		},
		{
			"name":        "字符串不等于",
			"op":          RequestCondOperatorNeqString,
			"description": "使用字符串对比参数值是否不相等于某个值",
		},
		{
			"name":        "在列表中",
			"op":          RequestCondOperatorIn,
			"description": "判断参数值在某个列表中",
		},
		{
			"name":        "不在列表中",
			"op":          RequestCondOperatorNotIn,
			"description": "判断参数值不在某个列表中",
		},
		{
			"name":        "整数等于",
			"op":          RequestCondOperatorEqInt,
			"description": "将参数转换为整数数字后进行对比",
		},
		{
			"name":        "浮点数等于",
			"op":          RequestCondOperatorEqFloat,
			"description": "将参数转换为可以有小数的浮点数字进行对比",
		},
		{
			"name":        "数字大于",
			"op":          RequestCondOperatorGtFloat,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "数字大于等于",
			"op":          RequestCondOperatorGteFloat,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "数字小于",
			"op":          RequestCondOperatorLtFloat,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "数字小于等于",
			"op":          RequestCondOperatorLteFloat,
			"description": "将参数转换为数字进行对比",
		},
		{
			"name":        "IP等于",
			"op":          RequestCondOperatorEqIP,
			"description": "将参数转换为IP进行对比",
		},
		{
			"name":        "IP大于",
			"op":          RequestCondOperatorGtIP,
			"description": "将参数转换为IP进行对比",
		},
		{
			"name":        "IP大于等于",
			"op":          RequestCondOperatorGteIP,
			"description": "将参数转换为IP进行对比",
		},
		{
			"name":        "IP小于",
			"op":          RequestCondOperatorLtIP,
			"description": "将参数转换为IP进行对比",
		},
		{
			"name":        "IP小于等于",
			"op":          RequestCondOperatorLteIP,
			"description": "将参数转换为IP进行对比",
		},
		{
			"name":        "IP范围",
			"op":          RequestCondOperatorIPInRange,
			"description": "IP在某个范围之内，范围格式可以是英文逗号分隔的ip1,ip2，或者CIDR格式的ip/bits",
		},
	}
}
