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

// TaskType holds the setup and reduce functions for a type of task
type TaskType struct {
	Name       string
	Initialize func(UUID string, config TaskConfig, input map[string]interface{}) *Intent
	Setup      func(i *Intent) (*TaskInstance, []*Intent)
	Reduce     func() *Intent
}
