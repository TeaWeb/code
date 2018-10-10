package configs

import "testing"

func TestSharedAdminConfig(t *testing.T) {
	adminConfig := SharedAdminConfig()
	t.Logf("%#v", adminConfig)
}
