package teastats

import (
	"testing"
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/TeaWeb/code/tealogs"
	"time"
)

func TestDailyUVStatAccessLogExist(t *testing.T) {
	// 是否已存在
	result := findCollection("accessLogs", nil).FindOne(context.Background(), bson.NewDocument(bson.EC.String("remoteAddr", "127.0.0.1")), findopt.Projection(map[string]int{
		"id": 1,
	}))

	existAccessLog := map[string]interface{}{}
	if result.Decode(existAccessLog) == mongo.ErrNoDocuments {
		t.Log("not exist")
	} else {
		t.Log("exist", existAccessLog)
	}
}

func TestDailyUVStat_Parse(t *testing.T) {
	accessLog := &tealogs.AccessLog{
		RemoteAddr: "127.0.0.1",
		SentHeader: map[string][]string{
			"Content-Type": {"text/html"},
		},
	}

	uv := &DailyUVStat{}
	uv.Process(accessLog)

	time.Sleep(100 * time.Millisecond)
}
