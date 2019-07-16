package tealogs

import (
	"errors"
	"fmt"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/logs"
	"github.com/pquerna/ffjson/ffjson"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ElasticSearch存储策略
type ESStorage struct {
	Storage `yaml:", inline"`

	Endpoint    string `yaml:"endpoint" json:"endpoint"`
	Index       string `yaml:"index" json:"index"`
	MappingType string `yaml:"mappingType" json:"mappingType"`
}

// 开启
func (this *ESStorage) Start() error {
	if len(this.Endpoint) == 0 {
		return errors.New("'endpoint' should not be nil")
	}
	if !regexp.MustCompile(`(?i)^(http|https)://`).MatchString(this.Endpoint) {
		this.Endpoint = "http://" + this.Endpoint
	}
	if len(this.Index) == 0 {
		return errors.New("'index' should not be nil")
	}
	if len(this.MappingType) == 0 {
		return errors.New("'mappingType' should not be nil")
	}
	return nil
}

// 写入日志
func (this *ESStorage) Write(accessLogs []*AccessLog) error {
	if len(accessLogs) == 0 {
		return nil
	}

	bulk := &strings.Builder{}
	id := time.Now().UnixNano()
	indexName := this.FormatVariables(this.Index)
	typeName := this.FormatVariables(this.MappingType)
	for _, accessLog := range accessLogs {
		id++
		opData, err := ffjson.Marshal(map[string]interface{}{
			"index": map[string]interface{}{
				"_index": indexName,
				"_type":  typeName,
				"_id":    fmt.Sprintf("%d", id),
			},
		})
		if err != nil {
			logs.Error(err)
			continue
		}

		data, err := this.FormatAccessLogBytes(accessLog)
		if err != nil {
			logs.Error(err)
			continue
		}

		if this.Format != StorageFormatJSON {
			m := map[string]interface{}{
				"log": teautils.BytesToString(data),
			}
			mData, err := ffjson.Marshal(m)
			if err != nil {
				logs.Error(err)
				continue
			}

			bulk.Write(opData)
			bulk.WriteString("\n")
			bulk.Write(mData)
			bulk.WriteString("\n")
		} else {
			bulk.Write(opData)
			bulk.WriteString("\n")
			bulk.Write(data)
			bulk.WriteString("\n")
		}
	}

	if bulk.Len() == 0 {
		return nil
	}

	req, err := http.NewRequest(http.MethodPost, this.Endpoint+"/_bulk", strings.NewReader(bulk.String()))
	if err != nil {
		return err
	}
	client := teautils.SharedHttpClient(10 * time.Second)
	defer req.Body.Close()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyData, _ := ioutil.ReadAll(resp.Body)
		return errors.New("ElasticSearch response status code: " + fmt.Sprintf("%d", resp.StatusCode) + " content: " + string(bodyData))
	}

	return nil
}

// 关闭
func (this *ESStorage) Close() error {
	return nil
}
