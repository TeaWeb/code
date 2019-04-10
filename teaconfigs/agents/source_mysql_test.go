package agents

import (
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestMySQLSource_Execute(t *testing.T) {
	source := NewMySQLSource()
	source.TimeoutSeconds = 10
	source.Addr = "127.0.0.1"
	source.Username = "root"
	source.Password = ""
	source.DatabaseName = teamongo.DatabaseName
	source.SQL = "SELECT * FROM tea_accessLogs"
	values, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(values, t)
}
