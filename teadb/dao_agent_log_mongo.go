package teadb

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teadb/shared"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/logs"
)

type MongoAgentLogDAO struct {
}

func (this *MongoAgentLogDAO) Init() {

}

func (this *MongoAgentLogDAO) TableName(agentId string) string {
	return "logs.agent." + agentId
}

// 插入一条数据
func (this *MongoAgentLogDAO) InsertOne(agentId string, log *agents.ProcessLog) error {
	this.initTable(agentId)
	if log.Id.IsZero() {
		log.Id = shared.NewObjectId()
	}
	return NewQuery(this.TableName(agentId)).
		InsertOne(log)
}

// 获取任务的日志
func (this *MongoAgentLogDAO) FindLatestTaskLogs(agentId string, taskId string, fromId string, size int) ([]*agents.ProcessLog, error) {
	result := []*agents.ProcessLog{}

	query := NewQuery(this.TableName(agentId))
	query.Attr("taskId", taskId).
		Desc("_id").
		Limit(size)

	if len(fromId) > 0 {
		lastObjectId, err := shared.ObjectIdFromHex(fromId)
		if err != nil {
			return result, err
		}
		query.Gt("_id", lastObjectId)
	}

	ones, err := query.FindOnes(new(agents.ProcessLog))
	if err != nil {
		return result, err
	}

	for _, one := range ones {
		result = append(result, one.(*agents.ProcessLog))
	}
	return result, nil
}

// 获取任务最后一次的执行日志
func (this *MongoAgentLogDAO) FindLatestTaskLog(agentId string, taskId string) (*agents.ProcessLog, error) {
	one, err := NewQuery(this.TableName(agentId)).
		Attr("taskId", taskId).
		Desc("_id").
		FindOne(new(agents.ProcessLog))
	if err != nil || one == nil {
		return nil, err
	}
	return one.(*agents.ProcessLog), nil
}

func (this *MongoAgentLogDAO) initTable(agentId string) {
	table := this.TableName(agentId)
	if isInitializedTable(table) {
		return
	}

	coll := teamongo.SharedCollection(table)
	err := coll.CreateIndex(shared.NewIndexField("agentId", true))
	if err != nil {
		logs.Error(err)
	}
	err = coll.CreateIndex(shared.NewIndexField("taskId", true))
	if err != nil {
		logs.Error(err)
	}
}
