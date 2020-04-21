package tasks

import (
	"fmt"

	"github.com/google/uuid"
)

// TaskInstance models an ongoing task
type TaskInstance struct {
	UUID              string
	TaskType          string
	TaskSpecification TaskType `json:"-"`
	Config            TaskConfig
	intentQueue       chan *Intent `json:"-"` // Channel of incoming tasks
	PartialResults    []interface{}
	State             interface{}
	RequestTimes      []RequestTime
}

type RequestTime struct {
	workerUUID    string
	start         int64
	end           int64
	intermediates []IntermediateTime
}

type IntermediateTime struct {
	name  string
	start int64
	end   int64
}

// Start new instance of a task
func (ti *TaskInstance) Start(mapTaskQueue chan *Intent, input *map[string]interface{}) {
	go func() {
		for {
			intent := <-ti.intentQueue
			switch intent.IntentType {
			case "setup":
				fmt.Println("Setup")
				go ti.handleSetup(intent)
			case "map":
				fmt.Println("Map")
				// Find a better place to put this
				mapTaskQueue <- intent
				// Handle map task
			case "reduce":
				fmt.Println("Reduce")
			case "end":
				fmt.Println("End")
			case "default":
				fmt.Println("Default")
			}
		}
	}()
	// Send setup intent to this instance's intent queue
	ti.intentQueue <- ti.TaskSpecification.Initialize(ti.UUID, ti.Config, *input)
}

// Handle a setup intent
func (ti *TaskInstance) handleSetup(intent *Intent) {

	// Execute setup procedure for task
	taskInstance, mapIntents := ti.TaskSpecification.Setup(intent)

	// Generate a UUID for the task
	taskInstance.UUID = uuid.New().String()

	// Add map intents to the queue
	for _, mi := range mapIntents {
		ti.intentQueue <- mi
	}
}

// GetUUID implements a messenger subscriber method
func (ti *TaskInstance) GetUUID() string {
	return ti.UUID
}

// OnReceive implements a messenger subscriber method
func (ti *TaskInstance) OnReceive(topic string, m *map[string]interface{}) {
	fmt.Println("OnReceive", (*m)["start"], (*m)["end"])
}
