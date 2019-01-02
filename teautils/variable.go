package teautils

import (
	"regexp"
	"strings"
	"sync"
)

// 变量信息存储类型
type VariableHolder string

var variableMapping = map[string][]interface{}{}
var variableLocker = sync.Mutex{}
var regexpNamedVariable = regexp.MustCompile("\\${[\\w.-]+}")

// 分析变量
func ParseVariables(source string, replacer func(varName string) (value string)) string {
	variableLocker.Lock()
	defer variableLocker.Unlock()
	holders, found := variableMapping[source]
	if !found {
		indexes := regexpNamedVariable.FindAllStringIndex(source, -1)
		before := 0
		for _, loc := range indexes {
			holders = append(holders, source[before:loc[0]])
			holder := source[loc[0]+2 : loc[1]-1]
			holders = append(holders, VariableHolder(holder))
			before = loc[1]
		}
		if before < len(source) {
			holders = append(holders, source[before:])
		}
		variableMapping[source] = holders
	}
	result := strings.Builder{}
	for _, h := range holders {
		_, ok := h.(VariableHolder)
		if ok {
			key := string(h.(VariableHolder))
			result.WriteString(replacer(key))
		} else {
			result.WriteString(h.(string))
		}
	}
	return result.String()
}
