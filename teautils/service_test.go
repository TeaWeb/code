package teautils

import (
	"testing"
)

func TestServiceManager_Log(t *testing.T) {
	manager := NewServiceManager("TeaWeb", "TeaWeb  Server")
	manager.Log("Hello, World")
	manager.LogError("Hello, World")
}
