package agents

import (
	"github.com/iwind/TeaGo/assert"
	"github.com/iwind/TeaGo/files"
	"testing"
)

func TestScriptSource_Execute(t *testing.T) {
	a := assert.NewAssertion(t)

	tmpFile := files.NewFile("/tmp/test.sh")
	tmpFile.WriteString(`#!/usr/bin/env bash

echo ${WORLD}, ${NAME}
echo 10
`)
	tmpFile.Chmod(0777)

	defer tmpFile.Delete()

	source := NewScriptSource()
	source.Path = "/tmp/test.sh"
	source.AddEnv("NAME", "ZHANG SAN")
	source.DataFormat = SourceDataFormatSingeLine
	a.IsNil(source.Validate())

	data, err := source.Execute(map[string]string{
		"WORLD": "HELLO",
	})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(data)
	}
}
