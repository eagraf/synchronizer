package service

import (
	context "context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

// TODO take a look at how testing services is done as a whole. V jank rn.
// Mocking?

// TestServer implementation
type TestServerImpl struct{}

func (tsi TestServerImpl) TestRPC(ctx context.Context, in *Ping) (*Pong, error) {
	res := &Pong{
		Message: in.Message,
	}
	return res, nil
}

// Test service external API handler
func testEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

var sp *ServicePool

func createTestServerImpl(sp *ServicePool) (*Service, error) {
	rpc := TestServerImpl{}
	apiRouter := chi.NewRouter()
	apiRouter.Route("/test", func(r chi.Router) {
		apiRouter.Get("/", testEndpoint)
	})
	ts, err := sp.StartService("Test", rpc, apiRouter)
	return ts, err
}

func TestMain(m *testing.M) {
	topology := make(map[string]map[string]bool)
	topology["Test"] = map[string]bool{"Test": true}
	sp = NewServicePool(2000, topology)

	// Create new TestService
	_, err := createTestServerImpl(sp)
	if err != nil {
		// Server failed to start
		os.Exit(-1)
	}
	os.Exit(m.Run())
}

func TestCreateService(t *testing.T) {
	// Test the external API
	req, err := http.NewRequest("GET", "http://localhost:2000/", nil)
	if err != nil {
		t.Error("Error constructing test API request: " + err.Error())
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error("Request failed: " + err.Error())
	}

	// Test the RPC server
	conn, err := grpc.Dial("localhost:2001", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewTestClient(conn)

	// Make request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.TestRPC(ctx, &Ping{Message: "Hello"})
	if err != nil {
		t.Errorf("could not greet: %v", err)
	}
	if r.Message != "Hello" {
		t.Error("Incorrect response message")
	}
}

func TestConnect(t *testing.T) {

	s1, err := createTestServerImpl(sp)
	if err != nil {
		// Server failed to start
		t.Error("Failed to start first service")
	}
	s2, err := createTestServerImpl(sp)
	if err != nil {
		// Server failed to start
		t.Error("Failed to start second service")
	}

	err = connect(s1, s2)
	if err != nil {
		t.Error("Connect should not cause an error: " + err.Error())
	}
	if len(s1.peers["Test"]) != 1 {
		t.Error("Incorrect number of peers for service 1")
	}
	if len(s2.peers["Test"]) != 0 {
		t.Error("Incorrect number of peers for service 2")
	}
	err = connect(s2, s1)
	if err != nil {
		t.Error("Connect should not cause an error")
	}
	if len(s1.peers["Test"]) != 1 {
		t.Error("Incorrect number of peers for service 1")
	}
	if len(s2.peers["Test"]) != 1 {
		t.Error("Incorrect number of peers for service 2")
	}
	// Test sending message over connection
	// Use gRPC invoke
	cc := s1.peers["Test"][s2.ID].ClientConn
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply := Pong{}
	err = cc.Invoke(ctx, "/testservice.Test/TestRPC", &Ping{Message: "Hello"}, &reply)
	if err != nil {
		t.Errorf("could not greet: %v", err)
	}
	if reply.Message != "Hello" {
		t.Error("Incorrect response message")
	}
}

func TestConnectServices(t *testing.T) {
	// Add another service
	s1, err := createTestServerImpl(sp)
	if err != nil {
		t.Error("Error creating service 1")

	}
	s2, err := createTestServerImpl(sp)
	if err != nil {
		t.Error("Error creating service 2")
	}

	beforeCount := len(sp.Pool["Test"])
	// Connect to all existing services (including s2 and those started by other tests)
	sp.ConnectService(s1)

	// Based off the topology, s2 should now have a connection to s1
	if len(s2.peers["Test"]) != 1 {
		t.Error("Failed to connect second service to first one")
	}

	if len(s1.peers["Test"]) != beforeCount {
		t.Error("Service 1 did not connect to all peers")
	}
}

// Test telemetry

func TestTelemetryLogging(t *testing.T) {
	s1, err := createTestServerImpl(sp)
	if err != nil {
		// Server failed to start
		t.Error("Failed to start first service")
	}
	if s1.Logger == nil {
		t.Error("Service initiator failed to assign a logger")
	}
	s1.Log("NewTag", "Test Logging")
	if logs, ok := sp.logStore.tags["NewTag"]; ok != false {
		if len(logs) != 1 {
			t.Error("Incorrect length for NewTags list")
		}
	} else {
		t.Error("NewTags list never created")
	}
	if logs, ok := sp.logStore.tags[s1.ID]; ok != false {
		if len(logs) != 1 {
			t.Error("Incorrect length for service tag list")
		}
	} else {
		t.Error("service tag list never created")
	}

	s1.Log("NewTag", "Test Logging 2")
	if len(sp.logStore.tags["NewTag"]) != 2 {
		t.Error("NewTag list not long enough")
	}

}
