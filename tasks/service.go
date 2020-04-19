package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

// TaskService handles API calls related to tasks
type TaskService struct {
	taskScheduler *TaskScheduler
}

// GetTaskService returns an instance of the TaskService
func GetTaskService() *TaskService {
	ts := TaskService{
		taskScheduler: GetTaskSchedulerSingleton(),
	}
	return &ts
}

// GetTasks gets all ongoing tasks
func (ts *TaskService) GetTasks(w http.ResponseWriter, r *http.Request) {
	buffer, err := json.Marshal(ts.taskScheduler.CurrentTasks)
	if err != nil {
		http.Error(w, "Failed to marshal tasks", 500)
		return
	}

	w.Write(buffer)
}

// PostTask asks the synchronizer to begin performing a new task
func (ts *TaskService) PostTask(w http.ResponseWriter, r *http.Request) {

	// Read body into buffer
	var buffer []byte
	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	type PostTaskBody struct {
		TaskType   string
		Input      map[string]interface{}
		NumWorkers int
	}

	// Read the body into map
	var body PostTaskBody
	err = json.Unmarshal(buffer, &body)
	if err != nil {
		fmt.Println(err)
	}

	// TODO validate body input

	// Generate uuid
	uuid := uuid.New().String()

	var intent *Intent
	// Generate the initial setup intent
	switch body.TaskType {
	case "GOL":
		intent = ts.taskScheduler.TaskRegistry["GOL"].Initialize(uuid, TaskConfig{body.NumWorkers}, body.Input)
	default:
		http.Error(w, "Task type not recognized", 400)
		return
	}

	// Enqueue with the task scheduler
	taskScheduler.IntentQueue <- intent
	w.Write([]byte(uuid))
}
