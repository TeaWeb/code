package teadb

import (
	"github.com/iwind/TeaGo/logs"
	"golang.org/x/net/context"
	"testing"
)

func TestPostgresDriver_CheckTableExists(t *testing.T) {
	driver := new(PostgresDriver)
	driver.Init()
	{
		ok, err := driver.CheckTableExists("teaweb.logs.audit")
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Log("found")
		} else {
			t.Log("not found")
		}
	}

	{
		ok, err := driver.CheckTableExists("teaweb.logs.audit123")
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Log("found")
		} else {
			t.Log("not found")
		}
	}
}

func TestPostgresDriver_CreateTable(t *testing.T) {
	driver := new(PostgresDriver)
	err := driver.initDB()
	if err != nil {
		t.Fatal(err)
	}
	currentDB, err := driver.checkDB()
	if err != nil {
		t.Fatal(err)
	}
	logs.Println("create table")
	_, err = currentDB.ExecContext(context.Background(), `CREATE TABLE "public"."a" (
			"id" serial8 primary key,
			"_id" varchar(24)
		);
	`)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		logs.Println("drop table")
		_, err = currentDB.ExecContext(context.Background(), "DROP TABLE \"a\"")
		if err != nil {
			t.Fatal(err)
		}
		logs.Println("drop table end")
	}()

	logs.Println("query table")
	stmt, err := currentDB.PrepareContext(context.Background(), "SELECT * FROM \"a\"")
	if err != nil {
		t.Fatal(err)
	}

	rows, err := stmt.Query()
	if err != nil {
		t.Fatal(err)
	}

	_ = rows.Close()
	_ = stmt.Close()
}

func TestPostgresDriver_TestDSN(t *testing.T) {
	driver := new(PostgresDriver)
	{
		message, ok := driver.TestDSN("postgres://postgres:@127.0.0.1:5432/teaweb?sslmode=disable")
		t.Log(message, ok)
	}
	{
		message, ok := driver.TestDSN("postgres://postgres:123456@127.0.0.1:5432/teaweb123?sslmode=disable")
		t.Log(message, ok)
	}
	{
		message, ok := driver.TestDSN("postgres://postgres:123456@127.0.0.1:5432/teaweb?sslmode=disable")
		t.Log(message, ok)
	}
}
