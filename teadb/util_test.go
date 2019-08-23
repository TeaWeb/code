package teadb

import "testing"

func TestCurrentDB(t *testing.T) {
	db := SharedDB()
	t.Log(db)
}
