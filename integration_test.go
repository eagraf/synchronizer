package main

import (
	"testing"
	"time"

	"github.com/eagraf/synchronizer/aggregator"
	"github.com/eagraf/synchronizer/coordinator"
	"github.com/eagraf/synchronizer/dataserver"
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
		t.Errorf("Incorrect number of worker recvs: %d", count)
	}
	if exists := service.LogExists(sp, "GetWorkersRecv", "Receiving 3 workers from 3 selectors"); !exists {
		t.Errorf("Correct receiving log not found")
	}
}

// TODO func TestCoordinatorSchedule()
func TestCoordinatorSchedule(t *testing.T) {
	sp := service.NewServicePool(5100, service.DefaultTopology)
	// Three selectors, aggregators and data servers
	selector.NewSelector(sp)
	selector.NewSelector(sp)
	selector.NewSelector(sp)

	aggregator.NewAggregator(sp)
	aggregator.NewAggregator(sp)
	aggregator.NewAggregator(sp)

	dataserver.NewDataServer(sp)
	dataserver.NewDataServer(sp)
	dataserver.NewDataServer(sp)

	// Nine workers evenly distributed across selectors
	messenger.NewTestClient("http://localhost:5100/websocket/", "client1")
	messenger.NewTestClient("http://localhost:5102/websocket/", "client2")
	messenger.NewTestClient("http://localhost:5104/websocket/", "client3")
	messenger.NewTestClient("http://localhost:5100/websocket/", "client4")
	messenger.NewTestClient("http://localhost:5102/websocket/", "client5")
	messenger.NewTestClient("http://localhost:5104/websocket/", "client6")
	messenger.NewTestClient("http://localhost:5100/websocket/", "client7")
	messenger.NewTestClient("http://localhost:5102/websocket/", "client8")
	messenger.NewTestClient("http://localhost:5104/websocket/", "client9")

	// One coordinator
	coordinator.NewCoordinator(sp)
	time.Sleep(time.Second)

	if count := service.CountTags(sp, "GetWorkersSend"); count != 3 {
		t.Errorf("Incorrect number of worker sends: %d", count)
	}
	if count := service.CountTags(sp, "GetWorkersRecv"); count != 1 {
		t.Errorf("Incorrect number of worker recvs: %d", count)
	}
	if exists := service.LogExists(sp, "GetWorkersRecv", "Receiving 9 workers from 3 selectors"); !exists {
		t.Errorf("Correct receiving log not found")
	}

	// Test scheduling receives
	if count := service.CountTags(sp, "DataServerReceiveSchedule"); count != 3 {
		t.Errorf("Incorrect number of worker sends: %d", count)
	}
}
