package coordinator

import (
	"github.com/eagraf/synchronizer/service"
)

type scheduler interface {
	// TODO schedule needs to return list of unscheduled workers
	schedule(taskQueue []*Task, workerQueue []*service.WorkersResponse_Worker) *schedule
}

type schedule struct {
	assignments       map[string]map[string][]int // Key 1: Worker, Key 2: Job, Index: Task index
	unassignedWorkers map[string]bool
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
	// Assign tasks to workers
	c.scheduler.schedule(c.taskQueue, c.workers)
	// Task assignments need to be allocated to data servers and aggregators
}

type naiveScheduler struct{}

func (ns *naiveScheduler) schedule(taskQueue []*Task, workerQueue []*service.WorkersResponse_Worker) *schedule {
	res := &schedule{
		assignments:       make(map[string]map[string][]int),
		unassignedWorkers: make(map[string]bool),
	}
	for i, task := range taskQueue {
		worker := workerQueue[i%len(workerQueue)]

		if _, ok := res.assignments[worker.WorkerUUID]; ok == false {
			res.assignments[worker.WorkerUUID] = make(map[string][]int)
		}
		if _, ok := res.assignments[worker.WorkerUUID][task.JobUUID]; ok == false {
			res.assignments[worker.WorkerUUID][task.JobUUID] = make([]int, 0)
		}

		res.assignments[worker.WorkerUUID][task.JobUUID] = append(res.assignments[worker.WorkerUUID][task.JobUUID], task.TaskIndex)
	}

	// Populate unassigned workers
	for i := len(taskQueue); i < len(workerQueue); i++ {
		res.unassignedWorkers[workerQueue[i].WorkerUUID] = true
	}
	return res
}
