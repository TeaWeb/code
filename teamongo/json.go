package teamongo

import (
	"encoding/json"
	"errors"
	"github.com/iwind/TeaGo/types"
	"github.com/mongodb/mongo-go-driver/bson"
	"reflect"
)

func JSONObject(m map[string]interface{}) (result *bson.Document, err error) {
	result = bson.NewDocument()
	for key, item := range m {
		switch item := item.(type) {
		case string:
			result.Append(bson.EC.String(key, item))
		case int8, int16, int, int32, int64, uint, uint8, uint16, uint32, uint64:
			result.Append(bson.EC.Int64(key, types.Int64(item)))
		case bool:
			result.Append(bson.EC.Boolean(key, item))
		case float32, float64:
			result.Append(bson.EC.Double(key, types.Float64(item)))
		case map[string]interface{}:
			doc, err := JSONObject(item)
			if err != nil {
				return nil, err
			}
			result.Append(bson.EC.SubDocument(key, doc))
		case []interface{}:
			arr, err := JSONArray(item)
			if err != nil {
				return nil, err
			}
			result.Append(bson.EC.Array(key, arr))
		default:
			if item == nil {
				result.Append(bson.EC.Null(key))
				continue
			}
			err = errors.New("unknown json data type '" + reflect.TypeOf(item).String() + "'")
			return
		}
	}
	return
}

func JSONArray(list []interface{}) (result *bson.Array, err error) {
	result = bson.NewArray()
	for _, item := range list {
		switch item := item.(type) {
		case string:
			result.Append(bson.EC.String("", item).Value())
		case int8, int16, int, int32, int64, uint, uint8, uint16, uint32, uint64:
			result.Append(bson.EC.Int64("", types.Int64(item)).Value())
		case bool:
			result.Append(bson.EC.Boolean("", item).Value())
		case float32, float64:
			result.Append(bson.EC.Double("", types.Float64(item)).Value())
		case map[string]interface{}:
			doc, err := JSONObject(item)
			if err != nil {
				return nil, err
			}
			result.Append(bson.EC.SubDocument("", doc).Value())
		case []interface{}:
			arr, err := JSONArray(item)
			if err != nil {
				return nil, err
			}
			result.Append(bson.EC.Array("", arr).Value())
		default:
			if item == nil {
				result.Append(bson.EC.Null("").Value())
				continue
			}
			err = errors.New("unknown json data type '" + reflect.TypeOf(item).String() + "'")
			return
		}
	}

	return
}

func JSONArrayBytes(data []byte) (result *bson.Array, err error) {
	list := []interface{}{}
	err = json.Unmarshal(data, &list)
	if err != nil {
		return
	}

	return JSONArray(list)
}

func JSONObjectBytes(data []byte) (result *bson.Document, err error) {
	m := map[string]interface{}{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return
	}
	return JSONObject(m)
}
