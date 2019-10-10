package certutils

import (
	"github.com/TeaWeb/code/teatesting"
	"testing"
)

func TestRenewACMECerts(t *testing.T) {
	if teatesting.IsGlobal() {
		return
	}

	RenewACMECerts()
}
