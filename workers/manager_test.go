package workers

import (
	"testing"
)

// TestAddRemove tests adding and removing workers on the workerset
// TODO this might need a mutex
func TestAddRemove(t *testing.T) {
	wm := GetWorkerManager()
	numWorkers := len(wm.Workers)
	uuid := wm.AddWorker("cloud")
	if len(wm.Workers) != numWorkers+1 {
		t.Error("Worker was not successfully added")
	}
	wm.RemoveWorker(uuid)
	if len(wm.Workers) != numWorkers {
		t.Error("Worker was not successfully removed")
	}
}
