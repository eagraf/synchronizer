package main

import (
	"net"
	"testing"
)

// TestAddRemove tests adding and removing workers on the workerset
func TestAddRemove(t *testing.T) {
	wm := GetWorkerManager()
	uuid := wm.AddWorker(net.IPv4(128, 0, 0, 1))
	if len(wm.Workers) != 1 {
		t.Error("Worker was not successfully added")
	}
	wm.RemoveWorker(uuid)
	if len(wm.Workers) != 0 {
		t.Error("Worker was not successfully removed")
	}
}
