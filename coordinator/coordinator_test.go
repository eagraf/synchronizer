package coordinator

import (
	"testing"
	"time"

	"github.com/eagraf/synchronizer/service"
)

var globalCoordinator *Coordinator
var apiURL string

/*func TestMain(m *testing.M) {
	//	var _ service.SelectorServer = (*RPCService)(nil)

	sp := service.NewServicePool(4000, service.DefaultTopology)
	c, err := newCoordinator(sp)
	if err != nil {
		// Server failed to start
		os.Exit(-1)
	}
	globalCoordinator = c
	apiURL = "http://localhost:" + strconv.Itoa(globalCoordinator.service.APIPort)
	os.Exit(m.Run())
}*/

func TestStartCoordinator(t *testing.T) {
	sp := service.NewServicePool(4000, service.DefaultTopology)
	c, err := NewCoordinator(sp)
	if err != nil {
		t.Error(err.Error())
	}
	time.Sleep(c.interval + time.Second) // A little more than one interval
	// Two intervals should have passed
	if c.round != 2 {
		t.Errorf("Wrong number of rounds: %d", c.round)
	}
}

func TestSchedule(t *testing.T) {

}
