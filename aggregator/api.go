package aggregator

import (
	"fmt"
	"net/http"

	"github.com/eagraf/synchronizer/messenger"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// RegisterRoutes defines routes for REST API using the chi router
func registerRoutes(ds *Aggregator) http.Handler {
	r := chi.NewRouter()
	r.Route("/websocket", func(r chi.Router) {
		r.Get("/", ds.websocket)
	})
	return r
}

// HTTP endpoint that promotes an HTTP request to a full WebSocket connection
func (a *Aggregator) websocket(w http.ResponseWriter, r *http.Request) {
	// Input validation
	if id := r.Header.Get("clientID"); id == "" {
		http.Error(w, "Missing clientID header.", 400)
		return
	}
	err := a.messenger.AddConnection(r.Header.Get("clientID"), w, r)
	if err != nil {
		fmt.Println(err.Error())
		//http.Error(w, "Failed to add connection: "+err.Error(), 500)
		return
	}
	// Otherwise, websocket connection is managed by messenger
	// TODO should there be a return?

	// Add to workers map
	worker := &Worker{
		UUID: r.Header.Get("clientID"),
	}
	a.workers[worker.UUID] = worker

	// Subscribe to worker topic
	a.messenger.AddSubscription(worker.UUID, a)

	// Send registration message with session id
	mb := new(messenger.MessageBuilder)
	m, err := mb.NewMessage(MessageInitiateDataTransfer, uuid.New().String()).
		AddHeader("session_id", uuid.New().String()).
		Done()

	if err != nil {
		fmt.Println(err.Error())
		//http.Error(w, "Failed to add connection: "+err.Error(), 500)
		return
	}

	// Send registration response
	a.messenger.Send(r.Header.Get("clientID"), m)
}
