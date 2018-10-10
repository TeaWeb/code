package teamongo

import (
	"testing"
	"context"
)

func TestSharedClient(t *testing.T) {
	client := SharedClient()
	t.Log(client.Database("teadb").Collection("accessLog").Find(context.Background(), nil))
}
