package tasks

// TaskConfig details a set of configurations for a specific task run
type TaskConfig struct {
	NumWorkers int
}

// Intent is a request for the scheduler to perform an action
// Vaguely reminiscent of the Android Intent pattern
type Intent struct {
	IntentType string
	TaskType   string
	TaskUUID   string
	Config     TaskConfig
	Input      interface{}
}

// TaskInstance models an ongoing task
type TaskInstance struct {
	UUID           string
	Config         TaskConfig
	PartialResults []interface{}
	State          interface{}
}

// TaskType holds the setup and reduce functions for a type of task
type TaskType struct {
	Name   string
	Setup  func(i *Intent) (*TaskInstance, []*Intent)
	Reduce func() *Intent
}
