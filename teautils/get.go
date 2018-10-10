package teautils

import (
	"reflect"
)

func Get(object interface{}, keys []string) interface{} {
	if len(keys) == 0 {
		return object
	}

	if object == nil {
		return nil
	}

	firstKey := keys[0]
	keys = keys[1:]

	value := reflect.ValueOf(object)

	if !value.IsValid() {
		return nil
	}

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() == reflect.Struct {
		field := value.FieldByName(firstKey)
		if !field.IsValid() {
			return nil
		}

		if len(keys) == 0 {
			return field.Interface()
		}

		return Get(field.Interface(), keys)
	}

	if value.Kind() == reflect.Map {
		mapKey := reflect.ValueOf(firstKey)
		mapValue := value.MapIndex(mapKey)
		if !mapValue.IsValid() {
			return nil
		}

		if len(keys) == 0 {
			return mapValue.Interface()
		}

		return Get(mapValue.Interface(), keys)
	}

	return nil
}
