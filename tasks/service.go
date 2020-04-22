package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

var taskService *TaskService

// TaskService handles API calls related to tasks
type TaskService struct {
	TaskRegistry map[string]TaskType
	CurrentTasks map[string]*TaskInstance // Map of currently running tasks
	MapTaskQueue chan *Intent
}

// InitializeTaskService initializes the TaskService singleton
func InitializeTaskService(taskRegistry map[string]TaskType, mapTaskQueue chan *Intent) *TaskService {
	if taskService != nil {
		panic("TaskService has already been initialized")
	}
	ts := TaskService{
		CurrentTasks: make(map[string]*TaskInstance),
		TaskRegistry: taskRegistry,
		MapTaskQueue: mapTaskQueue,
	}
	taskService = &ts
	return taskService
}

// GetTaskServiceSingleton returns the TaskService singleton
func GetTaskServiceSingleton() *TaskService {
	if taskService == nil {
		panic("TaskService has not been initialized yet")
	}
	return taskService
}

// GetTasks gets all ongoing tasks
func (ts *TaskService) GetTasks(w http.ResponseWriter, r *http.Request) {
	buffer, err := json.Marshal(ts.CurrentTasks)
	if err != nil {
		fmt.Println(err)
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

	// Create new task instance
	ti := TaskInstance{
		UUID:              uuid,
		TaskType:          body.TaskType,
		TaskSpecification: ts.TaskRegistry[body.TaskType], // TaskRegistry
		Config: TaskConfig{
			body.NumWorkers,
		},
		intentQueue:    make(chan *Intent),
		PartialResults: make([]interface{}, 0),
		RequestTimes:   make([]RequestTime, 0),
	}

	// Start the task
	ti.Start(ts.MapTaskQueue, &body.Input)
	ts.CurrentTasks[uuid] = &ti

	// Write uuid as response
	w.Write([]byte(uuid))
}

// GetTaskInstance returns a current task instance
// TODO this probably should be refactored into a model type
func (ts *TaskService) GetTaskInstance(uuid string) (*TaskInstance, bool) {
	val, ok := ts.CurrentTasks[uuid]
	return val, ok
}
