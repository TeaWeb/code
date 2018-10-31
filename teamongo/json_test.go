package teamongo

import "testing"

func TestJSONArray(t *testing.T) {
	data := `[  "1", "2", "3", 1, { "name": "hello" } ]`
	arr, err := JSONArrayBytes([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(arr)
}

func TestJSONObject(t *testing.T) {
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

	arr, err := JSONObjectBytes([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(arr)
}