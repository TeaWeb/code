package teaproxy

import (
	"github.com/iwind/TeaGo/logs"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestNewClientPool(t *testing.T) {
	var threads = 1000
	var count = 1000

	var success = 0
	var successLocker = sync.Mutex{}

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

					req, err := http.NewRequest("GET", "http://127.0.0.1:9991/", nil)
					if err != nil {
						t.Fatal(err)
					}
					client := SharedClientPool.client("123456", "127.0.0.1:9991", 30*time.Second, 0, 0)
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
					} else {
						successLocker.Lock()
						success++
						successLocker.Unlock()
					}
				}(j)
			}
		}(i)
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)
			logs.Println("success:", success, "fails:", fails)
			if success+fails == threads*count {
				break
			}
		}
	}()

	wg.Wait()
	t.Log("finished, fails:", fails, "qps:", int(float64(threads*count)/time.Since(before).Seconds()))
}

func TestNewClientPool2(t *testing.T) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:9991/", nil)
	if err != nil {
		t.Fatal(err)
	}
	client := SharedClientPool.client("123456", "127.0.0.1:9991", 30*time.Second, 0, 0)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

func TestNewClientPool_Timeout(t *testing.T) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:9991/timeout120", nil)
	if err != nil {
		t.Fatal(err)
	}
	client := SharedClientPool.client("123456", "127.0.0.1:9991", 60*time.Second, 0, 0)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

func TestNewClientPool_Timeout2(t *testing.T) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:9991/timeout120", nil)
	if err != nil {
		t.Fatal(err)
	}
	client := SharedClientPool.client("123456", "127.0.0.1:9991", 60*time.Second, 10*time.Second, 0)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}

func TestNewClientPool_Cache(t *testing.T) {
	client1 := SharedClientPool.client("123456", "127.0.0.1:9991", 60*time.Second, 10*time.Second, 0)
	client2 := SharedClientPool.client("123456", "127.0.0.1:9991", 60*time.Second, 10*time.Second, 0)
	t.Log(client1)
	t.Log(client2)
}
