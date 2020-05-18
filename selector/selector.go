package selector

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"

	"github.com/eagraf/synchronizer/messenger"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

const (
	// MessageRegistrationResponse is a message type identifier for selector registration responses
	MessageRegistrationResponse string = "selector_registration_response"
	// MessageHealthCheck is a message type identifier for selector health checks
	MessageHealthCheck string = "selector_health_check"
	// MessageHandoff is a message type identifier for selector handoffs
	MessageHandoff string = "selector_handoff"
)

type Selector struct {
	workers   []Worker // Don't use a list of pointers so that workers can be easilly encoded
	messenger *messenger.Messenger
}

type Worker struct {
	UUID         string
	Healthy      bool
	Disconnected bool
}

type WorkerRequest struct {
	numRequested int
}

type WorkerResponse struct {
	workers []Worker
}

type HandoffRequest struct {
}

type HandoffResponse struct {
}

func newSelector(apiPort int, rpcPort int) (*Selector, error) {
	// Initialize selector
	var s *Selector = new(Selector)
	m := messenger.NewMessenger()
	s.messenger = m

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
func registerRoutes(s *Selector) http.Handler {

	r := chi.NewRouter()
	r.Route("/websocket", func(r chi.Router) {
		r.Get("/", s.websocket)
	})
	return r
}

// HTTP endpoint that promotes an HTTP request to a full WebSocket connection
func (s *Selector) websocket(w http.ResponseWriter, r *http.Request) {
	if id := r.Header.Get("clientID"); id == "" {
		http.Error(w, "Missing clientID header.", 400)
		return
	}
	err := s.messenger.AddConnection(r.Header.Get("clientID"), w, r)
	if err != nil {
		http.Error(w, "Failed to add connection: "+err.Error(), 500)
		return
	}
	// Otherwise, websocket connection is managed by messenger
	// TODO should there be a return?

	// Send registration message with session id
	mb := new(messenger.MessageBuilder)
	m, err := mb.NewMessage(MessageRegistrationResponse, uuid.New().String()).
		AddHeader("session_id", uuid.New().String()).
		Done()

	if err != nil {
		http.Error(w, "Failed to add connection: "+err.Error(), 500)
		return
	}

	// Send registration response
	s.messenger.Send(r.Header.Get("clientID"), m)
}

// RPC interface
func (s *Selector) GetWorkers(request WorkerRequest, reply *WorkerResponse) error {
	return nil
}

func (s *Selector) Handoff(request HandoffRequest, reply *HandoffRequest) error {
	return nil
}

func (s *Selector) GetIdentifier() {

}

func (s *Selector) OnReceive(topic string, message *messenger.Message) {

}

/*
 * Types of sends:
 *   (1) Registration Response (selector_registration_response)
 *   (2) Health Check (selector_health_check)
 *   (3) Handoff (selector_handoff)
 */
func (s *Selector) OnSend(topic string, message *messenger.Message) {

}

func (s *Selector) OnClose(topic string) {

}

func (s *Selector) sendRegistrationResponse() {

}

func (s *Selector) sendHealthCheck(workerUUID string) error {
	mb := new(messenger.MessageBuilder)
	m, err := mb.NewMessage(MessageHealthCheck, uuid.New().String()).Done()
	if err != nil {
		return err
	}
	s.messenger.Send(workerUUID, m)
	return nil
}

func (s *Selector) sendHandoff() {

}

func (s *Selector) getWorker(workerUUID string) *Worker {
	return nil
}
