package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/TeaWeb/code/teamongo"
	"testing"
	"time"
)

func TestCostAllPeriodFilter_Start(t *testing.T) {
	queue := new(Queue)
	queue.Start("123456")

	filter := new(CostPagePeriodFilter)
	filter.Start(queue, "cost.all.hour")

	{
		accessLog := &tealogs.AccessLog{}
		accessLog.Timestamp = time.Now().Unix()
		accessLog.RequestTime = 0.01
		filter.Filter(accessLog)
	}

	{
		accessLog := &tealogs.AccessLog{}
		accessLog.Timestamp = time.Now().Unix()
		accessLog.RequestTime = 0.02
		filter.Filter(accessLog)
	}

	{
		accessLog := &tealogs.AccessLog{}
		accessLog.Timestamp = time.Now().Unix()
		accessLog.RequestTime = 0.01
		filter.Filter(accessLog)
	}

	t.Log(filter.values)
	time.Sleep(1 * time.Second)
}

func TestCostAllPeriodFilter_FindAll(t *testing.T) {
	// init
	{
		query := teamongo.NewQuery("values.server.VEQ6mBKq7w7lFUzj", new(Value))
		query.Find()
	}

	before := time.Now()
	query := teamongo.NewQuery("values.server.VEQ6mBKq7w7lFUzj", new(Value))
	query.Attr("item", "request.all.second")
	query.Attr("timeFormat.minute", []string{
		"201903240756",
		"201903240757",
		"201903240758",
		"201903240759",
		"201903240800",
		"201903240801",
		"201903240802",
		"201903240803",
		"201903240804",
		"201903240805",
		"201903240806",
		"201903240807",
		"201903240808",
		"201903240809",
		"201903240810",
		"201903240811",
		"201903240812",
		"201903240813",
		"201903240814",
		"201903240815",
		"201903240816",
		"201903240817",
		"201903240818",
		"201903240819",
		"201903240820",
		"201903240821",
		"201903240822",
		"201903240823",
		"201903240824",
		"201903240825",
		"201903240826",
		"201903240827",
		"201903240828",
		"201903240829",
		"201903240830",
		"201903240831",
		"201903240832",
		"201903240833",
		"201903240834",
		"201903240835",
		"201903240836",
		"201903240837",
		"201903240838",
		"201903240839",
		"201903240840",
		"201903240841",
		"201903240842",
		"201903240843",
		"201903240844",
		"201903240845",
		"201903240846",
		"201903240847",
		"201903240848",
		"201903240849",
		"201903240850",
		"201903240851",
		"201903240852",
		"201903240853",
		"201903240854",
		"201903240855",
	})
	query.FindAll()
	t.Log(time.Since(before).Seconds(), "s")
}
