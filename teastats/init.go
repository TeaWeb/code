package teastats

import (
	"github.com/TeaWeb/code/tealogs"
)

func init() {
	tealogs.SharedLogger().AddProcessor(new(Processor))
}
