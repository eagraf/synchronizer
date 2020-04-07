package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/eagraf/synchronizer/workers"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

// RegisterRoutes defines routes for REST API using the chi router
// TODO i guess technically this should be singleton ensured as well
func RegisterRoutes() http.Handler {

	workerService := GetWorkerService()

	r := chi.NewRouter()
	r.Route("/health", func(r chi.Router) {
		r.Get("/", getHealth)
	})
	r.Route("/workers", func(r chi.Router) {
		r.Post("/", workerService.postWorker)
		r.Get("/", workerService.websocket)
		r.Delete("/{uuid}/", workerService.deleteWorker)
	})
	return r
}

// getHealth responds 200 if the service is up
func getHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}

// WorkerService is a service abstraction for endpoints involving workers
type WorkerService struct {
	wm *workers.WorkerManager
}

// GetWorkerService returns a singleton instance of the WorkerService
// TODO singleton guarantee
func GetWorkerService() *WorkerService {
	ws := WorkerService{
		workers.GetWorkerManager(),
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

func (ws *WorkerService) websocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the connection to
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

// PostWorker adds a new worker to the pool
func (ws *WorkerService) postWorker(w http.ResponseWriter, r *http.Request) {

	// Request body must follow this format.
	type PostBody struct {
		IP         string `json:"ip" bson:"ip"`
		WorkerType string `json:"workerType" bson:"workerType"`
	}

	var body PostBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid body: %s", err.Error()), 400)
		return
	}

	// Parse ip from body
	ip := net.ParseIP(body.IP)
	if ip == nil {
		http.Error(w, fmt.Sprintf("Invalid IP address: %s", string(body.IP)), 400)
		return
	}

	// TODO ensure that worker type is valid enumerable
	if body.WorkerType == "" {
		http.Error(w, fmt.Sprintf("Invalid worker type %s", string(body.WorkerType)), 400)
		return
	}

	// Add the worker
	uuid := ws.wm.AddWorker(ip, body.WorkerType)

	w.Write([]byte(uuid))
}

// DeleteWorker removes a worker from the pool
func (ws *WorkerService) deleteWorker(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	err := ws.wm.RemoveWorker(uuid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove worker: %s", err.Error()), 400)
		return
	}

	w.Write([]byte(""))
}
