package teadb

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"time"
)

type Query struct {
	table      string
	offset     int
	size       int
	operandMap map[string][]*Operand // field => operands
	sortFields []*SortField
	debug      bool
	timeout    time.Duration
}

func NewQuery(table string) *Query {
	query := &Query{
		table: table,
	}
	query.Init()
	return query
}

func (this *Query) Init() *Query {
	this.offset = -1
	this.size = -1
	this.operandMap = map[string][]*Operand{}
	return this
}

func (this *Query) Table(table string) *Query {
	this.table = table
	return this
}

func (this *Query) Debug() *Query {
	this.debug = true
	return this
}

func (this *Query) Timeout(timeout time.Duration) *Query {
	this.timeout = timeout
	return this
}

func (this *Query) Attr(field string, value interface{}) *Query {
	if types.IsSlice(value) {
		this.Op(field, OperandIn, value)
	} else {
		this.Op(field, OperandEq, value)
	}
	return this
}

func (this *Query) Op(field string, operandCode OperandCode, value interface{}) *Query {
	operands, ok := this.operandMap[field]
	if ok {
		this.operandMap[field] = append(operands, NewOperand(operandCode, value))
	} else {
		this.operandMap[field] = []*Operand{NewOperand(operandCode, value)}
	}
	return this
}

func (this *Query) Not(field string, value interface{}) *Query {
	if types.IsSlice(value) {
		this.Op(field, OperandNotIn, value)
	} else {
		this.Op(field, OperandEq, value)
	}
	return this
}

func (this *Query) Lt(field string, value interface{}) *Query {
	this.Op(field, OperandLt, value)
	return this
}

func (this *Query) Lte(field string, value interface{}) *Query {
	this.Op(field, OperandLte, value)
	return this
}

func (this *Query) Gt(field string, value interface{}) *Query {
	this.Op(field, OperandGt, value)
	return this
}

func (this *Query) Gte(field string, value interface{}) *Query {
	this.Op(field, OperandGte, value)
	return this
}

func (this *Query) Asc(field string) *Query {
	if this.hasSortField(field) {
		this.removeSortField(field)
	}
	this.sortFields = append(this.sortFields, &SortField{
		Name: field,
		Type: SortAsc,
	})
	return this
}

func (this *Query) Desc(field string) *Query {
	if this.hasSortField(field) {
		this.removeSortField(field)
	}
	this.sortFields = append(this.sortFields, &SortField{
		Name: field,
		Type: SortDesc,
	})
	return this
}

func (this *Query) Offset(offset int) *Query {
	this.offset = offset
	return this
}

func (this *Query) Limit(size int) *Query {
	this.size = size
	return this
}

func (this *Query) Node() *Query {
	node := teaconfigs.SharedNodeConfig()
	if node != nil {
		this.Attr("nodeId", node.Id)
	} else {
		this.Attr("nodeId", "")
	}
	return this
}

func (this *Query) FindOne(modelPtr interface{}) (interface{}, error) {
	return sharedDriver.FindOne(this, modelPtr)
}

func (this *Query) FindOnes(modelPtr interface{}) ([]interface{}, error) {
	return sharedDriver.FindOnes(this, modelPtr)
}

func (this *Query) InsertOne(modelPtr interface{}) error {
	return sharedDriver.InsertOne(this.table, modelPtr)
}

func (this *Query) InsertOnes(modelPtrSlice interface{}) error {
	return sharedDriver.InsertOnes(this.table, modelPtrSlice)
}

func (this *Query) Delete() error {
	return sharedDriver.DeleteOnes(this)
}

func (this *Query) Count() (int64, error) {
	return sharedDriver.Count(this)
}

func (this *Query) Sum(field string) (float64, error) {
	return sharedDriver.Sum(this, field)
}

func (this *Query) Min(field string) (float64, error) {
	return sharedDriver.Min(this, field)
}

func (this *Query) Max(field string) (float64, error) {
	return sharedDriver.Max(this, field)
}

func (this *Query) Avg(field string) (float64, error) {
	return sharedDriver.Avg(this, field)
}

func (this *Query) Group(field string, result map[string]Expr) ([]maps.Map, error) {
	return sharedDriver.Group(this, field, result)
}

func (this *Query) hasSortField(field string) bool {
	for _, f := range this.sortFields {
		if f.Name == field {
			return true
		}
	}
	return false
}

func (this *Query) removeSortField(field string) {
	result := []*SortField{}
	for _, f := range this.sortFields {
		if f.Name != field {
			result = append(result, f)
		}
	}
	this.sortFields = result
}
