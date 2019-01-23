package teamongo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"reflect"
	"sync"
	"time"
)

type ValueQueryAction = string

const (
	ValueQueryActionCount   = "count"
	ValueQueryActionSum     = "sum"
	ValueQueryActionAvg     = "avg"
	ValueQueryActionMin     = "min"
	ValueQueryActionMax     = "max"
	ValueQueryActionFind    = "find"
	ValueQueryActionFindAll = "findAll"
)

var valueCollectionsMap = map[string]*Collection{}
var valueCollectionsLocker sync.Mutex

type ValueQuery struct {
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

func NewValueQuery() *ValueQuery {
	return &ValueQuery{
		cond:   map[string]interface{}{},
		sorts:  []map[string]int{},
		offset: -1,
		size:   -1,
	}
}

func (this *ValueQuery) Debug() *ValueQuery {
	this.debug = true
	return this
}

func (this *ValueQuery) Agent(agentId string) *ValueQuery {
	this.agentId = agentId
	return this
}

func (this *ValueQuery) App(appId string) *ValueQuery {
	this.appId = appId
	return this
}

func (this *ValueQuery) Item(itemId string) *ValueQuery {
	this.Attr("itemId", itemId)
	return this
}

func (this *ValueQuery) Asc(field string) *ValueQuery {
	if len(field) == 0 {
		field = "_id"
	}
	this.sorts = append(this.sorts, map[string]int{
		field: 1,
	})
	return this
}

func (this *ValueQuery) Desc(field string) *ValueQuery {
	if len(field) == 0 {
		field = "_id"
	}
	this.sorts = append(this.sorts, map[string]int{
		field: -1,
	})
	return this
}

func (this *ValueQuery) Offset(offset int64) *ValueQuery {
	this.offset = offset
	return this
}

func (this *ValueQuery) Limit(size int64) *ValueQuery {
	this.size = size
	return this
}

func (this *ValueQuery) Group(group []string) *ValueQuery {
	this.group = group
	return this
}

func (this *ValueQuery) Attr(field string, value interface{}) *ValueQuery {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		this.Op("in", field, value)
	} else {
		this.Op("eq", field, value)
	}
	return this
}

// 设置日志ID
func (this *ValueQuery) Id(idString string) *ValueQuery {
	objectId, err := objectid.FromHex(idString)
	if err != nil {
		logs.Error(err)
		return this.Attr("_id", idString)
	}
	this.Attr("_id", objectId)
	return this
}

func (this *ValueQuery) Op(op string, field string, value interface{}) {
	_, found := this.cond[field]
	if found {
		this.cond[field].(map[string]interface{})[op] = value
	} else {
		this.cond[field] = map[string]interface{}{
			op: value,
		}
	}
}

func (this *ValueQuery) Not(field string, value interface{}) *ValueQuery {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		this.Op("nin", field, value)
	} else {
		this.Op("ne", field, value)
	}
	return this
}

func (this *ValueQuery) Lt(field string, value interface{}) *ValueQuery {
	this.Op("lt", field, value)
	return this
}

func (this *ValueQuery) Lte(field string, value interface{}) *ValueQuery {
	this.Op("lte", field, value)
	return this
}

func (this *ValueQuery) Gt(field string, value interface{}) *ValueQuery {
	this.Op("gt", field, value)
	return this
}

func (this *ValueQuery) Gte(field string, value interface{}) *ValueQuery {
	this.Op("gte", field, value)
	return this
}

func (this *ValueQuery) Action(action ValueQueryAction, ) *ValueQuery {
	this.action = action
	return this
}

func (this *ValueQuery) For(field string) *ValueQuery {
	this.forField = field
	return this
}

// 开始执行
func (this *ValueQuery) Execute() (interface{}, error) {
	if len(this.agentId) == 0 {
		return nil, errors.New("AgentId should be set")
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
func (this *ValueQuery) Find() (*agents.Value, error) {
	result, err := this.Action(ValueQueryActionFind).Execute()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*agents.Value), nil
}

// 插入新数据
func (this *ValueQuery) Insert(value *agents.Value) error {
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
		value.Id = objectid.New()
	}

	collectionName := "values.agent." + this.agentId
	coll := this.selectColl(collectionName)
	_, err := coll.InsertOne(context.Background(), *value)
	return err
}

// 删除数据
func (this *ValueQuery) Delete() error {
	if len(this.agentId) == 0 {
		return errors.New("AgentId should be set")
	}

	filter := this.buildFilter()

	collectionName := "values.agent." + this.agentId
	coll := this.selectColl(collectionName)
	_, err := coll.DeleteMany(context.Background(), filter)
	return err
}

func (this *ValueQuery) queryNumber(collectionName string) (float64, error) {
	if this.action == ValueQueryActionCount {
		coll := this.selectColl(collectionName)
		filter := this.buildFilter()
		i, err := coll.Count(context.Background(), filter)
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

func (this *ValueQuery) queryGroup(collectionName string) (result map[string]map[string]interface{}, err error) {
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

func (this *ValueQuery) findAll(collectionName string) (result []*agents.Value, err error) {
	coll := this.selectColl(collectionName)
	opts := []findopt.Find{}
	if this.offset > -1 {
		opts = append(opts, findopt.Skip(this.offset))
	}
	if this.size > -1 {
		opts = append(opts, findopt.Limit(this.size))
	}
	if len(this.sorts) > 0 {
		for _, sort := range this.sorts {
			for field, order := range sort {
				opts = append(opts, findopt.Sort(map[string]int{
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

func (this *ValueQuery) buildFilter() map[string]interface{} {
	filter := map[string]interface{}{}

	// cond
	if len(this.cond) > 0 {
		for field, cond := range this.cond {
			fieldQuery := map[string]interface{}{}
			for op, value := range cond.(map[string]interface{}) {
				if lists.Contains([]string{"eq", "lt", "lte", "gt", "gte", "in", "nin", "ne"}, op) {
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

func (this *ValueQuery) jsonEncode(i interface{}) (string, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (this *ValueQuery) jsonEncodeString(i interface{}) string {
	data, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(data)
}

func (this *ValueQuery) selectColl(collectionName string) *Collection {
	valueCollectionsLocker.Lock()
	defer valueCollectionsLocker.Unlock()

	coll, found := valueCollectionsMap[collectionName]
	if found {
		return coll
	}

	coll = FindCollection(collectionName)
	coll.CreateIndex(map[string]bool{
		"itemId": true,
	})
	coll.CreateIndex(map[string]bool{
		"_id": false,
	})
	valueCollectionsMap[collectionName] = coll
	return coll
}
