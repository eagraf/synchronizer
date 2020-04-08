package workers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

// WorkerService is a service abstraction for endpoints involving workers
type WorkerService struct {
	wm *WorkerManager
}

// GetWorkerService returns a singleton instance of the WorkerService
// TODO singleton guarantee
func GetWorkerService() *WorkerService {
	ws := WorkerService{
		GetWorkerManager(),
	}
	return &ws
}

// upgrader promotes a standard HTTP/HTTPS connection to a websocket connection
// TODO implement CheckOrigin
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true },
}

// RegisterWorker notifies the WorkerManager of a new worker, and promotes the request to a websocket connection
func (ws *WorkerService) RegisterWorker(w http.ResponseWriter, r *http.Request) {
	// TODO any necessary authentication code

	// Upgrade connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// Get worker type
	// TODO ensure worker type is valid
	query := r.URL.Query()
	workerType, ok := query["workertype"]
	if !ok {
		conn.WriteMessage(websocket.TextMessage, []byte("No workertype included"))
		conn.Close()
		return
	}

	//fmt.Println("I am here")
	uuid := ws.wm.AddWorker(workerType[0])
	conn.WriteMessage(websocket.TextMessage, []byte(uuid))

	// Handoff connection to messenger

}

// DeleteWorker removes a worker from the pool
func (ws *WorkerService) DeleteWorker(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	err := ws.wm.RemoveWorker(uuid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove worker: %s", err.Error()), 400)
		return
	}

	w.Write([]byte(""))
}
