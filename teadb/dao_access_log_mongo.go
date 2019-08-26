package teadb

import (
	"context"
	"github.com/TeaWeb/code/teadb/shared"
	"github.com/TeaWeb/code/tealogs/accesslogs"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"strings"
	"time"
)

type MongoAccessLogDAO struct {
}

func (this *MongoAccessLogDAO) Init() {

}

func (this *MongoAccessLogDAO) TableName(day string) string {
	return "logs." + day
}

// 写入一条日志
func (this *MongoAccessLogDAO) InsertOne(accessLog *accesslogs.AccessLog) error {
	if accessLog.Id.IsZero() {
		accessLog.Id = shared.NewObjectId()
	}
	return NewQuery(this.TableName(timeutil.Format("Ymd"))).
		InsertOne(accessLog)
}

// 写入一组日志
func (this *MongoAccessLogDAO) InsertAccessLogs(accessLogList []interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := teamongo.SharedCollection(this.TableName(timeutil.Format("Ymd"))).
		InsertMany(ctx, accessLogList)
	return err
}

func (this *MongoAccessLogDAO) FindAccessLogCookie(day string, logId string) (*accesslogs.AccessLog, error) {
	idObject, err := shared.ObjectIdFromHex(logId)
	if err != nil {
		return nil, err
	}

	one, err := NewQuery(this.TableName(day)).
		Attr("_id", idObject).
		Result("cookie").
		FindOne(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*accesslogs.AccessLog), nil
}

func (this *MongoAccessLogDAO) FindRequestHeaderAndBody(day string, logId string) (*accesslogs.AccessLog, error) {
	idObject, err := shared.ObjectIdFromHex(logId)
	if err != nil {
		return nil, err
	}
	one, err := NewQuery(this.TableName(day)).
		Attr("_id", idObject).
		Result("header", "requestData").
		FindOne(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*accesslogs.AccessLog), nil
}

func (this *MongoAccessLogDAO) FindResponseHeaderAndBody(day string, logId string) (*accesslogs.AccessLog, error) {
	idObject, err := shared.ObjectIdFromHex(logId)
	if err != nil {
		return nil, err
	}
	one, err := NewQuery(this.TableName(day)).
		Attr("_id", idObject).
		Result("sentHeader", "responseBodyData").
		FindOne(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*accesslogs.AccessLog), nil
}

func (this *MongoAccessLogDAO) ListAccessLogs(day string, serverId string, fromId string, onlyErrors bool, searchIP string, offset int, size int) ([]*accesslogs.AccessLog, error) {
	query := NewQuery(this.TableName(day))
	query.Attr("serverId", serverId)
	if len(fromId) > 0 {
		fromIdObject, err := shared.ObjectIdFromHex(fromId)
		if err != nil {
			return nil, err
		}
		query.Lt("_id", fromIdObject)
	}
	if onlyErrors {
		query.Or([]OperandMap{
			{
				"hasErrors": []*Operand{NewOperand(OperandEq, true)},
			},
			{
				"status": []*Operand{NewOperand(OperandGte, 400)},
			},
		})
	}
	if len(searchIP) > 0 {
		query.Attr("remoteAddr", searchIP)
	}
	query.Offset(offset)
	query.Limit(size)
	query.Desc("_id")
	ones, err := query.FindOnes(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}
	return result, nil
}

func (this *MongoAccessLogDAO) HasNextAccessLog(day string, serverId string, fromId string, onlyErrors bool, searchIP string) (bool, error) {
	query := NewQuery(this.TableName(day))
	query.Attr("serverId", serverId).
		Result("_id")
	if len(fromId) > 0 {
		fromIdObject, err := shared.ObjectIdFromHex(fromId)
		if err != nil {
			return false, err
		}
		query.Lt("_id", fromIdObject)
	}
	if onlyErrors {
		query.Or([]OperandMap{
			{
				"hasErrors": []*Operand{NewOperand(OperandEq, true)},
			},
			{
				"status": []*Operand{NewOperand(OperandGte, 400)},
			},
		})
	}
	if len(searchIP) > 0 {
		query.Attr("remoteAddr", searchIP)
	}

	one, err := query.FindOne(new(accesslogs.AccessLog))
	if err != nil {
		return false, err
	}
	return one != nil, nil
}

func (this *MongoAccessLogDAO) HasAccessLog(day string, serverId string) (bool, error) {
	query := NewQuery(this.TableName(day))
	one, err := query.Attr("serverId", serverId).
		Result("_id").
		FindOne(new(accesslogs.AccessLog))
	return one != nil, err
}

func (this *MongoAccessLogDAO) ListLatestAccessLogs(day string, serverId string, fromId string, onlyErrors bool, size int) ([]*accesslogs.AccessLog, error) {
	query := NewQuery(this.TableName(day))

	shouldReverse := true
	query.Attr("serverId", serverId)
	if len(fromId) > 0 {
		fromIdObject, err := shared.ObjectIdFromHex(fromId)
		if err != nil {
			return nil, err
		}
		query.Gt("_id", fromIdObject)
		query.Asc("_id")
	} else {
		query.Desc("_id")
		shouldReverse = false
	}
	if onlyErrors {
		query.Or([]OperandMap{
			{
				"hasErrors": []*Operand{NewOperand(OperandEq, true)},
			},
			{
				"status": []*Operand{NewOperand(OperandGte, 400)},
			},
		})
	}
	query.Limit(size)
	ones, err := query.FindOnes(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}

	if shouldReverse {
		lists.Reverse(ones)
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}

	return result, nil
}

func (this *MongoAccessLogDAO) ListTopAccessLogs(day string, size int) ([]*accesslogs.AccessLog, error) {
	ones, err := NewQuery(this.TableName(day)).
		Limit(size).
		Desc("_id").
		FindOnes(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}
	return result, nil
}

func (this *MongoAccessLogDAO) QueryAccessLogs(day string, serverId string, query *Query) ([]*accesslogs.AccessLog, error) {
	query.table = this.TableName(day)
	ones, err := query.
		Attr("serverId", serverId).
		FindOnes(new(accesslogs.AccessLog))
	if err != nil {
		return nil, err
	}

	result := []*accesslogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*accesslogs.AccessLog))
	}
	return result, nil
}

func (this *MongoAccessLogDAO) initTable(day string) {
	table := this.TableName(day)
	if isInitializedTable(table) {
		return
	}
	for _, fields := range [][]*shared.IndexField{
		{
			shared.NewIndexField("serverId", true),
		},
		{
			shared.NewIndexField("status", true),
			shared.NewIndexField("serverId", true),
		},
		{
			shared.NewIndexField("remoteAddr", true),
			shared.NewIndexField("serverId", true),
		},
		{
			shared.NewIndexField("hasErrors", true),
			shared.NewIndexField("serverId", true),
		},
	} {
		err := this.createIndex(day, fields)
		if err != nil {
			logs.Error(err)
		}
	}
}

func (this *MongoAccessLogDAO) createIndex(day string, fields []*shared.IndexField) error {
	if len(fields) == 0 {
		return nil
	}

	coll := teamongo.SharedCollection(this.TableName(day))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// 创建新的
	bsonDoc := bsonx.Doc{}
	for _, field := range fields {
		if field.Asc {
			bsonDoc = bsonDoc.Append(field.Name, bsonx.Int32(1))
		} else {
			bsonDoc = bsonDoc.Append(field.Name, bsonx.Int32(-1))
		}
	}

	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bsonDoc,
		Options: options.Index().SetBackground(true),
	})

	// 忽略可能产生的冲突错误
	if err != nil && strings.Contains(err.Error(), "existing") {
		err = nil
	}

	return err
}
