package teastats

import (
	"testing"
)

func TestFindFilter(t *testing.T) {
	t.Log(FindFilter("request.all.second"))
}

func TestRestartServerFilters(t *testing.T) {
	RestartServerFilters("123456", []string{"request.all.second", "request.all.minute", "request.all.minute"})
	RestartServerFilters("123456", []string{"request.all.second", "request.all.minute", "pv.all.second"})
}
