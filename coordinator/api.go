package coordinator

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
)

// RegisterRoutes defines routes for REST API using the chi router
func registerRoutes(c *Coordinator) http.Handler {
	r := chi.NewRouter()
	r.Route("/jobs", func(r chi.Router) {
		r.Post("/", c.createMapReduceJob)
	})
	return r
}

type MapReduceJob struct {
	JobType    string  `json:"jobType"`
	TaskSize   int     `json:"taskSize"`
	TaskNumber int     `json:"taskNumber"`
	Tasks      []*Task `json:"tasks"`
}

type Task struct {
	TaskIndex int `json:"taskIndex"`
	TaskSize  int `json:"taskSize"`
}

type MapReduceJobRequest struct {
	JobType    string `json:"jobType"`
	TaskSize   int    `json:"taskSize"`
	TaskNumber int    `json:"taskNumber"`
}

// creates a new map reduce job
func (c *Coordinator) createMapReduceJob(w http.ResponseWriter, r *http.Request) {

	// Read body into MapReduceJobRequest
	var body MapReduceJobRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error decoding body", 400)
	}
	err = json.Unmarshal(buf, &body)
	if err != nil {
		http.Error(w, "Incorrect fields in input", 400)
	}

	// Generate subtasks
	w.Write([]byte("hello"))
}
