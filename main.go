package main

import (
	"context"
	"net/http"

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
	/*fmt.Println("Starting synchronizer")

	var taskRegistry map[string]tasks.TaskType = make(map[string]tasks.TaskType, 0)
	taskRegistry["GOL"] = gameoflife.GOLTaskType

	messenger.InitializeMessenger()

	wm := workers.GetWorkerManager()
	wm.Start()

	//	_ = tasks.Start(taskRegistry, wm.MapTaskQueue)

	r := RegisterRoutes()
	log.Fatal(http.ListenAndServe(":2216", r))*/

	topology := make(map[string]map[string]bool)
	topology["Test"] = map[string]bool{"Test": true}
	_ = service.NewServicePool(topology)

	// Create new TestService
	/*rpc := TestServerImpl{}
	API := chi.NewRouter()
	API.Route("/test", func(r chi.Router) {
		API.Get("/", testEndpoint)
	})
	_, err := sp.StartService("Test", rpc, API)
	if err != nil {
		fmt.Println("fuck")
	}*/
}
