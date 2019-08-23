package teadb

import (
	"errors"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"reflect"
	"strings"
	"time"
)

type MongoDriver struct {
	accessLogDAO AccessLogDAO
	agentLogDAO  AgentLogDAO
	auditLogDAO  AuditLogDAO
	noticeDAO    NoticeDAO
	valueDAO     ValueDAO
}

func (this *MongoDriver) Init() {
	this.valueDAO = new(MongoValueDAO)
	this.valueDAO.Init()
}

func (this *MongoDriver) FindOne(query *Query, modelPtr interface{}) (interface{}, error) {
	if len(query.table) == 0 {
		return nil, errors.New("'table' should not be empty")
	}

	db := this.db()
	if db == nil {
		return nil, errors.New("can not select db")
	}

	opts := options.Find()
	if query.offset > -1 {
		opts.SetSkip(int64(query.offset))
	}
	opts.SetLimit(1)

	if len(query.sortFields) > 0 {
		s := map[string]int{}
		for _, f := range query.sortFields {
			if f.IsAsc() {
				s[f.Name] = 1
			} else {
				s[f.Name] = -1
			}
		}
		opts.SetSort(s)
	}

	cursor, err := db.Collection(query.table).Find(this.timeoutContext(5*time.Second), this.buildFilter(query), opts)
	if err != nil {
		if this.isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	defer func(cursor *mongo.Cursor) {
		_ = cursor.Close(context.Background())
	}(cursor)

	if !cursor.Next(context.Background()) {
		return nil, nil
	}

	err = cursor.Decode(modelPtr)
	if err != nil {
		if this.isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return modelPtr, nil
}

func (this *MongoDriver) FindOnes(query *Query, modelPtr interface{}) ([]interface{}, error) {
	if len(query.table) == 0 {
		return nil, errors.New("'table' should not be empty")
	}

	db := this.db()
	if db == nil {
		return nil, errors.New("can not select db")
	}

	// 查询选项
	opts := options.Find()
	if query.offset > -1 {
		opts.SetSkip(int64(query.offset))
	}
	if query.size > -1 {
		opts.SetLimit(int64(query.size))
	}

	if len(query.sortFields) > 0 {
		s := map[string]int{}
		for _, f := range query.sortFields {
			if f.IsAsc() {
				s[f.Name] = 1
			} else {
				s[f.Name] = -1
			}
		}
		opts.SetSort(s)
	}

	cursor, err := db.Collection(query.table).Find(this.timeoutContext(5*time.Second), this.buildFilter(query), opts)
	if err != nil {
		if this.isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	defer func(cursor *mongo.Cursor) {
		_ = cursor.Close(context.Background())
	}(cursor)

	modelType := reflect.TypeOf(modelPtr).Elem()
	result := []interface{}{}
	for cursor.Next(context.Background()) {
		m := reflect.New(modelType).Interface()
		err = cursor.Decode(m)
		if err != nil {
			if this.isNotFoundError(err) {
				continue
			}
			return nil, err
		}

		result = append(result, m)
	}
	return result, nil
}

func (this *MongoDriver) DeleteOnes(query *Query) error {
	if len(query.table) == 0 {
		return errors.New("'table' should not be empty")
	}

	_, err := this.db().Collection(query.table).DeleteMany(this.timeoutContext(5*time.Second), this.buildFilter(query))
	return err
}

func (this *MongoDriver) InsertOne(table string, modelPtr interface{}) error {
	if len(table) == 0 {
		return errors.New("'table' should not be empty")
	}
	if modelPtr == nil {
		return errors.New("insertOne: modelPtr should not be nil")
	}
	_, err := this.db().Collection(table).InsertOne(this.timeoutContext(5*time.Second), modelPtr)
	return err
}

func (this *MongoDriver) InsertOnes(table string, modelPtrSlice interface{}) error {
	if len(table) == 0 {
		return errors.New("'table' should not be empty")
	}
	if modelPtrSlice == nil {
		return nil
	}

	t := reflect.ValueOf(modelPtrSlice)
	if t.IsNil() {
		return nil
	}
	if t.Kind() != reflect.Slice {
		return errors.New("insertOnes: only slice is accepted")
	}

	s := []interface{}{}
	l := t.Len()
	for i := 0; i < l; i++ {
		s = append(s, t.Index(i).Interface())
	}

	_, err := this.db().Collection(table).InsertMany(this.timeoutContext(5*time.Second), s)
	return err
}

func (this *MongoDriver) Count(query *Query) (int64, error) {
	if len(query.table) == 0 {
		return 0, errors.New("'table' should not be empty")
	}

	db := this.db()
	if db == nil {
		return 0, errors.New("can not select db")
	}

	// 查询选项
	opts := options.Count()
	if query.offset > -1 {
		opts.SetSkip(int64(query.offset))
	}
	if query.size > -1 {
		opts.SetLimit(int64(query.size))
	}
	return this.db().Collection(query.table).CountDocuments(this.timeoutContext(10*time.Second), this.buildFilter(query), opts)
}

func (this *MongoDriver) Sum(query *Query, field string) (float64, error) {
	return this.aggregate("sum", query, field)
}

func (this *MongoDriver) Avg(query *Query, field string) (float64, error) {
	return this.aggregate("avg", query, field)
}

func (this *MongoDriver) Min(query *Query, field string) (float64, error) {
	return this.aggregate("min", query, field)
}

func (this *MongoDriver) Max(query *Query, field string) (float64, error) {
	return this.aggregate("max", query, field)
}

func (this *MongoDriver) Group(query *Query, field string, result map[string]Expr) ([]maps.Map, error) {
	group := map[string]interface{}{
		"_id": "$" + field,
	}

	for name, expr := range result {
		// 处理点符号
		name = strings.Replace(name, ".", "@", -1)

		switch e := expr.(type) {
		case *SumExpr:
			group[name] = map[string]interface{}{
				"$sum": "$" + e.Field,
			}
		case *AvgExpr:
			group[name] = map[string]interface{}{
				"$avg": "$" + e.Field,
			}
		case *MaxExpr:
			group[name] = map[string]interface{}{
				"$max": "$" + e.Field,
			}
		case *MinExpr:
			group[name] = map[string]interface{}{
				"$min": "$" + e.Field,
			}
		case string:
			group[name] = map[string]interface{}{
				"$first": "$" + e,
			}
		}
	}

	sorts := map[string]interface{}{}
	if len(query.sortFields) > 0 {
		for _, sortField := range query.sortFields {
			if sortField.IsAsc() {
				sorts[sortField.Name] = 1
			} else {
				sorts[sortField.Name] = -1
			}
		}
	}

	pipelines := []interface{}{
		map[string]interface{}{
			"$match": this.buildFilter(query),
		},
		map[string]interface{}{
			"$group": group,
		},
	}
	if len(sorts) > 0 {
		pipelines = append(pipelines, map[string]interface{}{
			"$sort": sorts,
		})
	}

	if query.debug {
		logs.Println("===pipelines===")
		logs.PrintAsJSON(pipelines)
	}

	cursor, err := this.db().Collection(query.table).Aggregate(this.timeoutContext(30*time.Second), pipelines)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = cursor.Close(context.Background())
	}()

	ones := []maps.Map{}

	for cursor.Next(context.Background()) {
		m := maps.Map{}
		err = cursor.Decode(&m)
		if err != nil {
			return nil, err
		}

		// 处理@符号（从上面的点符号转换过来）
		for k, v := range m {
			if strings.Contains(k, "@") {
				this.setMapValue(m, strings.Split(k, "@"), v)
				delete(m, k)
			}
		}

		ones = append(ones, m)
	}

	return ones, nil
}

func (this *MongoDriver) db() *mongo.Database {
	client := teamongo.SharedClient()
	if client == nil {
		return nil
	}
	return client.Database(teamongo.DatabaseName)
}

func (this *MongoDriver) buildFilter(query *Query) (filter map[string]interface{}) {
	filter = map[string]interface{}{}
	if len(query.operandMap) > 0 {
		for field, operands := range query.operandMap {
			fieldQuery := map[string]interface{}{}
			for _, op := range operands {
				switch op.Code {
				case OperandEq:
					fieldQuery["$eq"] = op.Value
				case OperandLt:
					fieldQuery["$lt"] = op.Value
				case OperandLte:
					fieldQuery["$lte"] = op.Value
				case OperandGt:
					fieldQuery["$gt"] = op.Value
				case OperandGte:
					fieldQuery["$gte"] = op.Value
				case OperandIn:
					fieldQuery["$in"] = op.Value
				case OperandNotIn:
					fieldQuery["$nin"] = op.Value
				case OperandNeq:
					fieldQuery["$ne"] = op.Value
				}
			}
			if len(fieldQuery) > 0 {
				filter[field] = fieldQuery
			}
		}
	}
	return
}

func (this *MongoDriver) AccessLogDAO() AccessLogDAO {
	return this.accessLogDAO
}

func (this *MongoDriver) AgentLogDAO() AgentLogDAO {
	return this.agentLogDAO
}

func (this *MongoDriver) AuditLogDAO() AuditLogDAO {
	return this.auditLogDAO
}

func (this *MongoDriver) NoticeDAO() NoticeDAO {
	return this.noticeDAO
}

func (this *MongoDriver) ValueDAO() ValueDAO {
	return this.valueDAO
}

func (this *MongoDriver) isNotFoundError(err error) bool {
	return err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments
}

func (this *MongoDriver) timeoutContext(timeout time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return ctx
}

func (this *MongoDriver) aggregate(funcName string, query *Query, field string) (float64, error) {
	filter := this.buildFilter(query)

	pipelines, err := teamongo.BSONArrayBytes([]byte(`[
	{
		"$match": ` + stringutil.JSONEncode(filter) + `
	},
	{
		"$group": {
			"_id": null,
			"result": { "$` + funcName + `": ` + stringutil.JSONEncode("$"+field) + `}
		}
	}
]`))
	if err != nil {
		return 0, err
	}

	cursor, err := this.db().Collection(query.table).Aggregate(this.timeoutContext(30*time.Second), pipelines)
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = cursor.Close(context.Background())
	}()

	m := maps.Map{}
	if !cursor.Next(context.Background()) {
		return 0, nil
	}
	err = cursor.Decode(&m)
	if err != nil {
		return 0, err
	}

	return m.GetFloat64("result"), nil
}

func (this *MongoDriver) setMapValue(m maps.Map, keys []string, v interface{}) {
	l := len(keys)
	if l == 0 {
		return
	}
	if l == 1 {
		m[keys[0]] = v
		return
	}
	subM, ok := m[keys[0]]
	if ok {
		subV, ok := subM.(maps.Map)
		if ok {
			this.setMapValue(subV, keys[1:], v)
		} else {
			m[keys[0]] = maps.Map{}
			this.setMapValue(m[keys[0]].(maps.Map), keys[1:], v)
		}
	} else {
		m[keys[0]] = maps.Map{}
		this.setMapValue(m[keys[0]].(maps.Map), keys[1:], v)
	}
}
