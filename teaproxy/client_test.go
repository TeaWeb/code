package teaproxy

import (
	"context"
	"crypto/tls"
	"github.com/iwind/TeaGo/logs"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
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
						failLocker.Lock()
						fails++
						failLocker.Unlock()
						return
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

func BenchmarkNewClientPool(b *testing.B) {
	client := SharedClientPool.client("123456", "127.0.0.1:9991", 10*time.Second, 10*time.Second, 0)

	for i := 0; i < b.N; i ++ {
		req, err := http.NewRequest("GET", "http://127.0.0.1:9991/", nil)
		if err != nil {
			b.Fatal(err)
		}
		resp, err := client.Do(req)
		if err == nil {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}
	}
}

func BenchmarkNewClientPool2(b *testing.B) {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				// 握手配置
				c, err := (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 10 * time.Minute,
					DualStack: true,
				}).DialContext(ctx, network, addr)
				return c, err
			},
			Proxy:                 nil,
			MaxIdleConns:          1024,
			MaxIdleConnsPerHost:   1024,
			IdleConnTimeout:       2 * time.Minute,
			ExpectContinueTimeout: 1 * time.Second,
			TLSHandshakeTimeout:   0, // 不限
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DisableCompression: true,
		},
	}

	for i := 0; i < b.N; i ++ {
		req, err := http.NewRequest("GET", "http://127.0.0.1:9991/", nil)
		if err != nil {
			b.Fatal(err)
		}
		resp, err := client.Do(req)
		if err == nil {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}
	}
}

func TestClient_Proxy(t *testing.T) {
	mux := http.NewServeMux()

	u, _ := url.Parse("http://127.0.0.1:9991")
	proxy := httputil.NewSingleHostReverseProxy(u)

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		proxy.ServeHTTP(writer, request)
	})
	server := &http.Server{
		Handler: mux,
		Addr:    "127.0.0.1:8890",
	}
	server.ListenAndServe()
}

func TestMultipleClient(t *testing.T) {
	count := 100
	n := 1000
	t.Log(n, "threads")

	wg := sync.WaitGroup{}
	wg.Add(n)

	dials := int32(0)
	success := int64(0)

	client1 := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				atomic.AddInt32(&dials, 1)

				return (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 10 * time.Minute,
					DualStack: true,
				}).DialContext(ctx, network, addr)
			},
			Proxy:                 nil,
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   1024,
			IdleConnTimeout:       2 * time.Minute,
			ExpectContinueTimeout: 1 * time.Second,
			TLSHandshakeTimeout:   0, // 不限
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DisableCompression: true,
		},
	}

	before := time.Now()
	for j := 0; j < n; j ++ {
		go func() {
			defer wg.Done()

			for i := 0; i < count; i ++ {
				req, err := http.NewRequest("GET", "http://127.0.0.1:9991/", nil)
				if err != nil {
					t.Fatal(err)
				}
				resp, err := client1.Do(req)
				if err == nil {
					data := []byte{}
					b := make([]byte, resp.ContentLength)
					for {
						n, err := resp.Body.Read(b)
						if n > 0 {
							data = append(data, b[:n]...)
						}
						if err != nil {
							break
						}
					}

					resp.Body.Close()
					atomic.AddInt64(&success, 1)
				}
			}
		}()
	}
	wg.Wait()

	t.Log(float64(count*n) / time.Since(before).Seconds())
	t.Log(dials, "dials", "success:", success)
}
