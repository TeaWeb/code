package teautils

import (
	"github.com/iwind/TeaGo/types"
	"math"
	"reflect"
)

// 去除导致不能转换特殊内容的问题
func ConvertJSONObjectSafely(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Map:
		v := reflect.ValueOf(obj)
		result := map[string]interface{}{}
		for _, k := range v.MapKeys() {
			k1 := k.Interface()
			v1 := v.MapIndex(k)

			// NaN
			if v1.Interface() != nil && reflect.TypeOf(v1.Interface()).Kind() == reflect.Float64 {
				if math.IsNaN(v1.Interface().(float64)) {
					result[types.String(k1)] = float64(0)
					continue
				}
			}

			// interface{} key => string key
			result[types.String(k1)] = ConvertJSONObjectSafely(v1.Interface())
		}
		return result
	case reflect.Slice:
		v := reflect.ValueOf(obj)
		result := []interface{}{}
		count := v.Len()
		for i := 0; i < count; i ++ {
			v1 := v.Index(i)
			result = append(result, ConvertJSONObjectSafely(v1.Interface()))
		}
		return result
	default:
		return obj
	}
	return obj
}
