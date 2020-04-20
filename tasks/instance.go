package tasks

import (
	"fmt"
	"time"

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
	start         time.Time
	end           time.Time
	intermediates []IntermediateTime
}

type IntermediateTime struct {
	name  string
	start time.Time
	end   time.Time
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

// OnReceive implements a messenger subscriber
func (ti *TaskInstance) OnReceive(m *map[string]interface{}) {
	fmt.Println("OnReceive", m)

}
