package teaplugins

import (
	"github.com/TeaWeb/plugin/messages"
	plugins2 "github.com/TeaWeb/plugin/plugins"
	"os"
	"testing"
)

func TestLoader_Load(t *testing.T) {
	loader := NewLoader(os.Getenv("GOPATH") + "/src/github.com/TeaWeb/plugin/main/demo.plugin")
	err := loader.Load()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoader_CallAction(t *testing.T) {
	loader := NewLoader("")
	action := new(messages.RegisterPluginAction)
	action.Plugin = new(plugins2.Plugin)
	err := loader.CallAction(action)
	if err != nil {
		t.Fatal(err)
	}
}
