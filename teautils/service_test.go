package teautils

import (
	"github.com/TeaWeb/code/teaconst"
	"testing"
)

func TestServiceManager_Log(t *testing.T) {
	manager := NewServiceManager(teaconst.TeaProductName, teaconst.TeaProductName+" Server")
	manager.Log("Hello, World")
	manager.LogError("Hello, World")
}
