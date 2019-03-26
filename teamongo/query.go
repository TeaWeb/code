package teamongo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"sync"
	"time"
)

type QueryAction = string

const (
	QueryActionCount   = "count"
	QueryActionSum     = "sum"
	QueryActionAvg     = "avg"
	QueryActionMin     = "min"
	QueryActionMax     = "max"
	QueryActionFind    = "find"
	QueryActionFindAll = "findAll"
)

var queryCollectionsMap = map[string]*Collection{}
var queryCollectionsLocker sync.Mutex

type Query struct {
	collectionName string
	modelType      reflect.Type
	action         string
	group          []string
	cond           map[string]interface{}
	orCond         []map[string]interface{}
	forField       string

	sorts  []map[string]int
	offset int64
	size   int64

	result []string

	debug bool
}

func NewQuery(collectionName string, modelPtr interface{}) *Query {
	return &Query{
		collectionName: collectionName,
		modelType:      reflect.TypeOf(modelPtr),
		cond:           map[string]interface{}{},
		sorts:          []map[string]int{},
		offset:         -1,
		size:           -1,
	}
}

func (this *Query) Debug() *Query {
	this.debug = true
	return this
}

func (this *Query) Item(itemId string) *Query {
	this.Attr("itemId", itemId)
	return this
}

func (this *Query) Asc(field string) *Query {
	if len(field) == 0 {
		field = "_id"
	}
	this.sorts = append(this.sorts, map[string]int{
		field: 1,
	})
	return this
}

func (this *Query) Desc(field string) *Query {
	if len(field) == 0 {
		field = "_id"
	}
	this.sorts = append(this.sorts, map[string]int{
		field: -1,
	})
	return this
}

func (this *Query) AscPk() *Query {
	return this.Asc("_id")
}

func (this *Query) DescPk() *Query {
	return this.Desc("_id")
}

func (this *Query) Offset(offset int64) *Query {
	this.offset = offset
	return this
}

func (this *Query) Limit(size int64) *Query {
	this.size = size
	return this
}

func (this *Query) Group(group []string) *Query {
	this.group = group
	return this
}

func (this *Query) Attr(field string, value interface{}) *Query {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		this.Op("in", field, value)
	} else {
		this.Op("eq", field, value)
	}
	return this
}

func (this *Query) Or(conds ...map[string]interface{}) *Query {
	this.orCond = conds
	return this
}

// 设置日志ID
func (this *Query) Id(idString string) *Query {
	objectId, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		logs.Error(err)
		return this.Attr("_id", idString)
	}
	this.Attr("_id", objectId)
	return this
}

func (this *Query) Op(op string, field string, value interface{}) {
	_, found := this.cond[field]
	if found {
		this.cond[field].(map[string]interface{})[op] = value
	} else {
		this.cond[field] = map[string]interface{}{
			op: value,
		}
	}
}

func (this *Query) Not(field string, value interface{}) *Query {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		this.Op("nin", field, value)
	} else {
		this.Op("ne", field, value)
	}
	return this
}

func (this *Query) Lt(field string, value interface{}) *Query {
	this.Op("lt", field, value)
	return this
}

func (this *Query) Lte(field string, value interface{}) *Query {
	this.Op("lte", field, value)
	return this
}

func (this *Query) Gt(field string, value interface{}) *Query {
	this.Op("gt", field, value)
	return this
}

func (this *Query) Gte(field string, value interface{}) *Query {
	this.Op("gte", field, value)
	return this
}

func (this *Query) Result(field ...string) *Query {
	this.result = append(this.result, field ...)
	return this
}

func (this *Query) Action(action QueryAction, ) *Query {
	this.action = action
	return this
}

func (this *Query) For(field string) *Query {
	this.forField = field
	return this
}

// 开始执行
func (this *Query) Execute() (interface{}, error) {
	if this.action == QueryActionFindAll {
		result := []interface{}{}
		ones, err := this.FindAll()
		if err != nil {
			return nil, err
		}
		result = append(result, ones ...)
		return result, nil
	} else if this.action == QueryActionFind {
		result := []interface{}{}
		ones, err := this.FindAll()
		if err != nil {
			return nil, err
		}
		result = append(result, ones ...)
		if len(result) == 0 {
			return nil, nil
		}
		return result[0], nil
	} else if len(this.group) > 0 { // 按某个字段分组
		return this.queryGroup()
	} else { // count
		return this.queryNumber()
	}
}

// 查找单个数据
func (this *Query) Find() (interface{}, error) {
	result, err := this.Limit(1).Action(QueryActionFind).Execute()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result, nil
}

// 数字
func (this *Query) Count() (int64, error) {
	count, err := this.Action(QueryActionCount).Execute()
	if err != nil {
		return 0, err
	}
	return types.Int64(count), err
}

// 插入新数据
func (this *Query) Insert(value interface{}) error {
	if value == nil {
		return errors.New("value should not be nil")
	}

	coll := this.Coll()
	_, err := coll.InsertOne(this.context(3*time.Second), value)
	return err
}

// 更新数据
func (this *Query) Update(updates maps.Map) error {
	filter := this.buildFilter()

	coll := this.Coll()
	_, err := coll.UpdateMany(context.Background(), filter, updates, nil)
	return err
}

// 删除数据
func (this *Query) Delete() error {
	filter := this.buildFilter()

	coll := this.Coll()
	_, err := coll.DeleteMany(context.Background(), filter)
	return err
}

func (this *Query) queryNumber() (float64, error) {
	if this.action == QueryActionCount {
		coll := this.Coll()
		filter := this.buildFilter()
		i, err := coll.CountDocuments(context.Background(), filter)
		if err != nil {
			return 0, err
		}
		return float64(i), nil
	} else {
		result, err := this.queryGroup()
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

func (this *Query) queryGroup() (result map[string]map[string]interface{}, err error) {
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
	if this.action == QueryActionCount {
		countField = map[string]interface{}{
			"$sum": 1,
		}
	} else if this.action == QueryActionMin {
		if len(this.forField) == 0 {
			return nil, errors.New("should specify field for the action")
		}
		countField = map[string]interface{}{
			"$min": "$" + this.forField,
		}
	} else if this.action == QueryActionMax {
		if len(this.forField) == 0 {
			return nil, errors.New("should specify field for the action")
		}
		countField = map[string]interface{}{
			"$max": "$" + this.forField,
		}
	} else if this.action == QueryActionAvg {
		if len(this.forField) == 0 {
			return nil, errors.New("should specify field for the action")
		}
		countField = map[string]interface{}{
			"$avg": "$" + this.forField,
		}
	} else if this.action == QueryActionSum {
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

	cursor, err := this.Coll().Aggregate(context.Background(), pipelines)
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

func (this *Query) FindAll() (result []interface{}, err error) {
	coll := this.Coll()
	opts := []*options.FindOptions{}
	if this.offset > -1 {
		opts = append(opts, options.Find().SetSkip(this.offset))
	}
	if this.size > -1 {
		opts = append(opts, options.Find().SetLimit(this.size))
	}
	if len(this.result) > 0 {
		projection := map[string]interface{}{}
		for _, field := range this.result {
			projection[field] = 1
		}
		opts = append(opts, options.Find().SetProjection(projection))
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

	result = []interface{}{}
	for cursor.Next(context.Background()) {
		ptrValue := reflect.New(this.modelType.Elem())
		ptr := ptrValue.Interface()
		err := cursor.Decode(ptr)
		if err != nil {
			logs.Error(err)
			continue
		}

		// TODO DECODE

		result = append(result, ptr)
	}

	return result, nil
}

// 选择集合
func (this *Query) Coll() *Collection {
	queryCollectionsLocker.Lock()
	defer queryCollectionsLocker.Unlock()

	coll, found := queryCollectionsMap[this.collectionName]
	if found {
		return coll
	}

	coll = FindCollection(this.collectionName)
	queryCollectionsMap[this.collectionName] = coll
	return coll
}

func (this *Query) buildFilter() map[string]interface{} {
	filter := map[string]interface{}{}

	// cond
	if len(this.cond) > 0 {
		for field, cond := range this.cond {
			fieldQuery := map[string]interface{}{}
			for op, value := range cond.(map[string]interface{}) {
				if field == "_id" {
					if valueString, ok := value.(string); ok {
						idValue, err := primitive.ObjectIDFromHex(valueString)
						if err == nil {
							value = idValue
						}
					}
				}
				if lists.ContainsString([]string{"eq", "lt", "lte", "gt", "gte", "in", "nin", "ne"}, op) {
					fieldQuery["$"+op] = value
				}
			}
			if len(fieldQuery) > 0 {
				filter[field] = fieldQuery
			}
		}
	}

	// or
	if len(this.orCond) > 0 {
		filter["$or"] = this.orCond
	}

	if this.debug {
		logs.PrintAsJSON(filter)
	}

	return filter
}

func (this *Query) jsonEncode(i interface{}) (string, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (this *Query) jsonEncodeString(i interface{}) string {
	data, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(data)
}

func (this *Query) context(timeout time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return ctx
}
