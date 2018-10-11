package configs

import (
	"testing"
)

func TestSharedMongoConfig(t *testing.T) {
	t.Logf("%#v", SharedMongoConfig())
}
