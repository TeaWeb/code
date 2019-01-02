package teautils

import "testing"

func TestParseVariables(t *testing.T) {
	v := ParseVariables("hello, ${name}", func(s string) string {
		return "Lu"
	})
	t.Log(v)
}
