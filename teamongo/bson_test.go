package teamongo

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"testing"
)

func TestBSONArray(t *testing.T) {
	data := `[  "1", "2", "3", 1, { "name": "hello" } ]`
	arr, err := BSONArrayBytes([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(arr)
}

func TestBSONObject(t *testing.T) {
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

	arr, err := BSONObjectBytes([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(arr)
}

func TestBSONDecode(t *testing.T) {
	{
		doc := bson.NewDocument()
		doc.Append(bson.EC.Int32("a", 1))
		doc.Append(bson.EC.Int32("b", 2))
		t.Log(BSONDecode(doc))
	}

	{
		doc, err := BSONArray([]interface{}{1, 2, 3, 4, 5})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(BSONDecode(doc))
	}

	{
		doc := bson.NewDocument()
		doc.Append(bson.EC.Int32("a", 1))
		doc.Append(bson.EC.Int32("b", 2))

		arrayDoc, err := BSONArray([]interface{}{1, 2, 3, 4, 5})
		if err != nil {
			t.Fatal(err)
		}
		doc.Append(bson.EC.Array("c", arrayDoc))
		t.Log(BSONDecode(doc))
	}

	t.Log(BSONDecode(1))
	t.Log(BSONDecode("abc"))
}
