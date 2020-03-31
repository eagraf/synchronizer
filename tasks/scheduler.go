package tasks

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

// TaskScheduler manages the distribution of tasks to workers via a queue
type TaskScheduler struct {
	IntentQueue  chan *Intent             // Channel of incoming tasks
	CurrentTasks map[string]*TaskInstance // Map of currently running tasks
	TaskRegistry map[string]TaskType
	logger       *log.Logger
}

// Start listening for intents originated by the API/scheduler
func Start(taskRegistry map[string]TaskType) *TaskScheduler {

	// Create TaskScheduler
	var ts = TaskScheduler{
		IntentQueue:  make(chan *Intent),
		CurrentTasks: make(map[string]*TaskInstance),
		TaskRegistry: taskRegistry,
		logger:       log.New(os.Stdout, "TaskScheduler: ", log.LstdFlags),
	}

	// Listen loop
	// As of right now, this is the most naive algorithm possible for scheduling tasks
	go func() {
		for {
			intent := <-ts.IntentQueue
			switch intent.IntentType {
			case "setup":
				fmt.Println("Setup")
				go ts.handleSetup(intent)
			case "map":
				fmt.Println("Map")
			case "reduce":
				fmt.Println("Reduce")
			case "end":
				fmt.Println("End")
			case "default":
				fmt.Println("Default")
			}
		}
	}()

	return &ts
}

// Handle a setup intent
func (ts *TaskScheduler) handleSetup(intent *Intent) {
	// Make sure task type is registered
	_, ok := ts.TaskRegistry[intent.TaskType]
	if !ok {
		ts.logger.Println("Error: Task type is not in registry.")
		return
	}

	// Execute setup procedure for task
	taskInstance, mapIntents := ts.TaskRegistry[intent.TaskType].Setup(intent)

	// Generate a UUID for the task
	taskInstance.UUID = uuid.New().String()
	ts.CurrentTasks[taskInstance.UUID] = taskInstance

	// Add map intents to the queue
	for _, mi := range mapIntents {
		ts.IntentQueue <- mi
	}
}
