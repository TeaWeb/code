package teamongo

import (
	"context"
	"testing"
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
