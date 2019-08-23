package db

import (
	"testing"
)

func TestSharedMongoConfig(t *testing.T) {
	t.Logf("%#v", SharedMongoConfig())
}
