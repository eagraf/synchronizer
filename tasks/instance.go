package tasks

import (
	"fmt"
	"time"

	"github.com/eagraf/synchronizer/messenger"
	"github.com/google/uuid"
)

// TaskInstance models an ongoing task
type TaskInstance struct {
	UUID              string
	TaskType          string
	TaskSpecification TaskType `json:"-"`
	Config            TaskConfig
	intentQueue       chan *Intent  `json:"-"` // Channel of incoming tasks
	PartialResults    []interface{} `json:"-"`
	State             interface{}
	Input             interface{}
	RequestTimes      []RequestTime
	StartTime         int64
	EndTime           int64
	subscriptions     map[string]bool
}

type RequestTime struct {
	WorkerUUID    string
	Start         int64
	End           int64
	Intermediates []IntermediateTime
}

type IntermediateTime struct {
	Name  string
	Start int64
	End   int64
}

// Start new instance of a task
func (ti *TaskInstance) Start(mapTaskQueue chan *Intent, input *map[string]interface{}) {
	// TODO find a better place to initalize this
	ti.subscriptions = make(map[string]bool)
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

// AddSubscription implements a messenger subscriber method
func (ti *TaskInstance) AddSubscription(topic string) {
	ti.subscriptions[topic] = true
}

// GetUUID implements a messenger subscriber method
func (ti *TaskInstance) GetUUID() string {
	return ti.UUID
}

// OnReceive implements a messenger subscriber method
func (ti *TaskInstance) OnReceive(topic string, m *map[string]interface{}) {

	ti.PartialResults = append(ti.PartialResults, m)
	// TODO this is super hacky dont do this
	//ti.PartialResults[topic+"/"+len(ti.PartialResults)] = m

	// Update the request time list
	ti.RequestTimes = append(ti.RequestTimes, RequestTime{
		WorkerUUID: topic,
		Start:      int64((*m)["outer_start"].(int64)),
		End:        int64((*m)["outer_end"].(int64)),
		Intermediates: []IntermediateTime{
			IntermediateTime{
				Name:  "execution",
				Start: int64((*m)["start"].(float64)),
				End:   int64((*m)["end"].(float64)),
			},
		},
	})

	if len(ti.PartialResults) == ti.Config.NumWorkers {
		ti.EndTime = time.Now().UnixNano() / int64(time.Millisecond)
		// Stop listening to all subscriptions
		for s := range ti.subscriptions {
			messenger.GetMessengerSingleton().RemoveSubscriber(ti, s)
		}

		// Trigger reduce intent
	}
}

// OnSend implements a messenger subscriber method
func (ti *TaskInstance) OnSend(topic string) {
	if ti.StartTime == 0 {
		ti.StartTime = time.Now().UnixNano() / int64(time.Millisecond)
	}
}

// OnClose implements a messenger subscriber method
// TODO might want to unify this with messenger.RemoveSubscriber
func (ti *TaskInstance) OnClose(topic string) {
	delete(ti.subscriptions, topic)

	// TODO how does the task handle this?
}
