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

	c, _ := coordinator.NewCoordinator(sp)
	time.Sleep(time.Second)
	if len(c.Workers()) != 3 {
		t.Errorf("Incorrect number of workers %d", len(c.Workers()))
	}

}
