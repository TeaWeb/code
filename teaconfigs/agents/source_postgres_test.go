package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestSourcePostgres(t *testing.T) {
	source := NewPostgreSQLSource()
	source.Username = "postgres"
	source.Password = "123456"
	source.SQL = "SELECT * FROM settings"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}
