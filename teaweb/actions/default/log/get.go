package log

import (
	"fmt"
	"github.com/TeaWeb/code/teacharts"
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"math"
	"net/http"
	"time"
)

type GetAction actions.Action

// 获取日志
func (this *GetAction) Run(params struct {
	FromId       string
	Size         int64 `default:"10"`
	BodyFetching bool
}) {
	requestBodyFetching = params.BodyFetching
	requestBodyTime = time.Now()

	logger := tealogs.SharedLogger()
	accessLogs := lists.NewList(logger.ReadNewLogs(params.FromId, params.Size))
	result := accessLogs.Map(func(k int, v interface{}) interface{} {
		accessLog := v.(tealogs.AccessLog)
		return map[string]interface{}{
			"id":             accessLog.Id.Hex(),
			"requestTime":    accessLog.RequestTime,
			"request":        accessLog.Request,
			"requestURI":     accessLog.RequestURI,
			"requestMethod":  accessLog.RequestMethod,
			"remoteAddr":     accessLog.RemoteAddr,
			"remotePort":     accessLog.RemotePort,
			"userAgent":      accessLog.UserAgent,
			"host":           accessLog.Host,
			"status":         accessLog.Status,
			"statusMessage":  fmt.Sprintf("%d", accessLog.Status) + " " + http.StatusText(accessLog.Status),
			"timeISO8601":    accessLog.TimeISO8601,
			"timeLocal":      accessLog.TimeLocal,
			"requestScheme":  accessLog.Scheme,
			"proto":          accessLog.Proto,
			"contentType":    accessLog.SentContentType(),
			"bytesSent":      accessLog.BytesSent,
			"backendAddress": accessLog.BackendAddress,
			"fastcgiAddress": accessLog.FastcgiAddress,
			"extend":         accessLog.Extend,
			"referer":        accessLog.Referer,
		}
	})

	this.Data["logs"] = result.Slice

	fromTime := time.Now().Add(-24 * time.Hour)
	toTime := time.Now()

	countSuccess := logger.CountSuccessLogs(fromTime.Unix(), toTime.Unix(), "")
	countFail := logger.CountFailLogs(fromTime.Unix(), toTime.Unix())
	total := countSuccess + countFail
	this.Data["countSuccess"] = countSuccess
	this.Data["countFail"] = countFail
	this.Data["total"] = total

	// qps chart
	var qps = logger.QPS()
	qpsChart := teacharts.NewGaugeChart()
	qpsChart.Name = "实时QPS"
	qpsChart.Detail = ""
	qpsChart.Id = "qps-chart"
	qpsChart.Value = float64(qps)
	qpsChart.Unit = "Req/s"
	if qps < 100 {
		qpsChart.Max = 100
	} else if qps < 1000 {
		qpsChart.Max = 1000
	} else {
		qpsChart.Max = 10000
	}
	this.Data["qpsChart"] = qpsChart

	// bandwidth chart
	var bandWidth = logger.OutputBandWidth()
	bandwidthChart := teacharts.NewGaugeChart()
	bandwidthChart.Id = "bandwidth-chart"
	bandwidthChart.Name = "实时带宽"
	if bandWidth < 1024 {
		bandwidthChart.Detail = ""
		bandwidthChart.Value = float64(bandWidth)
		bandwidthChart.Unit = "Byte/s"
		max := math.Ceil(bandwidthChart.Value/float64(10)) * float64(10)
		if max == 0 {
			max = 100
		}
		bandwidthChart.Max = max
	} else if bandWidth < 1024*1024 {
		bandwidthChart.Detail = ""
		bandwidthChart.Value = float64(bandWidth) / 1024
		bandwidthChart.Unit = "KB/s"
		max := math.Ceil(bandwidthChart.Value/float64(10)) * float64(10)
		if max == 0 {
			max = 100
		}
		bandwidthChart.Max = max
	} else {
		bandwidthChart.Detail = ""
		bandwidthChart.Value = float64(bandWidth) / 1024 / 1024
		bandwidthChart.Unit = "MB/s"
		max := math.Ceil(bandwidthChart.Value/float64(10)) * float64(10)
		if max == 0 {
			max = 10
		}
		bandwidthChart.Max = max
	}
	this.Data["bandwidthChart"] = bandwidthChart

	// request chart，使用缓存
	cacheChart, found := this.Cache().Get("requestChart")
	if found {
		this.Data["requestChart"] = cacheChart
	} else {
		requestChart := teacharts.NewLineChart()
		requestChart.Name = "近一个小时请求数趋势"
		requestChart.Id = "request-line-chart"

		// 最近60分钟请求
		values := []interface{}{}
		labels := []string{}
		fromTime = time.Now().Add(-1 * time.Hour).Add(- time.Duration(time.Now().Second()) * time.Second)
		count := 0
		for {
			countRequest := tealogs.SharedLogger().CountSuccessLogs(fromTime.Unix(), fromTime.Unix()+60, "")
			values = append(values, countRequest)
			labels = append(labels, "")
			fromTime = fromTime.Add(1 * time.Minute)

			count ++

			if count >= 60 {
				break
			}
		}
		requestChart.AddLine(&teacharts.Line{
			Values: values,
			Color:  teacharts.ColorBlue,
			Filled: true,
			Name:   "",
		})
		requestChart.Labels = labels
		this.Data["requestChart"] = requestChart

		this.Cache().Set("requestChart", requestChart, 5*time.Second)
	}

	this.Success()
}
