package teacluster

import "testing"

func TestManager_Start(t *testing.T) {
	t.Fatal(SharedManager.Start())
}

func TestManager_PullItems(t *testing.T) {
	SharedManager.PullItems()
}
