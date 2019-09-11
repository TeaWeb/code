package teadb

import (
	"database/sql"
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestMySQLDriver_Open(t *testing.T) {
	dbInstance, err := sql.Open("mysql", "root:abcdef@tcp(127.0.0.1:3306)/teaweb123?charset=utf8mb4&timeout=30s")
	if err != nil {
		t.Log("error:", err.Error())
		return
	}
	_ = dbInstance.Close()
	t.Log("ok")
}

func TestMySQLDriver_buildWhere(t *testing.T) {
	q := NewQuery("myTable")
	q.Attr("name", "lu")
	q.Attr("age", 10)
	q.Gt("timestamp", "gt")
	q.Lt("timestamp", "lt")
	q.Gte("timestamp", "gte")
	q.Lte("timestamp", "lte")
	q.Not("timestamp", "not")
	q.Attr("a", []string{"a", "b", "c"})
	q.Attr("timestamp", nil)
	q.Or([]OperandMap{
		{
			"timestamp": {
				{
					Code:  OperandEq,
					Value: "123",
				},
			},
		},
		{
			"timestamp": {
				{
					Code:  OperandGt,
					Value: "456",
				},
				{
					Code:  OperandNotIn,
					Value: []int{1, 2, 3},
				},
			},
		},
	})

	driver := new(MySQLDriver)
	driver.Init()
	paramsHolder := NewSQLParamsHolder(driver.driver)
	where, err := driver.buildWhere(q.operandMap, nil, paramsHolder)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("where:", where)
	logs.PrintAsJSON(paramsHolder.Params, t)
}

func TestMySQLDriver_buildWhere_Or(t *testing.T) {
	q := NewQuery("myTable")
	q.Or([]OperandMap{
		{
			"timestamp": {
				{
					Code:  OperandEq,
					Value: "123",
				},
			},
		},
		{
			"timestamp": {
				{
					Code:  OperandGt,
					Value: "456",
				},
				{
					Code:  OperandNotIn,
					Value: []int{1, 2, 3},
				},
			},
		},
		{
			"timestamp": {
				{
					Code:  OperandLt,
					Value: 1024,
				},
			},
		},
	})

	driver := new(MySQLDriver)
	paramsHolder := NewSQLParamsHolder(driver.driver)
	where, err := driver.buildWhere(q.operandMap, nil, paramsHolder)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("where:", where)
	logs.PrintAsJSON(paramsHolder.Params, t)
}

func TestMySQLDriver_asSQL_SELECT(t *testing.T) {
	{
		q := NewQuery("myTable")
		q.Attr("name", "lu")
		q.Attr("age", 10)

		driver := new(MySQLDriver)
		s, err := driver.asSQL(SQLSelect, q, NewSQLParamsHolder(driver.driver), "", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(s)
	}

	{
		q := NewQuery("myTable")
		q.Result("a", "b")
		q.Attr("name", "lu")
		q.Attr("age", 10)
		q.Offset(20)
		q.Limit(10)
		q.Desc("_id")
		q.Asc("createdAt")

		driver := new(MySQLDriver)
		s, err := driver.asSQL(SQLSelect, q, NewSQLParamsHolder(driver.driver), "", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(s)
	}
}

func TestMySQLDriver_asSQL_DELETE(t *testing.T) {
	{
		q := NewQuery("myTable")
		q.Attr("name", "lu")
		q.Attr("age", 10)

		driver := new(MySQLDriver)
		s, err := driver.asSQL(SQLDelete, q, NewSQLParamsHolder(driver.driver), "", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(s)
	}

	{
		q := NewQuery("myTable")
		q.Result("a", "b")
		q.Attr("name", "lu")
		q.Attr("age", 10)
		q.Offset(20)
		q.Limit(10)
		q.Desc("_id")
		q.Asc("createdAt")

		driver := new(MySQLDriver)
		s, err := driver.asSQL(SQLDelete, q, NewSQLParamsHolder(driver.driver), "", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(s)
	}
}

func TestMySQLDriver_asSQL_Update(t *testing.T) {
	{
		q := NewQuery("myTable")
		q.Attr("name", "lu")
		q.Attr("age", 10)

		driver := new(MySQLDriver)
		s, err := driver.asSQL(SQLUpdate, q, NewSQLParamsHolder(driver.driver), "", map[string]interface{}{
			"a": 1,
			"b": 2,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(s)
	}

	{
		q := NewQuery("myTable")
		q.Result("a", "b")
		q.Attr("name", "lu")
		q.Attr("age", 10)
		q.Offset(20)
		q.Limit(10)
		q.Desc("_id")
		q.Asc("createdAt")

		driver := new(MySQLDriver)
		s, err := driver.asSQL(SQLUpdate, q, NewSQLParamsHolder(driver.driver), "", map[string]interface{}{
			"a": 1,
			"b": 2,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(s)
	}
}

func TestMySQLDriver_TestDSN(t *testing.T) {
	driver := new(MySQLDriver)
	{
		message, ok := driver.TestDSN("root:abcdef@tcp(127.0.0.1:3306)/teaweb123?charset=utf8mb4&timeout=30s")
		t.Log(message, ok)
	}
	{
		message, ok := driver.TestDSN("root:123456@tcp(127.0.0.1:3306)/teaweb123?charset=utf8mb4&timeout=30s")
		t.Log(message, ok)
	}
	{
		message, ok := driver.TestDSN("root:123456@tcp(127.0.0.1:3306)/teaweb?charset=utf8mb4&timeout=30s")
		t.Log(message, ok)
	}
}

func TestMySQLDriver_Ping(t *testing.T) {
	driver := new(MySQLDriver)
	err := driver.initDB()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(driver.Test())
}
