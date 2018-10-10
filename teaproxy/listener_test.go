package teaproxy

import (
	"testing"
	"net/http"
	"os"
	"bufio"
)

func TestNetClient(t *testing.T) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", "http://shop.balefm.cn/files/songs/20180608/8f600c49fc24edb7.mp3", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp.ContentLength)
}

func TestNetReadResponse(t *testing.T) {
	fp, err := os.OpenFile("/Users/liuxiangchao/Documents/Projects/pp/apps/TeaWeb/src/main/tmp/music/d03b8fb5e35237b4fbac40eb84db71c0.cache", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer fp.Close()
	response, err := http.ReadResponse(bufio.NewReader(fp), nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(response.StatusCode, response.Proto, response.Body)

}
