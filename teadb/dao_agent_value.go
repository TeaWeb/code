package teadb

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
)

type AgentValueDAOInterface interface {
	Init()
	TableName(agentId string) string
	Insert(agentId string, value *agents.Value) error
	ClearItemValues(agentId string, appId string, itemId string, level notices.NoticeLevel) error
	FindLatestItemValue(agentId string, appId string, itemId string) (*agents.Value, error)
	FindLatestItemValueNoError(agentId string, appId string, itemId string) (*agents.Value, error)

	// 取得最近的数值记录
	FindLatestItemValues(agentId string, appId string, itemId string, noticeLevel notices.NoticeLevel, lastId string, size int) ([]*agents.Value, error)

	ListItemValues(agentId string, appId string, itemId string, noticeLevel notices.NoticeLevel, lastId string, offset int, size int) ([]*agents.Value, error)
	QueryValues(query *Query) ([]*agents.Value, error)
	GroupValuesByTime(query *Query, timeField string, result map[string]Expr) ([]*agents.Value, error)
	DropAgentTable(agentId string) error
}
