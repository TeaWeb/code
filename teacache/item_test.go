package teacache

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestItem_Encode(t *testing.T) {
	item := &Item{
		Header: []byte("Hello"),
		Body:   []byte("World"),
	}

	a := assert.NewAssertion(t).Quiet()
	a.Log(string(item.Encode()))

	newItem := &Item{}
	newItem.Decode(item.Encode())
	a.Equals(string(newItem.Header), "Hello")
	a.Equals(string(newItem.Body), "World")
}
