package teatesting

import (
	"github.com/iwind/TeaGo/logs"
	"net/http"
	"time"
)

// 需要测试HTTP Server支持
func RequireHTTPServer() bool {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:9991", nil)
	if err != nil {
		logs.Error(err)
		return false
	}
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		go StartTestServer()
		time.Sleep(1 * time.Second)

		// test again
		resp, err := c.Do(req)
		if err != nil {
			return false
		}
		_ = resp.Body.Close()
		return true
	} else {
		_ = resp.Body.Close()
		return true
	}
}

// 需要TeaWeb Server支持
func RequireTeaWebServer() bool {
	// TODO
	return false
}

// 需要Nginx Status支持
func RequireNginxStatus() bool {
	// TODO
	return false
}

// 需要MySQL支持
func RequireMySQL() bool {
	// TODO
	return false
}

// 需要Postgres支持
func RequirePostgres() bool {
	// TODO
	return false
}

// 需要Docker支持
func RequireDocker() bool {
	// TODO
	return false
}

// 需要端口支持
func RequirePort(port int) bool {
	// TODO
	return false
}

// 需要DNS支持
func RequireDNS() bool {
	// TODO
	return false
}

// 需要Redis支持
func RequireRedis() bool {
	// TODO
	return false
}

// 需要ES支持
func RequireElasticSearch() bool {
	// TODO
	return false
}

// 检查Fastcgi
func RequireFascgi() bool {
	// TODO
	return false
}

// 检查数据库
func RequireDBAvailable() bool {
	// TODO
	if IsGlobal() {
		return false
	}
	return true
}
