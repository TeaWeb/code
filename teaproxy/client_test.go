package teaproxy

import (
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestNewClientPool(t *testing.T) {
	var threads = 1000
	var count = 1000
	var fails = 0
	var failLocker = sync.Mutex{}
	var before = time.Now()
	wg := sync.WaitGroup{}
	wg.Add(threads * count)
	for i := 0; i < threads; i ++ {
		go func(i int) {
			for j := 0; j < count; j ++ {
				func(j int) {
					defer wg.Done()

					req, err := http.NewRequest("GET", "http://127.0.0.1:8800/@oldapi/ping", nil)
					if err != nil {
						t.Fatal(err)
					}
					client := SharedClientPool.client("127.0.0.1:8800", 30*time.Second, 16)
					resp, err := client.Do(req)
					if err != nil {
						t.Fatal(err)
					}
					defer resp.Body.Close()

					_, err = ioutil.ReadAll(resp.Body)
					if err != nil {
						failLocker.Lock()
						fails++
						failLocker.Unlock()
					}
				}(j)
			}
		}(i)
	}

	wg.Wait()
	t.Log("finished, fails:", fails, int(float64(threads*count)/time.Since(before).Seconds()))
}
