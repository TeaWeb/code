package tealogs

import (
	"context"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/time"
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

func TestQuery_UpgradeDate(t *testing.T) {
	timeFrom := time.Now().AddDate(0, 0, -60)
	timeTo := time.Now()

	days := []string{}
	for {
		if !timeTo.After(timeFrom) {
			break
		}
		days = append(days, timeutil.Format("Ymd", timeFrom))
		timeFrom = timeFrom.AddDate(0, 0, 1)
	}

	for _, day := range days {
		logs.Println(day)
		coll := teamongo.FindCollection("logs." + day)
		cursor, err := coll.Find(context.Background(), map[string]interface{}{})
		if err != nil {
			t.Fatal(err)
		}
		for cursor.Next(context.Background()) {
			m := map[string]interface{}{}
			cursor.Decode(&m)
			timestamp := types.Int64(m["timestamp"])
			now := time.Unix(timestamp, 0)
			_, err := coll.UpdateOne(context.Background(), map[string]interface{}{
				"_id": m["_id"],
			}, map[string]interface{}{
				"$unset": map[string]interface{}{
					"date": "",
				},
				"$set": map[string]interface{}{
					"time": time.Unix(timestamp, 0),
					"timeFormat": map[string]interface{}{
						"year":   timeutil.Format("Y", now),
						"month":  timeutil.Format("Ym", now),
						"day":    timeutil.Format("Ymd", now),
						"hour":   timeutil.Format("YmdH", now),
						"minute": timeutil.Format("YmdHi", now),
						"second": timeutil.Format("YmdHis", now),
					},
				},
			})
			if err != nil {
				cursor.Close(context.Background())
				t.Fatal(err)
			}
		}
		cursor.Close(context.Background())
	}

	t.Log(days)
	t.Log("finished")
}
