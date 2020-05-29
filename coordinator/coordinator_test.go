package coordinator

import (
	"fmt"
	"testing"
	"time"

	"github.com/eagraf/synchronizer/messenger"
	"github.com/eagraf/synchronizer/selector"
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
	time.Sleep(2 * c.interval) // A little more than one interval
	// Two intervals should have passed
	if c.round != 2 {
		t.Errorf("Wrong number of rounds: %d", c.round)
	}
}

func TestSchedule(t *testing.T) {

}

// Not really necessary anymore... eh
func TestMultiCast(t *testing.T) {
	sp := service.NewServicePool(4010, service.DefaultTopology)
	selector.NewSelector(sp)
	selector.NewSelector(sp)
	selector.NewSelector(sp)
	c, _ := NewCoordinator(sp)

	replys := make([]interface{}, 4)
	for i := range replys {
		replys[i] = &service.WorkersResponse{}
	}
	peers, _ := c.service.AllPeersOfType("Selector")
	fmt.Println(len(peers))
	req := &service.WorkersRequest{
		All:    true,
		Number: 0,
	}
	responses, errs := c.service.MultiCast(peers, service.SelectorGetWorkers, req, replys)
	if len(errs) != 0 {
		t.Error(errs[0].Error())
	}
	if len(responses) != 4 {
		t.Error("Incorrect number of responses")
	}
	fmt.Printf("Response %v\n", *(responses[0].(*service.WorkersResponse)))
}

func TestGetWorkersFromSelectors(t *testing.T) {
	sp := service.NewServicePool(5000, service.DefaultTopology)
	selector.NewSelector(sp)
	selector.NewSelector(sp)
	selector.NewSelector(sp)

	messenger.NewTestClient("http://localhost:5000/websocket/", "client1")
	messenger.NewTestClient("http://localhost:5002/websocket/", "client2")
	messenger.NewTestClient("http://localhost:5004/websocket/", "client3")

	c, _ := NewCoordinator(sp)
	time.Sleep(time.Second)
	if len(c.workers) != 3 {
		t.Errorf("Incorrect number of workers %d", len(c.workers))
	}

}
