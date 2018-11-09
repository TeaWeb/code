package teastats

import "testing"

func TestAPIDailyRequestsStat_ListLatestDays(t *testing.T) {
	stat := new(APIDailyRequestsStat)
	for _, v := range stat.ListLatestDays("lb001", 10) {
		t.Log(v)
	}
}

func TestAPIDailyRequestsStat_ListLatestDaysForAPI(t *testing.T) {
	stat := new(APIDailyRequestsStat)
	for _, v := range stat.ListLatestDaysForAPI("lb001", "/user/:id", 10) {
		t.Log(v)
	}
}
