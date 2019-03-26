package noticeutils

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"
)

type NoticeQueryAction = string

const (
	NoticeQueryActionCount   = "count"
	NoticeQueryActionSum     = "sum"
	NoticeQueryActionAvg     = "avg"
	NoticeQueryActionMin     = "min"
	NoticeQueryActionMax     = "max"
	NoticeQueryActionFind    = "find"
	NoticeQueryActionFindAll = "findAll"
)

type NoticeQuery struct {
	action    string
	agentCond *notices.AgentCond
	proxyCond *notices.ProxyCond
	group     []string
	cond      map[string]interface{}
	forField  string

	sorts  []map[string]int
	offset int64
	size   int64

	debug bool
}

func NewNoticeQuery() *NoticeQuery {
	return &NoticeQuery{
		cond:   map[string]interface{}{},
		sorts:  []map[string]int{},
		offset: -1,
		size:   -1,
	}
}

func (this *NoticeQuery) Debug() *NoticeQuery {
	this.debug = true
	return this
}

func (this *NoticeQuery) Agent(cond *notices.AgentCond) *NoticeQuery {
	this.agentCond = cond
	return this
}

func (this *NoticeQuery) Proxy(cond *notices.ProxyCond) *NoticeQuery {
	this.proxyCond = cond
	return this
}

func (this *NoticeQuery) Asc(field string) *NoticeQuery {
	if len(field) == 0 {
		field = "_id"
	}
	this.sorts = append(this.sorts, map[string]int{
		field: 1,
	})
	return this
}

func (this *NoticeQuery) Desc(field string) *NoticeQuery {
	if len(field) == 0 {
		field = "_id"
	}
	this.sorts = append(this.sorts, map[string]int{
		field: -1,
	})
	return this
}

func (this *NoticeQuery) Offset(offset int64) *NoticeQuery {
	this.offset = offset
	return this
}

func (this *NoticeQuery) Limit(size int64) *NoticeQuery {
	this.size = size
	return this
}

func (this *NoticeQuery) Group(group []string) *NoticeQuery {
	this.group = group
	return this
}

func (this *NoticeQuery) Attr(field string, value interface{}) *NoticeQuery {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		this.Op("in", field, value)
	} else {
		this.Op("eq", field, value)
	}
	return this
}

// 设置日志ID
func (this *NoticeQuery) Id(idString string) *NoticeQuery {
	objectId, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		logs.Error(err)
		return this.Attr("_id", idString)
	}
	this.Attr("_id", objectId)
	return this
}

func (this *NoticeQuery) Op(op string, field string, value interface{}) {
	_, found := this.cond[field]
	if found {
		this.cond[field].(map[string]interface{})[op] = value
	} else {
		this.cond[field] = map[string]interface{}{
			op: value,
		}
	}
}

func (this *NoticeQuery) Not(field string, value interface{}) *NoticeQuery {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		this.Op("nin", field, value)
	} else {
		this.Op("ne", field, value)
	}
	return this
}

func (this *NoticeQuery) Lt(field string, value interface{}) *NoticeQuery {
	this.Op("lt", field, value)
	return this
}

func (this *NoticeQuery) Lte(field string, value interface{}) *NoticeQuery {
	this.Op("lte", field, value)
	return this
}

func (this *NoticeQuery) Gt(field string, value interface{}) *NoticeQuery {
	this.Op("gt", field, value)
	return this
}

func (this *NoticeQuery) Gte(field string, value interface{}) *NoticeQuery {
	this.Op("gte", field, value)
	return this
}

func (this *NoticeQuery) Action(action NoticeQueryAction, ) *NoticeQuery {
	this.action = action
	return this
}

func (this *NoticeQuery) For(field string) *NoticeQuery {
	this.forField = field
	return this
}

// 开始执行
func (this *NoticeQuery) Execute() (interface{}, error) {
	collectionName := "notices"
	if this.action == NoticeQueryActionFindAll {
		result := []*notices.Notice{}
		ones, err := this.findAll(collectionName)
		if err != nil {
			return nil, err
		}
		result = append(result, ones ...)
		return result, nil
	} else if this.action == NoticeQueryActionFind {
		result := []*notices.Notice{}
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
func (this *NoticeQuery) Find() (*notices.Notice, error) {
	result, err := this.Action(NoticeQueryActionFind).Execute()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*notices.Notice), nil
}

// 查找多个数据
func (this *NoticeQuery) FindAll() ([]*notices.Notice, error) {
	ones, err := this.Action(NoticeQueryActionFindAll).Execute()
	if err != nil {
		return nil, err
	}
	return ones.([]*notices.Notice), err
}

// 插入新数据
func (this *NoticeQuery) Insert(notice *notices.Notice) error {
	if notice == nil {
		return errors.New("value should not be nil")
	}

	if notice.Id.IsZero() {
		notice.Id = primitive.NewObjectID()
	}

	collectionName := "notices"
	coll := teamongo.FindCollection(collectionName)
	_, err := coll.InsertOne(context.Background(), *notice)
	return err
}

// 删除数据
func (this *NoticeQuery) Delete() error {
	filter := this.buildFilter()

	collectionName := "notices"
	coll := teamongo.FindCollection(collectionName)
	_, err := coll.DeleteMany(context.Background(), filter)
	return err
}

// 修改数据
func (this *NoticeQuery) Update(values maps.Map) error {
	if values.Len() == 0 {
		return nil
	}

	collectionName := "notices"
	coll := teamongo.FindCollection(collectionName)

	filter := this.buildFilter()
	_, err := coll.UpdateMany(context.Background(), filter, values)
	return err
}

// 计算数量
func (this *NoticeQuery) Count() (count int64, err error) {
	c, err := this.Action(NoticeQueryActionCount).Execute()
	if err != nil {
		return 0, err
	}
	return types.Int64(c), err
}

func (this *NoticeQuery) queryNumber(collectionName string) (float64, error) {
	if this.action == NoticeQueryActionCount {
		coll := teamongo.FindCollection(collectionName)
		filter := this.buildFilter()
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		i, err := coll.CountDocuments(ctx, filter)
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

func (this *NoticeQuery) queryGroup(collectionName string) (result map[string]map[string]interface{}, err error) {
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
	if this.action == NoticeQueryActionCount {
		countField = map[string]interface{}{
			"$sum": 1,
		}
	} else if this.action == NoticeQueryActionMin {
		if len(this.forField) == 0 {
			return nil, errors.New("should specify field for the action")
		}
		countField = map[string]interface{}{
			"$min": "$" + this.forField,
		}
	} else if this.action == NoticeQueryActionMax {
		if len(this.forField) == 0 {
			return nil, errors.New("should specify field for the action")
		}
		countField = map[string]interface{}{
			"$max": "$" + this.forField,
		}
	} else if this.action == NoticeQueryActionAvg {
		if len(this.forField) == 0 {
			return nil, errors.New("should specify field for the action")
		}
		countField = map[string]interface{}{
			"$avg": "$" + this.forField,
		}
	} else if this.action == NoticeQueryActionSum {
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

	pipelines, err := teamongo.BSONArrayBytes([]byte(`[
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

	cursor, err := teamongo.FindCollection(collectionName).Aggregate(context.Background(), pipelines)
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

func (this *NoticeQuery) findAll(collectionName string) (result []*notices.Notice, err error) {
	coll := teamongo.FindCollection(collectionName)
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

	result = []*notices.Notice{}
	for cursor.Next(context.Background()) {
		m := &notices.Notice{}
		err := cursor.Decode(m)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}

func (this *NoticeQuery) buildFilter() map[string]interface{} {
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

	// proxy
	if this.proxyCond != nil {
		if len(this.proxyCond.ServerId) > 0 {
			filter["proxy.serverId"] = this.proxyCond.ServerId
		}
		if len(this.proxyCond.LocationId) > 0 {
			filter["proxy.locationId"] = this.proxyCond.LocationId
		}
		if len(this.proxyCond.RewriteId) > 0 {
			filter["proxy.rewriteId"] = this.proxyCond.RewriteId
		}
		if len(this.proxyCond.FastcgiId) > 0 {
			filter["proxy.fastcgiId"] = this.proxyCond.FastcgiId
		}
		if len(this.proxyCond.BackendId) > 0 {
			filter["proxy.backendId"] = this.proxyCond.BackendId
		}
	}

	// agent
	if this.agentCond != nil {
		if len(this.agentCond.AgentId) > 0 {
			filter["agent.agentId"] = this.agentCond.AgentId
		}
		if len(this.agentCond.AppId) > 0 {
			filter["agent.appId"] = this.agentCond.AppId
		}
		if len(this.agentCond.TaskId) > 0 {
			filter["agent.taskId"] = this.agentCond.TaskId
		}
		if len(this.agentCond.ItemId) > 0 {
			filter["agent.itemId"] = this.agentCond.ItemId
		}
	}

	if this.debug {
		logs.PrintAsJSON(filter)
	}

	return filter
}

func (this *NoticeQuery) jsonEncode(i interface{}) (string, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (this *NoticeQuery) jsonEncodeString(i interface{}) string {
	data, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(data)
}
