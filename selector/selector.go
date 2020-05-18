package selector

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"

	"github.com/go-chi/chi"
)

type selector struct {
	workers []worker // Don't use a list of pointers so that workers can be easilly encoded
}

type worker struct {
	UUID string
}

type workerRequest struct {
	numRequested int
}

type workerReply struct {
	workers []worker
}

type handoffRequest struct {
}

type handoffReply struct {
}

func newSelector(apiPort int, rpcPort int) (*selector, error) {
	// Initialize selector
	var s *selector = new(selector)

	// Start api
	routes := registerRoutes(s)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(apiPort))
	if err != nil {
		log.Fatal("listen error (api):", err)
		return nil, err
	}
	go http.Serve(listener, routes)

	// Start rpc
	rpc.Register(s)
	rpc.HandleHTTP()
	listener, err = net.Listen("tcp", ":"+strconv.Itoa(rpcPort))
	if err != nil {
		log.Fatal("listen error (rpc):", err)
		return nil, err
	}
	go http.Serve(listener, nil)

	// Return selector
	return s, nil
}

// RegisterRoutes defines routes for REST API using the chi router
func registerRoutes(s *selector) http.Handler {

	r := chi.NewRouter()
	r.Route("/websocket", func(r chi.Router) {
		r.Get("/", s.websocket)
	})
	return r
}

// HTTP endpoint that promotes an HTTP request to a full WebSocket connection
func (s *selector) websocket(w http.ResponseWriter, r *http.Request) {
}

// RPC interface
func (s *selector) GetWorkers(request workerRequest, reply *workerReply) error {
	return nil
}

func (s *selector) Handoff(request handoffRequest, reply *handoffRequest) error {
	return nil
}
