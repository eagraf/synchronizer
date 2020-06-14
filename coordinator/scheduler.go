package coordinator

import (
	"github.com/eagraf/synchronizer/service"
)

type scheduler interface {
	scheduleWorkers(taskQueue []*Task, workerQueue []*service.WorkersResponse_Worker) *workerSchedule
	scheduleduleDataServers(jobs []*MapReduceJob, dataServers []*dataServer) *dataServerSchedule
	scheduleAggregators(jobs []*MapReduceJob, aggregators []*aggregator) *aggregatorSchedule
}

type schedule struct {
	workerSchedule     workerSchedule
	dataServerSchedule dataServerSchedule
	aggregatorSchedule aggregatorSchedule
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
	ID string
}

type aggregator struct {
	ID string
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

func (c *Coordinator) schedule() *schedule {
	res := new(schedule)

	// Assign tasks to workers
	ws := c.scheduler.scheduleWorkers(c.taskQueue, c.workers)
	// Clear taskQueue and workers
	c.taskQueue = nil
	c.workers = nil

	// Map active jobs into  a slice
	jobs := make([]*MapReduceJob, len(c.activeJobs))
	i := 0
	for _, j := range c.activeJobs {
		jobs[i] = j
		i++
	}

	// Task assignments need to be allocated to data servers and aggregators
	dataServerConnections, err := c.service.AllPeersOfType("Data Server")
	if err != nil {
		// Handle somehow
	}
	// Map connections into dataserver struct
	dataServers := make([]*dataServer, len(dataServerConnections))
	i = 0
	for _, ds := range dataServerConnections {
		dataServers[i] = &dataServer{ID: ds.Service.ID}
		i++
	}
	// Schedule data servers
	dss := c.scheduler.scheduleduleDataServers(jobs, dataServers)

	aggregatorConnections, err := c.service.AllPeersOfType("Aggregator")
	if err != nil {
		// Log it?
	}
	// Map connections into aggregator struct
	aggregators := make([]*aggregator, len(aggregatorConnections))
	i = 0
	for _, ag := range aggregatorConnections {
		aggregators[i] = &aggregator{ID: ag.Service.ID}
		i++
	}
	// Schedule aggregators
	as := c.scheduler.scheduleAggregators(jobs, aggregators)

	res.workerSchedule = *ws
	res.dataServerSchedule = *dss
	res.aggregatorSchedule = *as
	return res
}
