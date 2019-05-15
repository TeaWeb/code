package teamongo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"sync"
	"time"
)

type AgentValueQueryAction = string

const (
	ValueQueryActionCount   = "count"
	ValueQueryActionSum     = "sum"
	ValueQueryActionAvg     = "avg"
	ValueQueryActionMin     = "min"
	ValueQueryActionMax     = "max"
	ValueQueryActionFind    = "find"
	ValueQueryActionFindAll = "findAll"
)

var agentValueCollectionsMap = map[string]*Collection{}
var agentValueCollectionsLocker sync.Mutex

type AgentValueQuery struct {
	action   string
	agentId  string
	appId    string
	group    []string
	cond     map[string]interface{}
	forField string

	sorts  []map[string]int
	offset int64
	size   int64

	debug bool
}

func NewAgentValueQuery() *AgentValueQuery {
	return &AgentValueQuery{
		cond:   map[string]interface{}{},
		sorts:  []map[string]int{},
		offset: -1,
		size:   -1,
	}
}

func (this *AgentValueQuery) Debug() *AgentValueQuery {
	this.debug = true
	return this
}

func (this *AgentValueQuery) Agent(agentId string) *AgentValueQuery {
	this.agentId = agentId
	return this
}

func (this *AgentValueQuery) App(appId string) *AgentValueQuery {
	this.appId = appId
	return this
}

func (this *AgentValueQuery) Item(itemId string) *AgentValueQuery {
	this.Attr("itemId", itemId)
	return this
}

func (this *AgentValueQuery) Asc(field string) *AgentValueQuery {
	if len(field) == 0 {
		field = "_id"
	}
	this.sorts = append(this.sorts, map[string]int{
		field: 1,
	})
	return this
}

func (this *AgentValueQuery) Desc(field string) *AgentValueQuery {
	if len(field) == 0 {
		field = "_id"
	}
	this.sorts = append(this.sorts, map[string]int{
		field: -1,
	})
	return this
}

func (this *AgentValueQuery) Offset(offset int64) *AgentValueQuery {
	this.offset = offset
	return this
}

func (this *AgentValueQuery) Limit(size int64) *AgentValueQuery {
	this.size = size
	return this
}

func (this *AgentValueQuery) Group(group []string) *AgentValueQuery {
	this.group = group
	return this
}

func (this *AgentValueQuery) Attr(field string, value interface{}) *AgentValueQuery {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		this.Op("in", field, value)
	} else {
		this.Op("eq", field, value)
	}
	return this
}

// 设置日志ID
func (this *AgentValueQuery) Id(idString string) *AgentValueQuery {
	objectId, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		logs.Error(err)
		return this.Attr("_id", idString)
	}
	this.Attr("_id", objectId)
	return this
}

func (this *AgentValueQuery) Op(op string, field string, value interface{}) {
	_, found := this.cond[field]
	if found {
		this.cond[field].(map[string]interface{})[op] = value
	} else {
		this.cond[field] = map[string]interface{}{
			op: value,
		}
	}
}

func (this *AgentValueQuery) Not(field string, value interface{}) *AgentValueQuery {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		this.Op("nin", field, value)
	} else {
		this.Op("ne", field, value)
	}
	return this
}

func (this *AgentValueQuery) Lt(field string, value interface{}) *AgentValueQuery {
	this.Op("lt", field, value)
	return this
}

func (this *AgentValueQuery) Lte(field string, value interface{}) *AgentValueQuery {
	this.Op("lte", field, value)
	return this
}

func (this *AgentValueQuery) Gt(field string, value interface{}) *AgentValueQuery {
	this.Op("gt", field, value)
	return this
}

func (this *AgentValueQuery) Gte(field string, value interface{}) *AgentValueQuery {
	this.Op("gte", field, value)
	return this
}

func (this *AgentValueQuery) Action(action AgentValueQueryAction, ) *AgentValueQuery {
	this.action = action
	return this
}

func (this *AgentValueQuery) For(field string) *AgentValueQuery {
	this.forField = field
	return this
}

// 开始执行
func (this *AgentValueQuery) Execute() (interface{}, error) {
	if len(this.agentId) == 0 {
		return nil, errors.New("AgentId should be set")
	}

	node := teaconfigs.SharedNodeConfig()
	if node != nil {
		this.Attr("nodeId", node.Id)
	}

	collectionName := "values.agent." + this.agentId
	if this.action == ValueQueryActionFindAll {
		result := []*agents.Value{}
		ones, err := this.findAll(collectionName)
		if err != nil {
			return nil, err
		}
		result = append(result, ones ...)
		return result, nil
	} else if this.action == ValueQueryActionFind {
		result := []*agents.Value{}
		ones, err := this.findAll(collectionName)
		if err != nil {
			return nil, err
		}
		result = append(result, ones ...)
		if len(result) == 0 {
			return nil, nil
		}
		return result[0], nil
	} else if len(this.group) > 0 { // 按某个字段分组
		return this.queryGroup(collectionName)
	} else { // count
		return this.queryNumber(collectionName)
	}
}

// 查找单个数据
func (this *AgentValueQuery) Find() (*agents.Value, error) {
	result, err := this.Limit(1).Action(ValueQueryActionFind).Execute()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*agents.Value), nil
}

// 查找多个数据
func (this *AgentValueQuery) FindAll() ([]*agents.Value, error) {
	result, err := this.Action(ValueQueryActionFindAll).Execute()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.([]*agents.Value), nil
}

// 插入新数据
func (this *AgentValueQuery) Insert(value *agents.Value) error {
	if value == nil {
		return errors.New("value should not be nil")
	}
	if len(this.agentId) == 0 {
		if len(value.AgentId) > 0 {
			this.agentId = value.AgentId
		} else {
			return errors.New("AgentId should be set")
		}
	}

	if value.Value == nil {
		value.Value = 0
	}

	if value.Id.IsZero() {
		value.Id = primitive.NewObjectID()
	}

	collectionName := "values.agent." + this.agentId
	coll := this.selectColl(collectionName)
	_, err := coll.InsertOne(context.Background(), *value)
	return err
}

// 删除数据
func (this *AgentValueQuery) Delete() error {
	if len(this.agentId) == 0 {
		return errors.New("AgentId should be set")
	}

	filter := this.buildFilter()

	collectionName := "values.agent." + this.agentId
	coll := this.selectColl(collectionName)
	_, err := coll.DeleteMany(context.Background(), filter)
	return err
}

func (this *AgentValueQuery) queryNumber(collectionName string) (float64, error) {
	if this.action == ValueQueryActionCount {
		coll := this.selectColl(collectionName)
		filter := this.buildFilter()
		i, err := coll.CountDocuments(context.Background(), filter)
		if err != nil {
			return 0, err
		}
		return float64(i), nil
	} else {
		result, err := this.queryGroup(collectionName)
		if err != nil {
			return 0, err
		}
		if len(result) == 0 {
			return 0, nil
		}
		for _, v := range result {
			return types.Float64(v["count"]), nil
		}
	}
	return 0, nil
}

func (this *AgentValueQuery) queryGroup(collectionName string) (result map[string]map[string]interface{}, err error) {
	result = map[string]map[string]interface{}{}

	var groupId interface{} = nil
	if len(this.group) == 1 {
		groupId = "$" + this.group[0]
	} else if len(this.group) > 1 {
		groupId = lists.MapString(this.group, func(k int, v interface{}) interface{} {
			return "$" + v.(string)
		})
	}

	groupIdString, err := this.jsonEncode(groupId)
	if err != nil {
		return nil, err
	}

	var countField interface{}
	if this.action == ValueQueryActionCount {
		countField = map[string]interface{}{
			"$sum": 1,
		}
	} else if this.action == ValueQueryActionMin {
		if len(this.forField) == 0 {
			return nil, errors.New("should specify field for the action")
		}
		countField = map[string]interface{}{
			"$min": "$" + this.forField,
		}
	} else if this.action == ValueQueryActionMax {
		if len(this.forField) == 0 {
			return nil, errors.New("should specify field for the action")
		}
		countField = map[string]interface{}{
			"$max": "$" + this.forField,
		}
	} else if this.action == ValueQueryActionAvg {
		if len(this.forField) == 0 {
			return nil, errors.New("should specify field for the action")
		}
		countField = map[string]interface{}{
			"$avg": "$" + this.forField,
		}
	} else if this.action == ValueQueryActionSum {
		if len(this.forField) == 0 {
			return nil, errors.New("should specify field for the action")
		}
		countField = map[string]interface{}{
			"$sum": "$" + this.forField,
		}
	} else {
		countField = map[string]interface{}{
			"$sum": 1,
		}
	}
	countFieldString, err := this.jsonEncode(countField)
	if err != nil {
		return nil, err
	}

	filter := this.buildFilter()
	filterString, err := this.jsonEncode(filter)
	if err != nil {
		return nil, err
	}

	pipelines, err := BSONArrayBytes([]byte(`[
	{
		"$match": ` + filterString + `
	},
	{
		"$group": {
			"_id": ` + groupIdString + `,
			"count": ` + countFieldString + `
		}
	}
]`))
	if err != nil {
		return nil, err
	}

	cursor, err := this.selectColl(collectionName).Aggregate(context.Background(), pipelines)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = cursor.Close(context.Background())
		if err != nil {
			logs.Error(err)
		}
	}()

	for cursor.Next(context.Background()) {
		m := map[string]interface{}{}
		err := cursor.Decode(&m)
		if err != nil {
			return nil, err
		}
		result[types.String(m["_id"])] = m
	}

	return
}

func (this *AgentValueQuery) findAll(collectionName string) (result []*agents.Value, err error) {
	coll := this.selectColl(collectionName)
	opts := []*options.FindOptions{}
	if this.offset > -1 {
		opts = append(opts, options.Find().SetSkip(this.offset))
	}
	if this.size > -1 {
		opts = append(opts, options.Find().SetLimit(this.size))
	}
	if len(this.sorts) > 0 {
		for _, sort := range this.sorts {
			for field, order := range sort {
				opts = append(opts, options.Find().SetSort(map[string]int{
					field: order,
				}))
			}
		}
	}
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	cursor, err := coll.Find(ctx, this.buildFilter(), opts ...)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = cursor.Close(context.Background())
		if err != nil {
			logs.Error(err)
		}
	}()

	result = []*agents.Value{}
	for cursor.Next(context.Background()) {
		m := &agents.Value{}
		err := cursor.Decode(m)

		// m.Value处理，因为m.Value是一个interface{}，在Decode的时候有可能会变成*bson.Document
		if m.Value != nil {
			m.Value, err = BSONDecode(m.Value)
			if err != nil {
				logs.Error(err)
			}
		}

		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}

func (this *AgentValueQuery) buildFilter() map[string]interface{} {
	filter := map[string]interface{}{}

	// cond
	if len(this.cond) > 0 {
		for field, cond := range this.cond {
			fieldQuery := map[string]interface{}{}
			for op, value := range cond.(map[string]interface{}) {
				if lists.ContainsString([]string{"eq", "lt", "lte", "gt", "gte", "in", "nin", "ne"}, op) {
					fieldQuery["$"+op] = value
				}
			}
			if len(fieldQuery) > 0 {
				filter[field] = fieldQuery
			}
		}
	}

	// app id
	if len(this.appId) > 0 {
		filter["appId"] = this.appId
	}

	if this.debug {
		logs.PrintAsJSON(filter)
	}

	return filter
}

func (this *AgentValueQuery) jsonEncode(i interface{}) (string, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (this *AgentValueQuery) jsonEncodeString(i interface{}) string {
	data, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(data)
}

func (this *AgentValueQuery) selectColl(collectionName string) *Collection {
	agentValueCollectionsLocker.Lock()
	defer agentValueCollectionsLocker.Unlock()

	coll, found := agentValueCollectionsMap[collectionName]
	if found {
		return coll
	}

	coll = FindCollection(collectionName)
	coll.CreateIndex(map[string]bool{
		"itemId": true,
	})
	coll.CreateIndex(map[string]bool{
		"appId":  true,
		"itemId": true,
	})
	agentValueCollectionsMap[collectionName] = coll
	return coll
}
