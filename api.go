package main

import (
	"net/http"

	"github.com/eagraf/synchronizer/gameoflife"
	"github.com/eagraf/synchronizer/tasks"
	"github.com/eagraf/synchronizer/workers"
	"github.com/go-chi/chi"
)

// RegisterRoutes defines routes for REST API using the chi router
// TODO i guess technically this should be singleton ensured as well
func RegisterRoutes() http.Handler {
	var taskRegistry map[string]tasks.TaskType = make(map[string]tasks.TaskType, 0)
	taskRegistry["GOL"] = gameoflife.GOLTaskType

	workerService := workers.GetWorkerService()
	taskService := tasks.GetTaskService(taskRegistry, workers.GetWorkerManager().MapTaskQueue)

	r := chi.NewRouter()
	r.Route("/health", func(r chi.Router) {
		r.Get("/", getHealth)
	})
	r.Route("/workers", func(r chi.Router) {
		//r.Post("/", workerService.PostWorker)
		r.Get("/", workerService.RegisterWorker)
		r.Delete("/{uuid}/", workerService.DeleteWorker)
	})
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", taskService.PostTask)
		r.Get("/", taskService.GetTasks)
	})
	return r
}

// getHealth responds 200 if the service is up
func getHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}
