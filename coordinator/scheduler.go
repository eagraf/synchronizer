package coordinator

type Scheduler interface {
	Schedule()
}

// MapReduceJob is a basic job where tasks are easilly subdivided and distributed to workers
type MapReduceJob struct {
	JobUUID    string  `json:"jobUUID"`
	JobType    string  `json:"jobType"`
	TaskSize   int     `json:"taskSize"`
	TaskNumber int     `json:"taskNumber"`
	Tasks      []*Task `json:"tasks"`
}

// Task is a subunit of a job
type Task struct {
	JobUUID   string `json:"jobUUID"`
	TaskIndex int    `json:"taskIndex"`
	TaskSize  int    `json:"taskSize"`
}

func (c *Coordinator) schedule() {
	// Unimplimented
}
