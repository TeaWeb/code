package teamongo

import (
	"context"
	"fmt"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
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

func TestCollectionStat(t *testing.T) {
	db := SharedClient().Database(DatabaseName)
	cursor, err := db.ListCollections(context.Background(), maps.Map{})
	if err != nil {
		t.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		m := maps.Map{}
		err := cursor.Decode(&m)
		if err != nil {
			t.Fatal(err)
		}
		name := m["name"].(string)
		t.Logf("%#v", name)

		result := db.RunCommand(context.Background(), bsonx.Doc{{"collStats", bsonx.String(name)}, {"verbose", bsonx.Boolean(false)}})
		if result.Err() != nil {
			t.Fatal(result.Err())
		}

		m1 := maps.Map{}
		err = result.Decode(&m1)
		if err != nil {
			t.Fatal(err)
		}
		logs.PrintAsJSON(maps.Map{
			"count": m1.GetInt("count"),
			"size":  fmt.Sprintf("%.2fM", float64(m1.GetInt("size"))/1024/1024),
			"ok":    m1.GetInt("ok"),
		}, t)
	}
}
