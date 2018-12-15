package teacache

import (
	"bufio"
	"bytes"
	"github.com/TeaWeb/code/teaconfigs/shared"
	"github.com/TeaWeb/code/teaproxy"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"io"
	"net/http"
	"sync"
)

var cachePolicyMap = map[*shared.CachePolicy]ManagerInterface{}
var cachePolicyMapLocker = sync.RWMutex{}

func ProcessBeforeRequest(req *teaproxy.Request, writer *teaproxy.ResponseWriter) bool {
	cacheConfig := req.CachePolicy()
	if cacheConfig == nil || !cacheConfig.On {
		return true
	}

	cachePolicyMapLocker.RLock()
	cache, found := cachePolicyMap[cacheConfig]
	cachePolicyMapLocker.RUnlock()
	if !found {
		cache = NewManagerFromConfig(cacheConfig)
		if cache == nil {
			return true
		}
		cachePolicyMapLocker.Lock()
		cachePolicyMap[cacheConfig] = cache
		cachePolicyMapLocker.Unlock()
	}

	// key
	if len(cacheConfig.Key) == 0 {
		return true
	}
	key := req.Format(cacheConfig.Key)
	data, err := cache.Read(key)
	if err != nil {
		if err != ErrNotFound {
			logs.Error(err)
		} else {
			req.SetCacheEnabled()
			writer.SetBodyCopying(true)
		}
		return true
	}

	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(data[8:])), nil)
	if err != nil {
		logs.Error(err)
		return true
	}
	defer resp.Body.Close()

	writer.WriteHeader(resp.StatusCode)
	for k, vs := range resp.Header {
		for _, v := range vs {
			writer.Header().Add(k, v)
		}
	}
	io.Copy(writer, resp.Body)

	return false
}

func ProcessAfterRequest(req *teaproxy.Request, writer *teaproxy.ResponseWriter) bool {
	if !req.IsCacheEnabled() {
		return true
	}

	cacheConfig := req.CachePolicy()
	if cacheConfig == nil {
		return true
	}

	//check status
	if writer.StatusCode() == http.StatusNotModified { // 如果没有修改就不会有body，会有陷阱，所以这里不加入缓存
		return true
	}
	if len(cacheConfig.Status) == 0 {
		cacheConfig.Status = []int{http.StatusOK}
	}
	if !lists.Contains(cacheConfig.Status, writer.StatusCode()) {
		return true
	}

	cachePolicyMapLocker.RLock()
	cache, found := cachePolicyMap[cacheConfig]
	cachePolicyMapLocker.RUnlock()
	if !found {
		return true
	}

	key := req.Format(cacheConfig.Key)
	headerData := writer.HeaderData()
	item := &Item{
		Header: headerData,
		Body:   writer.Body(),
	}
	if len(headerData) == 0 {
		return true
	}
	data := item.Encode()
	if cacheConfig.MaxDataSize() > 0 && float64(len(data)) > cacheConfig.MaxDataSize() {
		return true
	}
	err := cache.Write(key, data)
	if err != nil {
		logs.Error(err)
	}
	return true
}
