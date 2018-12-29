package teaconfigs

import (
	"github.com/go-yaml/yaml"
	"testing"
)

func TestServerBackendConfig(t *testing.T) {
	yamlData, err := yaml.Marshal(new(ServerBackendConfig))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(yamlData))
}
