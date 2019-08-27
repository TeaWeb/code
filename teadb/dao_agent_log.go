package teadb

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
)

type AgentLogDAOInterface interface {
	Init()

	// 插入一条数据
	InsertOne(agentId string, log *agents.ProcessLog) error

	// 获取最新任务的日志
	FindLatestTaskLogs(agentId string, taskId string, fromId string, size int) ([]*agents.ProcessLog, error)

	// 获取任务最后一次的执行日志
	FindLatestTaskLog(agentId string, taskId string) (*agents.ProcessLog, error)
}
