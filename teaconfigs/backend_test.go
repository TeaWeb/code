package teaconfigs

import (
	"github.com/go-yaml/yaml"
	"testing"
)

func TestBackendConfig(t *testing.T) {
	yamlData, err := yaml.Marshal(new(BackendConfig))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(yamlData))
}
