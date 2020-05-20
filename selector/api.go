package selector

import (
	"net/http"

	"github.com/eagraf/synchronizer/messenger"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

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
	// Input validation
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

	// Add to workers map
	worker := &Worker{
		UUID:         r.Header.Get("clientID"),
		Healthy:      true,
		Disconnected: false,
	}
	s.workers[worker.UUID] = worker

	// Subscribe to worker topic
	s.messenger.AddSubscription(worker.UUID, s)

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
