package teamongo

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"testing"
)

func TestSharedClient(t *testing.T) {
	client := SharedClient()
	t.Log(client.Database("teadb").Collection("accessLog").Find(context.Background(), nil))
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

	arr := bson.NewDocument()
	err := bson.UnmarshalExtJSON([]byte(data), true, &arr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(arr)
}
