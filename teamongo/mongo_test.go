package teamongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestSharedClient(t *testing.T) {
	client := SharedClient()
	t.Log(client.Database("teadb").Collection("accessLog").Find(context.Background(), map[string]interface{}{}))
}

func TestUnmarshalJSON(t *testing.T) {
	data := `{
		"$group": {
			"_id": null,
			"total": {
				"$sum": "$count"
			}
		},
		"$match": {
			"serverId": "123",
			"day": {
				"$in": [ "20181010", "20181011" ]
			}
		}
	}`
	t.Log(data)

	value, err := BSONArrayBytes([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(value)
}

func TestFindCollection(t *testing.T) {
	coll := FindCollection("logs.20190302")
	//opts := options.FindOne().SetHint(map[string]interface{}{
	//	"serverId": 1,
	//})
	{
		before := time.Now()
		cursor, err := coll.Find(context.Background(), map[string]interface{}{
			"serverId": "VEQ6mBKq7w7lFUzj",
		}, options.Find().SetLimit(1).SetHint(map[string]interface{}{
			"serverId": 1,
		}))
		if err != nil {
			t.Fatal(err)
		}
		defer cursor.Close(context.Background())
		if cursor.Next(context.Background()) {
			t.Log("has")
		} else {
			t.Log("not has")
		}
		t.Log(time.Since(before).Seconds(), "s")
	}
}
