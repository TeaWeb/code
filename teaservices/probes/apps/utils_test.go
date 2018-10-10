package apps

import "testing"

func Test_ps(t *testing.T) {
	p := ps("mysql", []string{"mysqld_safe$"}, true)
	t.Log(p)

	for _, p1 := range p {
		t.Log(p1.Cmdline())
	}
}
