package selector

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/eagraf/synchronizer/messenger"
)

// Selector global variable
var globalSelector *Selector

// Setup common server for all tests to use
func TestMain(m *testing.M) {
	s, err := newSelector(3000, 3001)
	if err != nil {
		// Server failed to start
		os.Exit(-1)
	}
	globalSelector = s
	os.Exit(m.Run())
}

func TestNewSelector(t *testing.T) {
	// Test that external API started
	req, err := http.NewRequest("GET", "http://localhost:3000", nil)
	if err != nil {
		t.Error("Error constructing test API request: " + err.Error())
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error("Request failed: " + err.Error())
	}

	// Test that RPC started
	req, err = http.NewRequest("GET", "http://localhost:3001", nil)
	if err != nil {
		t.Error("Error constructing test API request: " + err.Error())
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error("Request failed: " + err.Error())
	}
}

func TestWorkerRegistration(t *testing.T) {
	tc, err := messenger.NewTestClient("http://localhost:3000/websocket/", "test-client-id")
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
}

func TestHealthCheck(t *testing.T) {
	tc, _ := messenger.NewTestClient("http://localhost:3000/websocket/", "test-client-id")

	// Receive connection response
	_, err := tc.Receive()
	if err != nil {
		t.Error("Error recieved instead of registration response: " + err.Error())
	}

	// Send healthcheck
	globalSelector.sendHealthCheck("test-client-id")
	hc, err := tc.Receive()
	if err != nil {
		t.Error("Error recieved instead of registration response: " + err.Error())
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
	if globalSelector.getWorker("test-client-id").Healthy == false {
		t.Error("Worker is not healthy")
	}
}

func TestHealthCheckTimeout(t *testing.T) {
	tc, _ := messenger.NewTestClient("http://localhost:3000/websocket/", "test-client-id")

	// Receive connection response
	_, err := tc.Receive()
	if err != nil {
		t.Error("Error recieved instead of registration response: " + err.Error())
	}

	// Send healthcheck
	globalSelector.sendHealthCheck("test-client-id")
	hc, err := tc.Receive()
	if err != nil {
		t.Error("Error recieved instead of registration response: " + err.Error())
	}
	if hc.GetMetadata().MessageType != MessageHealthCheck {
		t.Error("Incorrect message type")
	}

	time.Sleep(5 * time.Second)

	// Check worker status in selector
	if globalSelector.getWorker("test-client-id").Healthy == true {
		t.Error("Worker is healthy")
	}
}

func TestWorkerDisconnect(t *testing.T) {
	tc, _ := messenger.NewTestClient("http://localhost:3000/websocket/", "test-client-id")

	// Receive connection response
	_, err := tc.Receive()
	if err != nil {
		t.Error("Error recieved instead of registration response: " + err.Error())
	}

	err = tc.Close()
	if err != nil {
		t.Error("Error closing client: " + err.Error())
	}

	// Check worker status in selector
	if globalSelector.getWorker("test-client-id").Disconnected == false {
		t.Error("Worker is not disconnected")
	}

}

func TestRPCGetWorkers(t *testing.T) {
	// Connect three clients
	messenger.NewTestClient("http://localhost:3000/websocket/", "client1")
	messenger.NewTestClient("http://localhost:3000/websocket/", "client2")
	messenger.NewTestClient("http://localhost:3000/websocket/", "client3")

	req := WorkerRequest{}
	res := new(WorkerResponse)
	err := globalSelector.GetWorkers(req, res)
	if err != nil {
		t.Error("RPC failed: " + err.Error())
	}

	if len(res.workers) < 3 {
		t.Error("Incorrect number of workers")
	}

}

// TODO testing handoff will be difficult because it requires interactions with other services

func TestRPCHandoff(t *testing.T) {

}
