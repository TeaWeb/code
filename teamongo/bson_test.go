package teamongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestBSONDecode(t *testing.T) {
	{
		doc := primitive.D{}
		doc = append(doc, primitive.E{
			Key:   "abc",
			Value: "123",
		})
		t.Log(BSONDecode(doc))
	}

	t.Log(BSONDecode(1))
	t.Log(BSONDecode("abc"))
}
