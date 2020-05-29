package coordinator

type Scheduler interface {
	Schedule()
}

type job struct {
	uuid  string
	tasks map[string]*task
}

type task struct {
	uuid  string
	jobID string
}

func (c *Coordinator) schedule() {
	// Unimplimented
}
