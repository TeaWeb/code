package tealogs

import (
	"github.com/TeaWeb/code/tealogs/accesslogs"
	"github.com/TeaWeb/code/teatesting"
	"testing"
	"time"
)

func TestESStorage_Write(t *testing.T) {
	if !teatesting.RequireElasticSearch() {
		return
	}

	storage := &ESStorage{
		Storage: Storage{
		},
		Endpoint:    "http://127.0.0.1:9200",
		Index:       "logs",
		MappingType: "accessLogs",
	}
	err := storage.Start()
	if err != nil {
		t.Fatal(err)
	}

	{
		storage.Format = StorageFormatJSON
		storage.Template = `${timeLocal} "${requestMethod} ${requestPath}"`
		err = storage.Write([]*accesslogs.AccessLog{
			{
				RequestMethod: "POST",
				RequestPath:   "/1",
				TimeLocal:     time.Now().Format("2/Jan/2006:15:04:05 -0700"),
				Header: map[string][]string{
					"Content-Type": {"text/html"},
				},
			},
			{
				RequestMethod: "GET",
				RequestPath:   "/2",
				TimeLocal:     time.Now().Format("2/Jan/2006:15:04:05 -0700"),
				Header: map[string][]string{
					"Content-Type": {"text/css"},
				},
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	err = storage.Close()
	if err != nil {
		t.Fatal(err)
	}
}
