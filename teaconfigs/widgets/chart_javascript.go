package widgets

import (
	"github.com/iwind/TeaGo/utils/string"
	"regexp"
	"strings"
)

// Javascript
type JavascriptChart struct {
	Code string `yaml:"code" json:"code"`
}

func (this *JavascriptChart) AsJavascript(options map[string]interface{}) (code string, err error) {
	code = this.Code

	code = regexp.MustCompile("(\\w+)\\.render\\(\\)").ReplaceAllStringFunc(code, func(s string) string {
		index := strings.Index(s, ".")
		varName := s[0:index]
		return varName + ".options = " + stringutil.JSONEncode(options) + ";\n" + s
	})

	return code, nil
}
