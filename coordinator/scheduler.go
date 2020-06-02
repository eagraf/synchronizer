package coordinator

import (
	"github.com/eagraf/synchronizer/service"
)

type scheduler interface {
	scheduleWorkers(taskQueue []*Task, workerQueue []*service.WorkersResponse_Worker) *workerSchedule
	scheduleduleDataServers(jobs []*MapReduceJob, dataServers []*dataServer) *dataServerSchedule
	scheduleAggregator(jobs []*MapReduceJob, aggregators []*aggregator) *aggregatorSchedule
}

type workerAssignments = map[string]map[string][]int // Key 1: Worker, Key 2: Job, Index: Task index

type workerSchedule struct {
	assignments       workerAssignments
	unassignedWorkers map[string]bool
}

type dataServerSchedule struct {
	assignments map[string][]string // Map of data server ids, list of jobUUIDs
}

type aggregatorSchedule struct {
	assignments map[string][]string // Map of data server ids, list of jobUUIDs
}

// Helper types that contain important scheduling information
type dataServer struct {
	UUID string
}

type aggregator struct {
	UUID string
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
	c.scheduler.scheduleWorkers(c.taskQueue, c.workers)
	// Clear taskQueue and workers
	c.taskQueue = nil
	c.workers = nil
	// Task assignments need to be allocated to data servers and aggregators
}

type naiveScheduler struct{}

func (ns *naiveScheduler) scheduleWorkers(taskQueue []*Task, workerQueue []*service.WorkersResponse_Worker) *workerSchedule {
	res := &workerSchedule{
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

func (ns *naiveScheduler) scheduleDataServers(jobs []*MapReduceJob, dataServers []*dataServer) *dataServerSchedule {
	res := &dataServerSchedule{
		assignments: make(map[string][]string),
	}

	// Initialize each data server in assignments map
	for _, ds := range dataServers {
		res.assignments[ds.UUID] = make([]string, 0)
	}
	// Populate assignments
	for i, job := range jobs {
		dsUUID := dataServers[i%len(dataServers)].UUID
		res.assignments[dsUUID] = append(res.assignments[dsUUID], job.JobUUID)
	}
	return res
}

func (ns *naiveScheduler) scheduleAggregators(jobs []*MapReduceJob, aggregators []*aggregator) *aggregatorSchedule {
	res := &aggregatorSchedule{
		assignments: make(map[string][]string),
	}

	// Initialize each data server in assignments map
	for _, ds := range aggregators {
		res.assignments[ds.UUID] = make([]string, 0)
	}
	// Populate assignments
	for i, job := range jobs {
		agUUID := aggregators[i%len(aggregators)].UUID
		res.assignments[agUUID] = append(res.assignments[agUUID], job.JobUUID)
	}
	return res

}
