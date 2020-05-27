package selector

import (
	"context"
	"fmt"
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
var globalSelector *Selector
var apiURL string
var rpcURL string

// Setup common server for all tests to use
func TestMain(m *testing.M) {
	var _ service.SelectorServer = (*RPCService)(nil)

	sp := service.NewServicePool(service.DefaultTopology)
	s, err := newSelector(sp)
	if err != nil {
		// Server failed to start
		os.Exit(-1)
	}
	globalSelector = s
	apiURL = "http://localhost:" + strconv.Itoa(globalSelector.service.APIPort)
	rpcURL = "localhost:" + strconv.Itoa(globalSelector.service.RPCPort)
	os.Exit(m.Run())
}

func TestRPCServerImplementation(t *testing.T) {
	var _ service.SelectorServer = (*RPCService)(nil)

	rs := RPCService{}
	st := reflect.TypeOf(rs)
	if !st.Implements(reflect.TypeOf((*service.SelectorServer)(nil)).Elem()) {
		t.Error("Interface fails to implement SelectorServer")
	}
}

func TestNewSelector(t *testing.T) {
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
	c := service.NewSelectorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.GetWorkers(ctx, &service.WorkersRequest{})
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
	if m2.GetMetadata().MessageType != MessageRegistrationResponse {
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
	if len(globalSelector.workers) != 1 {
		t.Error("Incorrect number of workers")
	}

	tc.Close()
}

func TestHealthCheck(t *testing.T) {
	tc, _ := messenger.NewTestClient(apiURL+"/websocket/", "health-check-id")

	// Receive connection response
	_, err := tc.Receive()
	if err != nil {
		t.Error("Error recieved instead of registration response: " + err.Error())
	}

	// Send healthcheck
	globalSelector.sendHealthCheck("health-check-id")
	hc, err := tc.Receive()
	if err != nil {
		t.Error("Error recieved instead of health check: " + err.Error())
	}
	if hc.GetMetadata().MessageType != MessageHealthCheck {
		t.Error("Incorrect message type")
	}

	// Send response to healthcheck
	mb := new(messenger.MessageBuilder)
	m, _ := mb.NewMessage("worker_health_check_response", hc.GetMetadata().Request).Done()
	err = tc.Send(m)
	if err != nil {
		t.Error("Error sending health check response: " + err.Error())
	}

	// Check worker status in selector
	worker, _ := globalSelector.getWorker("health-check-id")
	if worker.Healthy == false {
		t.Error("Worker is not healthy")
	}

	tc.Close()
}

func TestHealthCheckTimeout(t *testing.T) {
	tc, _ := messenger.NewTestClient(apiURL+"/websocket/", "health-check-timeout-id")

	// Receive connection response
	_, err := tc.Receive()
	if err != nil {
		t.Error("Error recieved instead of registration response: " + err.Error())
	}

	// Send healthcheck
	globalSelector.sendHealthCheck("health-check-timeout-id")
	hc, err := tc.Receive()
	if err != nil {
		t.Error("Error recieved instead of registration response: " + err.Error())
	}
	if hc.GetMetadata().MessageType != MessageHealthCheck {
		t.Error("Incorrect message type")
	}

	time.Sleep(6 * time.Second)

	// Check worker status in selector
	worker, _ := globalSelector.getWorker("health-check-timeout-id")
	if worker.Healthy == true {
		t.Error("Worker is healthy")
	}

	tc.Close()
}

func TestWorkerDisconnect(t *testing.T) {
	tc, _ := messenger.NewTestClient(apiURL+"/websocket/", "worker-disconnect-id")

	// Receive connection response
	_, err := tc.Receive()
	if err != nil {
		t.Error("Error recieved instead of registration response: " + err.Error())
	}
	err = tc.Close()
	time.Sleep(time.Second) // Given some delay
	if err != nil {
		t.Error("Error closing client: " + err.Error())
	}

	// Check worker status in selector
	worker, _ := globalSelector.getWorker("worker-disconnect-id")
	if worker.Disconnected == false {
		t.Error("Worker is not disconnected")
	}
}

func TestRPCGetWorkers(t *testing.T) {
	// Connect three clients
	tc1, _ := messenger.NewTestClient(apiURL+"/websocket/", "client1")
	tc2, _ := messenger.NewTestClient(apiURL+"/websocket/", "client2")
	tc3, _ := messenger.NewTestClient(apiURL+"/websocket/", "client3")

	// Test that RPC started
	conn, err := grpc.Dial(rpcURL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := service.NewSelectorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	workers, err := c.GetWorkers(ctx, &service.WorkersRequest{})
	if err != nil {
		t.Errorf("could not greet: %v", err)
	}
	fmt.Println(workers)

	tc1.Close()
	tc2.Close()
	tc3.Close()
}

// TODO testing handoff will be difficult because it requires interactions with other services

func TestRPCHandoff(t *testing.T) {

}
