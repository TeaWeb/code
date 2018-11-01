package teaconfigs

import "testing"

func TestAPIMatch(t *testing.T) {
	api := NewAPI()
	api.Path = "/user/:id"
	err := api.Validate()
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(api.Match("/hello"))
	t.Log(api.Match("/user"))
	t.Log(api.Match("/user/"))
	t.Log(api.Match("/user/123"))

	api.Path = "/user/:id/:name"
	api.Validate()
	t.Log(api.Match("/user/123/liu"))
}
