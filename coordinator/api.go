package coordinator

import (
	"net/http"

	"github.com/go-chi/chi"
)

// RegisterRoutes defines routes for REST API using the chi router
func registerRoutes(c *Coordinator) http.Handler {
	r := chi.NewRouter()
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", c.createTask)
	})
	return r
}

func (c *Coordinator) createTask(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
