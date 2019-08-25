package teadb

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/lists"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoAccessLogDAO struct {
}

func (this *MongoAccessLogDAO) Init() {

}

func (this *MongoAccessLogDAO) TableName(day string) string {
	return "logs." + day
}

func (this *MongoAccessLogDAO) FindAccessLogCookie(day string, logId string) (*tealogs.AccessLog, error) {
	idObject, err := primitive.ObjectIDFromHex(logId)
	if err != nil {
		return nil, err
	}

	one, err := NewQuery(this.TableName(day)).
		Attr("_id", idObject).
		Result("cookie").
		FindOne(new(tealogs.AccessLog))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*tealogs.AccessLog), nil
}

func (this *MongoAccessLogDAO) FindRequestHeaderAndBody(day string, logId string) (*tealogs.AccessLog, error) {
	idObject, err := primitive.ObjectIDFromHex(logId)
	if err != nil {
		return nil, err
	}
	one, err := NewQuery(this.TableName(day)).
		Attr("_id", idObject).
		Result("header", "requestData").
		FindOne(new(tealogs.AccessLog))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*tealogs.AccessLog), nil
}

func (this *MongoAccessLogDAO) FindResponseHeaderAndBody(day string, logId string) (*tealogs.AccessLog, error) {
	idObject, err := primitive.ObjectIDFromHex(logId)
	if err != nil {
		return nil, err
	}
	one, err := NewQuery(this.TableName(day)).
		Attr("_id", idObject).
		Result("sentHeader", "responseBodyData").
		FindOne(new(tealogs.AccessLog))
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*tealogs.AccessLog), nil
}

func (this *MongoAccessLogDAO) ListAccessLogs(day string, serverId string, fromId string, onlyErrors bool, searchIP string, offset int, size int) ([]*tealogs.AccessLog, error) {
	query := NewQuery(this.TableName(day))
	query.Attr("serverId", serverId)
	if len(fromId) > 0 {
		fromIdObject, err := primitive.ObjectIDFromHex(fromId)
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
	ones, err := query.FindOnes(new(tealogs.AccessLog))
	if err != nil {
		return nil, err
	}

	result := []*tealogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*tealogs.AccessLog))
	}
	return result, nil
}

func (this *MongoAccessLogDAO) HasNextAccessLog(day string, serverId string, fromId string, onlyErrors bool, searchIP string) (bool, error) {
	query := NewQuery(this.TableName(day))
	query.Attr("serverId", serverId)
	if len(fromId) > 0 {
		fromIdObject, err := primitive.ObjectIDFromHex(fromId)
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

	one, err := query.FindOne(new(tealogs.AccessLog))
	if err != nil {
		return false, err
	}
	return one != nil, nil
}

func (this *MongoAccessLogDAO) ListLatestAccessLogs(day string, serverId string, fromId string, onlyErrors bool, size int) ([]*tealogs.AccessLog, error) {
	query := NewQuery(this.TableName(day))

	shouldReverse := true
	query.Attr("serverId", serverId)
	if len(fromId) > 0 {
		fromIdObject, err := primitive.ObjectIDFromHex(fromId)
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
	ones, err := query.FindOnes(new(tealogs.AccessLog))
	if err != nil {
		return nil, err
	}

	if shouldReverse {
		lists.Reverse(ones)
	}

	accessLogs := []*tealogs.AccessLog{}
	for _, one := range ones {
		accessLogs = append(accessLogs, one.(*tealogs.AccessLog))
	}

	return accessLogs, nil
}

func (this *MongoAccessLogDAO) ListTopAccessLogs(day string, size int) ([]*tealogs.AccessLog, error) {
	ones, err := NewQuery(this.TableName(day)).
		Limit(size).
		Desc("_id").
		FindOnes(new(tealogs.AccessLog))
	if err != nil {
		return nil, err
	}

	result := []*tealogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*tealogs.AccessLog))
	}
	return result, nil
}

func (this *MongoAccessLogDAO) QueryAccessLogs(day string, serverId string, query *Query) ([]*tealogs.AccessLog, error) {
	query.table = this.TableName(day)
	ones, err := query.
		Attr("serverId", serverId).
		FindOnes(new(tealogs.AccessLog))
	if err != nil {
		return nil, err
	}

	result := []*tealogs.AccessLog{}
	for _, one := range ones {
		result = append(result, one.(*tealogs.AccessLog))
	}
	return result, nil
}
