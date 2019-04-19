package teawaf

type RuleOperator = string

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
)
