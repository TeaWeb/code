package teamongo

import (
	"encoding/json"
	"errors"
	"github.com/iwind/TeaGo/types"
	"github.com/mongodb/mongo-go-driver/bson"
	"reflect"
)

func BSONObject(m map[string]interface{}) (result *bson.Document, err error) {
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
			doc, err := BSONObject(item)
			if err != nil {
				return nil, err
			}
			result.Append(bson.EC.SubDocument(key, doc))
		case []interface{}:
			arr, err := BSONArray(item)
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

func BSONArray(list []interface{}) (result *bson.Array, err error) {
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
			doc, err := BSONObject(item)
			if err != nil {
				return nil, err
			}
			result.Append(bson.EC.SubDocument("", doc).Value())
		case []interface{}:
			arr, err := BSONArray(item)
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

func BSONArrayBytes(data []byte) (result *bson.Array, err error) {
	list := []interface{}{}
	err = json.Unmarshal(data, &list)
	if err != nil {
		return
	}

	return BSONArray(list)
}

func BSONObjectBytes(data []byte) (result *bson.Document, err error) {
	m := map[string]interface{}{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return
	}
	return BSONObject(m)
}

func BSONDecode(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case *bson.Document:
		m := map[string]interface{}{}
		iterator := v.Iterator()
		for iterator.Next() {
			e := iterator.Element()
			elementValue, err := BSONDecode(e.Value())
			if err != nil {
				return nil, err
			}
			m[e.Key()] = elementValue
		}
		return m, nil
	case *bson.Array:
		s := []interface{}{}
		count := v.Len()
		for i := 0; i < count; i ++ {
			e, err := v.Lookup(uint(i))
			if err != nil {
				return nil, err
			}
			elementValue, err := BSONDecode(e)
			if err != nil {
				return nil, err
			}
			s = append(s, elementValue)
		}
		return s, nil
	case *bson.Value:
		v2 := v.Interface()
		return BSONDecode(v2)
	case map[string]interface{}:
		for itemKey, itemValue := range v {
			r, err := BSONDecode(itemValue)
			if err != nil {
				return nil, err
			}
			v[itemKey] = r
		}
	}

	return value, nil
}
