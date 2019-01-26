package teamongo

import "github.com/TeaWeb/code/teaconfigs/audits"

// 审计日志
func NewAuditsQuery() *Query {
	return NewQuery("logs.audit", new(audits.Log))
}
