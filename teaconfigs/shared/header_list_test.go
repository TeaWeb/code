package shared

import (
	"fmt"
	"github.com/TeaWeb/code/teautils"
	"testing"
	"time"
)

func TestHeaderList_FormatHeaders(t *testing.T) {

	list := &HeaderList{}

	for i := 0; i < 5; i ++ {
		list.AddHeader(&HeaderConfig{
			On:    true,
			Name:  "A" + fmt.Sprintf("%d", i),
			Value: "ABCDEFGHIJ${name}KLM${hello}NEFGHIJILKKKk",
		})
	}

	list.ValidateHeaders()

	b := time.Now()
	count := 1000000
	for j := 0; j < count; j ++ {
		list.FormatHeaders(func(source string) string {
			return teautils.ParseVariables(source, func(varName string) (value string) {
				return "abc"
			})
		})

		//logs.PrintAsJSON(newList, t)
	}
	t.Logf("%f qps", float64(count)/time.Since(b).Seconds())
}
