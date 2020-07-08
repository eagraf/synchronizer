package main

import (
	"testing"
	"time"

	"github.com/eagraf/synchronizer/coordinator"
	"github.com/eagraf/synchronizer/messenger"
	"github.com/eagraf/synchronizer/selector"
	"github.com/eagraf/synchronizer/service"
)

func TestCoordinatorGetWorkersFromSelectors(t *testing.T) {
	sp := service.NewServicePool(5000, service.DefaultTopology)
	selector.NewSelector(sp)
	selector.NewSelector(sp)
	selector.NewSelector(sp)

	messenger.NewTestClient("http://localhost:5000/websocket/", "client1")
	messenger.NewTestClient("http://localhost:5002/websocket/", "client2")
	messenger.NewTestClient("http://localhost:5004/websocket/", "client3")

	coordinator.NewCoordinator(sp)
	time.Sleep(time.Second)

	if count := service.CountTags(sp, "GetWorkersSend"); count != 3 {
		t.Errorf("Incorrect number of worker sends: %d", count)
	}
	if count := service.CountTags(sp, "GetWorkersRecv"); count != 1 {
		t.Errorf("Incorrect number of worker sends: %d", count)
	}
	if exists := service.LogExists(sp, "GetWorkersRecv", "Receiving 3 workers from 3 selectors"); !exists {
		t.Errorf("Correct receiving log not found")
	}
}

// TODO func TestCoordinatorSchedule()
