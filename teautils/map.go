package teautils

import (
	"github.com/go-yaml/yaml"
	"github.com/pquerna/ffjson/ffjson"
)

func MapToObjectYAML(from map[string]interface{}, toPtr interface{}) error {
	data, err := yaml.Marshal(from)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, toPtr)
	return err
}

func MapToObjectJSON(from map[string]interface{}, toPtr interface{}) error {
	data, err := ffjson.Marshal(from)
	if err != nil {
		return err
	}

	err = ffjson.Unmarshal(data, toPtr)
	return err
}
