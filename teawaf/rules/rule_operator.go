package rules

type RuleOperator = string
type RuleCaseInsensitive = string

const (
	RuleOperatorGt          = "gt"
	RuleOperatorGte         = "gte"
	RuleOperatorLt          = "lt"
	RuleOperatorLte         = "lte"
	RuleOperatorEq          = "eq"
	RuleOperatorNeq         = "neq"
	RuleOperatorEqString    = "eq string"
	RuleOperatorNeqString   = "neq string"
	RuleOperatorMatch       = "match"
	RuleOperatorNotMatch    = "not match"
	RuleOperatorContains    = "contains"
	RuleOperatorNotContains = "not contains"
	RuleOperatorPrefix      = "prefix"
	RuleOperatorSuffix      = "suffix"
	RuleOperatorHasKey      = "has key" // has key in slice or map
	RuleOperatorVersionGt   = "version gt"
	RuleOperatorVersionLt   = "version lt"

	RuleCaseInsensitiveNone = "none"
	RuleCaseInsensitiveYes  = "yes"
	RuleCaseInsensitiveNo   = "no"
)

type RuleOperatorDefinition struct {
	Name            string
	Code            string
	Description     string
	CaseInsensitive RuleCaseInsensitive // default caseInsensitive setting
}

var AllRuleOperators = []*RuleOperatorDefinition{
	{
		Name:            "数值大于",
		Code:            RuleOperatorGt,
		Description:     "使用数值对比大于",
		CaseInsensitive: RuleCaseInsensitiveNone,
	},
	{
		Name:            "数值大于等于",
		Code:            RuleOperatorGte,
		Description:     "使用数值对比大于等于",
		CaseInsensitive: RuleCaseInsensitiveNone,
	},
	{
		Name:            "数值小于",
		Code:            RuleOperatorLt,
		Description:     "使用数值对比小于",
		CaseInsensitive: RuleCaseInsensitiveNone,
	},
	{
		Name:            "数值小于等于",
		Code:            RuleOperatorLte,
		Description:     "使用数值对比小于等于",
		CaseInsensitive: RuleCaseInsensitiveNone,
	},
	{
		Name:            "数值等于",
		Code:            RuleOperatorEq,
		Description:     "使用数值对比等于",
		CaseInsensitive: RuleCaseInsensitiveNone,
	},
	{
		Name:            "数值不等于",
		Code:            RuleOperatorNeq,
		Description:     "使用数值对比不等于",
		CaseInsensitive: RuleCaseInsensitiveNone,
	},
	{
		Name:            "字符串等于",
		Code:            RuleOperatorEqString,
		Description:     "使用字符串对比等于",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "字符串不等于",
		Code:            RuleOperatorNeqString,
		Description:     "使用字符串对比不等于",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "正则匹配",
		Code:            RuleOperatorMatch,
		Description:     "使用正则表达式匹配，在头部使用(?i)表示不区分大小写，<a href=\"http://teaos.cn/doc/regexp/Regexp.md\" target=\"_blank\">正则表达式语法 &raquo;</a>",
		CaseInsensitive: RuleCaseInsensitiveYes,
	},
	{
		Name:            "正则不匹配",
		Code:            RuleOperatorNotMatch,
		Description:     "使用正则表达式不匹配，在头部使用(?i)表示不区分大小写，<a href=\"http://teaos.cn/doc/regexp/Regexp.md\" target=\"_blank\">正则表达式语法 &raquo;</a>",
		CaseInsensitive: RuleCaseInsensitiveYes,
	},
	{
		Name:            "包含字符串",
		Code:            RuleOperatorContains,
		Description:     "包含某个字符串",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "不包含字符串",
		Code:            RuleOperatorNotContains,
		Description:     "不包含某个字符串",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "包含前缀",
		Code:            RuleOperatorPrefix,
		Description:     "包含某个前缀",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "包含后缀",
		Code:            RuleOperatorSuffix,
		Description:     "包含某个后缀",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "包含索引",
		Code:            RuleOperatorHasKey,
		Description:     "对于一组数据拥有某个键值或者索引",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
	{
		Name:            "版本号大于",
		Code:            RuleOperatorVersionGt,
		Description:     "对于版本号大于",
		CaseInsensitive: RuleCaseInsensitiveNo,
	},
}
