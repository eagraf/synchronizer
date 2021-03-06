package aggregator

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/eagraf/synchronizer/messenger"
	"github.com/eagraf/synchronizer/service"
	"google.golang.org/grpc"
)

// Selector global variable
var globalAggregator *Aggregator
var apiURL string
var rpcURL string

// Setup common server for all tests to use
func TestMain(m *testing.M) {
	var _ service.AggregatorServer = (*RPCService)(nil)

	sp := service.NewServicePool(2400, service.DefaultTopology)
	a, err := NewAggregator(sp)
	if err != nil {
		// Server failed to start
		os.Exit(-1)
	}
	globalAggregator = a
	apiURL = "http://localhost:" + strconv.Itoa(globalAggregator.service.APIPort)
	rpcURL = "localhost:" + strconv.Itoa(globalAggregator.service.RPCPort)
	os.Exit(m.Run())
}

func TestRPCServerImplementation(t *testing.T) {
	var _ service.AggregatorServer = (*RPCService)(nil)

	rs := RPCService{}
	st := reflect.TypeOf(rs)
	if !st.Implements(reflect.TypeOf((*service.AggregatorServer)(nil)).Elem()) {
		t.Error("Interface fails to implement SelectorServer")
	}
}

func TestNewAggregator(t *testing.T) {
	// Test that external API started
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		t.Error("Error constructing test API request: " + err.Error())
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error("Request failed: " + err.Error())
	}

	// Test that RPC started
	conn, err := grpc.Dial(rpcURL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := service.NewAggregatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.ReceiveSchedule(ctx, &service.AggregatorReceiveScheduleRequest{})
	if err != nil {
		t.Errorf("could not greet: %v", err)
	}
}

func TestWorkerRegistration(t *testing.T) {
	tc, err := messenger.NewTestClient(apiURL+"/websocket/", "worker-registration-id")
	if err != nil {
		t.Error("Failed to establish websocket connection: " + err.Error())
	}

	// Receive connection response
	m2, err := tc.Receive()
	if err != nil {
		t.Error("Error recieved instead of registration response: " + err.Error())
	}
	if m2.GetMetadata().MessageType != MessageInitiateDataTransfer {
		t.Error("Incorrect message type")
	}
	// Test message contents
	id, ok := m2.GetHeader("session_id")
	if ok == false {
		t.Error("No session_id header")
	}
	if len(id.(string)) != 36 {
		t.Error("id not valid UUID")
	}
	// Test workers map
	if len(globalAggregator.workers) != 1 {
		t.Error("Incorrect number of workers")
	}

	tc.Close()
}
