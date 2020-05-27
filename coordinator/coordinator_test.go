package coordinator

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/eagraf/synchronizer/service"
)

var globalCoordinator *Coordinator
var apiURL string

func TestMain(m *testing.M) {
	//	var _ service.SelectorServer = (*RPCService)(nil)

	sp := service.NewServicePool(service.DefaultTopology)
	c, err := newCoordinator(sp)
	if err != nil {
		// Server failed to start
		os.Exit(-1)
	}
	globalCoordinator = c
	apiURL = "http://localhost:" + strconv.Itoa(globalCoordinator.service.APIPort)
	os.Exit(m.Run())
}

func TestStartCoordinator(t *testing.T) {
	time.Sleep(globalCoordinator.interval + time.Second) // A little more than one interval
	// Two intervals should have passed
	if globalCoordinator.round != 2 {
		t.Error("Wrong number of rounds")
	}
}
