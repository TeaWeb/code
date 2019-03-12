package teautils

import (
	"github.com/go-yaml/yaml"
	"github.com/pquerna/ffjson/ffjson"
)

// 通过YAML把map转换为object
func MapToObjectYAML(fromMap map[string]interface{}, toPtr interface{}) error {
	data, err := yaml.Marshal(fromMap)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, toPtr)
	return err
}

// 通过JSON把map转换为object
func MapToObjectJSON(fromMap map[string]interface{}, toPtr interface{}) error {
	data, err := ffjson.Marshal(fromMap)
	if err != nil {
		return err
	}

	err = ffjson.Unmarshal(data, toPtr)
	return err
}

// 通过JSON把object转换为map
func ObjectToMapJSON(fromPtr interface{}, toMap *map[string]interface{}) error {
	data, err := ffjson.Marshal(fromPtr)
	if err != nil {
		return err
	}

	err = ffjson.Unmarshal(data, toMap)
	return err
}
