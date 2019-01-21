package agents

import (
	"encoding/json"
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"strings"
)

// 数据格式
type SourceDataFormat = uint8

const (
	SourceDataFormatSingeLine    = 1 // 单行
	SourceDataFormatMultipleLine = 2 // 多行
	SourceDataFormatJSON         = 3 // JSON
	SourceDataFormatYAML         = 4 // YAML
)

// 所有的数据格式
func AllSourceDataFormats() []maps.Map {
	return []maps.Map{
		{
			"name": "单行数据",
			"code": SourceDataFormatSingeLine,
		},
		{
			"name": "多行数据",
			"code": SourceDataFormatMultipleLine,
		},
		{
			"name": "JSON数据",
			"code": SourceDataFormatJSON,
		},
		{
			"name": "YAML数据",
			"code": SourceDataFormatYAML,
		},
	}
}

// 取得单个数据格式
func FindSourceDataFormat(dataFormat SourceDataFormat) maps.Map {
	for _, m := range AllSourceDataFormats() {
		if types.Uint8(m["code"]) == dataFormat {
			return m
		}
	}
	return nil
}

// 解码数据
func DecodeSource(data []byte, format SourceDataFormat) (value interface{}, err error) {
	switch format {
	case SourceDataFormatSingeLine:
		return strings.TrimSpace(string(data)), nil
	case SourceDataFormatMultipleLine:
		s := strings.TrimSpace(string(data))
		if len(s) == 0 {
			return []string{}, nil
		}
		lines := strings.Split(s, "\n")
		for index, line := range lines {
			line = strings.TrimSpace(line)
			lines[index] = line
		}
		return lines, nil
	case SourceDataFormatJSON:
		v := map[string]interface{}{}
		err := json.Unmarshal(data, &v)
		if err != nil {
			return nil, err
		}
		return v, nil
	case SourceDataFormatYAML:
		v := map[string]interface{}{}
		err := yaml.Unmarshal(data, &v)
		if err != nil {
			return nil, err
		}
		return v, nil
	}

	return nil, errors.New("data format should be specified")
}
