package tealogs

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
	"time"
)

func TestQuery_Count(t *testing.T) {
	query := NewQuery()
	query.From(time.Now().AddDate(0, 0, -6))
	query.To(time.Now())
	query.Attr("scheme", []string{"http", "https"})
	query.Gte("status", 404)
	//query.Debug()
	v, err := query.Action(QueryActionCount).Execute()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(v)
}

func TestQuery_DurationPerformance(t *testing.T) {
	query := NewQuery()
	query.From(time.Date(2018, 11, 22, 0, 0, 0, 0, time.Local))
	query.To(time.Date(2018, 11, 22, 23, 59, 59, 0, time.Local))
	query.Attr("timeFormat.hour", "2018112222")
	query.Attr("timeFormat.minute", "201811222246")
	//query.Attr("scheme", []string{"http", "https"})
	query.Debug()
	query.Duration(QueryDurationSecondly)
	v, err := query.Action(QueryActionCount).Execute()
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(v)
}

func TestQuery_DurationHourly(t *testing.T) {
	query := NewQuery()
	query.From(time.Date(2018, 11, 22, 0, 0, 0, 0, time.Local))
	query.To(time.Now())
	//query.Attr("scheme", []string{"http", "https"})
	//query.Debug()
	query.Duration(QueryDurationHourly)
	v, err := query.Action(QueryActionCount).Execute()
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(v)
}

func TestQuery_Duration(t *testing.T) {
	query := NewQuery()
	query.From(time.Now().AddDate(0, 0, -1))
	query.To(time.Now())
	//query.Attr("scheme", []string{"http", "https"})
	//query.Debug()
	query.Duration(QueryDurationHourly)
	v, err := query.Action(QueryActionSum).For("requestTime").Execute()
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(v)
}

func TestQuery_Group(t *testing.T) {
	query := NewQuery()
	query.From(time.Date(2019, 1, 2, 0, 0, 0, 0, time.Local))
	query.To(time.Date(2019, 1, 2, 23, 59, 59, 0, time.Local))
	query.Attr("serverId", "VEQ6mBKq7w7lFUzj")
	query.Gt("status", 0)
	v, err := query.Action(QueryActionCount).Group([]string{"status"}).Execute()
	if err != nil {
		t.Fatal(err)
	}

	logs.PrintAsJSON(v)
}
