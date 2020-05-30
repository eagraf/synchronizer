package coordinator

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
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
	JobUUID    string  `json"jobUUID"`
	JobType    string  `json:"jobType"`
	TaskSize   int     `json:"taskSize"`
	TaskNumber int     `json:"taskNumber"`
	Tasks      []*Task `json:"tasks"`
}

type Task struct {
	JobUUID   string `json:"jobUUID"`
	TaskIndex int    `json:"taskIndex"`
	TaskSize  int    `json:"taskSize"`
}

type MapReduceJobRequest struct {
	JobType    string `json:"jobType"`
	TaskSize   int    `json:"taskSize"`
	TaskNumber int    `json:"taskNumber"`
}

type MapReduceJobResponse struct {
	JobUUID    string `json:"jobUUID"`
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

	// Add a new active job
	job := &MapReduceJob{
		JobUUID:    uuid.New().String(),
		JobType:    body.JobType,
		TaskNumber: body.TaskNumber,
		TaskSize:   body.TaskSize,
		Tasks:      make([]*Task, body.TaskNumber),
	}

	// Break the job into task components, and add to coordinators task queue
	c.jobMutex.Lock()
	for i := 0; i < body.TaskNumber; i++ {
		newTask := &Task{
			TaskIndex: i,
			TaskSize:  body.TaskSize,
		}

		job.Tasks[i] = newTask
		c.taskQueue = append(c.taskQueue, newTask)
	}
	c.activeJobs[job.JobUUID] = job
	c.jobMutex.Unlock()

	// Marshal the job into a response
	res, err := json.Marshal(job)
	if err != nil {
		http.Error(w, "Failed to generate response json", 500)
	}
	// Write the response
	w.Write(res)
}
