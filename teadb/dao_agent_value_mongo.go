package teadb

import (
	"errors"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teadb/shared"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/logs"
	"golang.org/x/net/context"
	"time"
)

type MongoAgentValueDAO struct {
}

func (this *MongoAgentValueDAO) Init() {
}

func (this *MongoAgentValueDAO) TableName(agentId string) string {
	return this.agentCollName(agentId)
}

func (this *MongoAgentValueDAO) Insert(agentId string, value *agents.Value) error {
	if value == nil {
		return errors.New("value should not be nil")
	}
	if len(agentId) == 0 {
		if len(value.AgentId) > 0 {
			agentId = value.AgentId
		} else {
			return errors.New("AgentId should be set")
		}
	}

	if value.Value == nil {
		value.Value = 0
	}

	if value.Id.IsZero() {
		value.Id = shared.NewObjectId()
	}

	coll := this.selectColl(this.agentCollName(agentId))
	_, err := coll.InsertOne(context.Background(), *value)
	return err
}

func (this *MongoAgentValueDAO) ClearItemValues(agentId string, appId string, itemId string, level notices.NoticeLevel) error {
	if len(agentId) == 0 {
		return errors.New("agentId should not be empty")
	}
	query := NewQuery(this.agentCollName(agentId)).
		Attr("appId", appId).
		Attr("itemId", itemId)
	if level > 0 {
		query.Attr("noticeLevel", level)
	}
	return query.Delete()
}

func (this *MongoAgentValueDAO) FindLatestItemValue(agentId string, appId string, itemId string) (*agents.Value, error) {
	query := NewQuery(this.agentCollName(agentId)).
		Attr("itemId", itemId).
		Node().
		Desc("createdAt")
	if len(appId) > 0 {
		query.Attr("appId", appId)
	}
	v, err := query.FindOne(new(agents.Value))
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return this.processValue(v.(*agents.Value)), nil
}

func (this *MongoAgentValueDAO) FindLatestItemValueNoError(agentId string, appId string, itemId string) (*agents.Value, error) {
	query := NewQuery(this.agentCollName(agentId)).
		Attr("itemId", itemId).
		Attr("error", "").
		Node().
		Desc("createdAt")
	if len(appId) > 0 {
		query.Attr("appId", appId)
	}
	v, err := query.FindOne(new(agents.Value))
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return this.processValue(v.(*agents.Value)), nil
}

// 取得最近的数值记录
func (this *MongoAgentValueDAO) FindLatestItemValues(agentId string, appId string, itemId string, noticeLevel notices.NoticeLevel, lastId string, size int) ([]*agents.Value, error) {
	query := NewQuery(this.agentCollName(agentId))
	query.Attr("appId", appId)
	query.Attr("itemId", itemId)
	query.Node()
	query.Limit(size)
	query.Desc("createdAt")

	if noticeLevel > 0 {
		if noticeLevel == notices.NoticeLevelInfo {
			query.Attr("noticeLevel", []interface{}{notices.NoticeLevelInfo, notices.NoticeLevelNone})
		} else {
			query.Attr("noticeLevel", noticeLevel)
		}
	}

	if len(lastId) > 0 {
		lastObjectId, err := shared.ObjectIdFromHex(lastId)
		if err != nil {
			return nil, err
		}
		query.Gt("_id", lastObjectId)
	}

	ones, err := query.FindOnes(new(agents.Value))
	if err != nil {
		return nil, err
	}
	result := []*agents.Value{}
	for _, one := range ones {
		result = append(result, this.processValue(one.(*agents.Value)))
	}
	return result, nil
}

func (this *MongoAgentValueDAO) ListItemValues(agentId string, appId string, itemId string, noticeLevel notices.NoticeLevel, lastId string, offset int, size int) ([]*agents.Value, error) {
	query := NewQuery(this.agentCollName(agentId))
	query.Attr("appId", appId)
	query.Attr("itemId", itemId)
	query.Node()
	query.Offset(offset)
	query.Limit(size)
	query.Desc("createdAt")

	if noticeLevel > 0 {
		if noticeLevel == notices.NoticeLevelInfo {
			query.Attr("noticeLevel", []interface{}{notices.NoticeLevelInfo, notices.NoticeLevelNone})
		} else {
			query.Attr("noticeLevel", noticeLevel)
		}
	}

	if len(lastId) > 0 {
		lastObjectId, err := shared.ObjectIdFromHex(lastId)
		if err != nil {
			return nil, err
		}
		query.Lt("_id", lastObjectId)
	}

	ones, err := query.FindOnes(new(agents.Value))
	if err != nil {
		return nil, err
	}
	result := []*agents.Value{}
	for _, one := range ones {
		result = append(result, this.processValue(one.(*agents.Value)))
	}
	return result, nil
}

func (this *MongoAgentValueDAO) QueryValues(query *Query) ([]*agents.Value, error) {
	ones, err := query.FindOnes(new(agents.Value))
	if err != nil {
		return nil, err
	}

	result := []*agents.Value{}
	for _, one := range ones {
		result = append(result, this.processValue(one.(*agents.Value)))
	}
	return result, nil
}

func (this *MongoAgentValueDAO) GroupValuesByTime(query *Query, timeField string, result map[string]Expr) ([]*agents.Value, error) {
	query.Asc("timeFormat." + timeField)
	result["timeFormat"] = "timeFormat"
	ones, err := query.Group("timeFormat."+timeField, result)
	if err != nil {
		return nil, err
	}

	values := []*agents.Value{}
	for _, one := range ones {
		value := agents.NewValue()
		timeFormat := one.GetMap("timeFormat")
		one.Delete("_id", "timeFormat")
		value.Value = one
		value.TimeFormat.Year = timeFormat.GetString("year")
		value.TimeFormat.Month = timeFormat.GetString("month")
		value.TimeFormat.Day = timeFormat.GetString("day")
		value.TimeFormat.Hour = timeFormat.GetString("hour")
		value.TimeFormat.Minute = timeFormat.GetString("minute")
		value.TimeFormat.Second = timeFormat.GetString("second")
		values = append(values, value)
	}
	return values, nil
}

func (this *MongoAgentValueDAO) DropAgentTable(agentId string) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return this.selectColl(this.agentCollName(agentId)).Drop(ctx)
}

func (this *MongoAgentValueDAO) selectColl(collName string) *teamongo.Collection {
	coll := teamongo.SharedCollection(collName)

	if isInitializedTable(collName) {
		return coll
	}
	err := coll.CreateIndex(
		shared.NewIndexField("appId", true),
		shared.NewIndexField("itemId", true),
		shared.NewIndexField("createdAt", false),
	)
	if err != nil {
		logs.Error(err)
	}
	err = coll.CreateIndex(
		shared.NewIndexField("appId", true),
		shared.NewIndexField("itemId", true),
		shared.NewIndexField("nodeId", true),
		shared.NewIndexField("createdAt", false),
	)
	if err != nil {
		logs.Error(err)
	}
	return coll
}

func (this *MongoAgentValueDAO) agentCollName(agentId string) string {
	return "values.agent." + agentId
}

func (this *MongoAgentValueDAO) processValue(ptr *agents.Value) *agents.Value {
	if ptr.Value == nil {
		return ptr
	}
	v, err := teamongo.BSONDecode(ptr.Value)
	if err == nil {
		ptr.Value = v
	} else {
		logs.Error(err)
	}
	return ptr
}
