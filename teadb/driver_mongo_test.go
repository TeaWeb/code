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
	filter, err := driver.buildFilter(q)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(filter, t)
}

func TestMongoDriver_setMapValue(t *testing.T) {
	m := maps.Map{}

	driver := new(MongoDriver)
	driver.setMapValue(m, []string{"a", "b", "c", "d", "e"}, 123)
	logs.PrintAsJSON(m, t)
}

func TestMongoDriver_connect(t *testing.T) {
	driver := new(MongoDriver)
	client, err := driver.connect()
	if err != nil {
		t.Log("ERROR:", err.Error())
		return
	}
	t.Log("client:", client)
}

func TestMongoDriver_Test(t *testing.T) {
	driver := new(MongoDriver)
	err := driver.Test()
	if err != nil {
		t.Log("ERROR:", err.Error())
		return
	}
	t.Log("client:", driver)
}
