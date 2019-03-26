package tealogs

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"time"
)

type QueryDuration = string

const (
	QueryDurationYearly   = "yearly"
	QueryDurationMonthly  = "monthly"
	QueryDurationDaily    = "daily"
	QueryDurationHourly   = "hourly"
	QueryDurationMinutely = "minutely"
	QueryDurationSecondly = "secondly"
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

type Query struct {
	action   string
	timeFrom time.Time
	timeTo   time.Time
	group    []string
	cond     map[string]interface{}
	duration QueryDuration
	forField string

	sorts  []map[string]int
	offset int64
	size   int64

	result []string

	debug bool
}

func NewQuery() *Query {
	return &Query{
		cond:   map[string]interface{}{},
		sorts:  []map[string]int{},
		offset: -1,
		size:   -1,
	}
}

func (this *Query) Debug() *Query {
	this.debug = true
	return this
}

func (this *Query) From(time time.Time) *Query {
	this.timeFrom = time
	return this
}

func (this *Query) To(time time.Time) *Query {
	this.timeTo = time
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

func (this *Query) Duration(duration QueryDuration) *Query {
	this.duration = duration
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
	// 时间段
	collectionNames := []string{}
	if this.timeFrom.Year() < 2000 {
		this.timeFrom = time.Now()
	}
	if this.timeTo.Year() < 2000 {
		this.timeTo = this.timeFrom
	}
	if this.timeTo.Before(this.timeFrom) {
		return nil, errors.New("timeTo should be after timeFrom")
	}

	collectionNames = append(collectionNames, "logs."+timeutil.Format("Ymd", this.timeFrom))
	startTime := this.timeFrom
	for {
		startTime = startTime.AddDate(0, 0, 1)
		if startTime.After(this.timeTo) {
			break
		}
		collectionNames = append(collectionNames, "logs."+timeutil.Format("Ymd", startTime))
	}

	if len(collectionNames) == 0 {
		return nil, errors.New("no data")
	}

	if this.action == QueryActionFindAll {
		result := []*AccessLog{}
		for _, collectionName := range collectionNames {
			ones, err := this.findAll(collectionName)
			if err != nil {
				return nil, err
			}
			result = append(result, ones ...)
		}
		return result, nil
	} else if this.action == QueryActionFind {
		result := []*AccessLog{}
		for _, collectionName := range collectionNames {
			ones, err := this.findAll(collectionName)
			if err != nil {
				return nil, err
			}
			result = append(result, ones ...)
		}
		if len(result) == 0 {
			return nil, nil
		}
		return result[0], nil
	} else if len(this.group) > 0 { // 按某个字段分组
		result := map[string]map[string]interface{}{}
		for _, collectionName := range collectionNames {
			n, err := this.queryGroup(collectionName)
			if err != nil {
				return nil, err
			}
			for k, v := range n {
				old, found := result[k]
				if found {
					result[k]["count"] = types.Float64(v["count"]) + types.Float64(old["count"])
				} else {
					result[k] = v
				}
			}
		}
		return result, nil
	} else if len(this.duration) > 0 { // 按时间周期归类 { duration => count }
		result := map[string]float64{}
		for _, collectionName := range collectionNames {
			n, err := this.queryDuration(collectionName)
			if err != nil {
				return nil, err
			}
			for durationName, v := range n {
				_, found := result[durationName]
				if found {
					result[durationName] += v
				} else {
					result[durationName] = v
				}
			}
		}
		return result, nil
	} else { // count
		result := float64(0)
		for _, collectName := range collectionNames {
			n, err := this.queryNumber(collectName)
			if err != nil {
				return nil, err
			}
			result += n
		}
		return result, nil
	}

	return nil, nil
}

// 查找单个数据
func (this *Query) Find() (*AccessLog, error) {
	result, err := this.Limit(1).Action(QueryActionFind).Execute()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*AccessLog), nil
}

// 查找多个数据
func (this *Query) FindAll() ([]*AccessLog, error) {
	result, err := this.Action(QueryActionFindAll).Execute()
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.([]*AccessLog), nil
}

func (this *Query) queryNumber(collectionName string) (float64, error) {
	if this.action == QueryActionCount {
		coll := teamongo.FindCollection(collectionName)
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

func (this *Query) queryDuration(collectionName string) (result map[string]float64, err error) {
	result = map[string]float64{}

	var groupId interface{}
	if this.duration == QueryDurationYearly {
		groupId = "$timeFormat.year"
	} else if this.duration == QueryDurationMonthly {
		groupId = "$timeFormat.month"
	} else if this.duration == QueryDurationDaily {
		groupId = "$timeFormat.day"
	} else if this.duration == QueryDurationHourly {
		groupId = "$timeFormat.hour"
	} else if this.duration == QueryDurationMinutely {
		groupId = "$timeFormat.minute"
	} else if this.duration == QueryDurationSecondly {
		groupId = "$timeFormat.second"
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

	// filter
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
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		m := map[string]interface{}{}
		err := cursor.Decode(&m)
		if err != nil {
			return nil, err
		}
		//logs.Println(m)
		result[types.String(m["_id"])] = types.Float64(m["count"])
	}

	return
}

func (this *Query) queryGroup(collectionName string) (result map[string]map[string]interface{}, err error) {
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
	defer cursor.Close(context.Background())

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

func (this *Query) findAll(collectionName string) (result []*AccessLog, err error) {
	coll := teamongo.FindCollection(collectionName)
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
	cursor, err := coll.Find(context.Background(), this.buildFilter(), opts ...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	result = []*AccessLog{}
	for cursor.Next(context.Background()) {
		m := &AccessLog{}
		err := cursor.Decode(m)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}

func (this *Query) buildFilter() map[string]interface{} {
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
