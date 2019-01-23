package teautils

// 转换JSON可编码的Map
func JSONMap(i interface{}) interface{} {
	switch x := i.(type) {
	case map[string]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k] = JSONMap(v)
		}
		return m2
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = JSONMap(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = JSONMap(v)
		}
	}
	return i
}
