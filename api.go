package main

import (
	"net/http"

	"github.com/eagraf/synchronizer/workers"
	"github.com/go-chi/chi"
)

// RegisterRoutes defines routes for REST API using the chi router
// TODO i guess technically this should be singleton ensured as well
func RegisterRoutes() http.Handler {

	workerService := workers.GetWorkerService()

	r := chi.NewRouter()
	r.Route("/health", func(r chi.Router) {
		r.Get("/", getHealth)
	})
	r.Route("/workers", func(r chi.Router) {
		r.Post("/", workerService.PostWorker)
		r.Get("/", workerService.Websocket)
		r.Delete("/{uuid}/", workerService.DeleteWorker)
	})
	return r
}

// getHealth responds 200 if the service is up
func getHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}
