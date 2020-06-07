package main

import (
	"context"
	"net/http"
	"time"

	"github.com/eagraf/synchronizer/coordinator"
	"github.com/eagraf/synchronizer/messenger"
	"github.com/eagraf/synchronizer/selector"
	"github.com/eagraf/synchronizer/service"
	"google.golang.org/grpc"
)

// TestServer implementation
type TestServerImpl struct{}

func (tsi *TestServerImpl) TestRPC(ctx context.Context, in *service.Ping, opts ...grpc.CallOption) (*service.Pong, error) {
	res := &service.Pong{
		Message: "Hello",
	}
	return res, nil
}

// Test service external API handler
func testEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hell"))
}

func main() {
	sp := service.NewServicePool(4000, service.DefaultTopology)
	selector.NewSelector(sp)
	selector.NewSelector(sp)
	selector.NewSelector(sp)
	coordinator.NewCoordinator(sp)

	time.Sleep(5 * time.Second)

	messenger.NewTestClient("http://localhost:4000/websocket/", "client1")
	messenger.NewTestClient("http://localhost:4002/websocket/", "client2")
	messenger.NewTestClient("http://localhost:4004/websocket/", "client3")
	select {}

}
