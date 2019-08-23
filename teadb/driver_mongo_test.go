package teadb

import (
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"testing"
)

func TestMongoDriver_buildFilter(t *testing.T) {
	q := new(Query)
	q.Init()
	q.Attr("name", "lu")
	q.Op("age", OperandGt, 1024)
	q.Op("age", OperandLt, 2048)
	q.Op("count", OperandEq, 3)

	driver := new(MongoDriver)
	logs.PrintAsJSON(driver.buildFilter(q), t)
}

func TestMongoDriver_setMapValue(t *testing.T) {
	m := maps.Map{}

	driver := new(MongoDriver)
	driver.setMapValue(m, []string{"a", "b", "c", "d", "e"}, 123)
	logs.PrintAsJSON(m, t)
}
