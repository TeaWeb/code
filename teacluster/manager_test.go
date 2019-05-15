package teacluster

import "testing"

func TestManager_Start(t *testing.T) {
	t.Fatal(ClusterManager.Start())
}

func TestManager_PullItems(t *testing.T) {
	ClusterManager.PullItems()
}
