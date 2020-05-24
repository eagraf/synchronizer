package service

import (
	context "context"
	fmt "fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

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

func TestMain(m *testing.M) {
	sp := NewServicePool()

	// Create new TestService
	rpc := TestServerImpl{}
	apiRouter := chi.NewRouter()
	apiRouter.Route("/test", func(r chi.Router) {
		apiRouter.Get("/", testEndpoint)
	})
	_, err := sp.StartService("Test", rpc, apiRouter)
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
	fmt.Println("hello")
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error("Request failed: " + err.Error())
	}
	fmt.Println("ello")

	// Test the RPC server
	conn, err := grpc.Dial("localhost:2001", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewTestClient(conn)

	// Contact the server and print out its response.
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
