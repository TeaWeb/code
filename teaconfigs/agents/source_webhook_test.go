package agents

import (
	"net/http"
	"testing"
)

func TestWebHookSource_ExecuteGet(t *testing.T) {
	webHook := NewWebHookSource()
	webHook.Method = http.MethodGet
	webHook.URL = "http://127.0.0.1:9991/webhook?hell=world"
	webHook.DataFormat = SourceDataFormatSingeLine
	err := webHook.Validate()
	if err != nil {
		t.Fatal(err)
	}
	result, err := webHook.Execute(map[string]string{
		"host": "127.0.0.1",
		"port": "3306",
	})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(result)
	}
}

func TestWebHookSource_ExecutePost(t *testing.T) {
	webHook := NewWebHookSource()
	webHook.Method = http.MethodPost
	webHook.URL = "http://127.0.0.1:9991/webhook?hell=world"
	webHook.DataFormat = SourceDataFormatSingeLine
	err := webHook.Validate()
	if err != nil {
		t.Fatal(err)
	}
	result, err := webHook.Execute(map[string]string{
		"host": "127.0.0.1",
		"port": "3306",
	})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(result)
	}
}
