package coordinator

import (
	"fmt"
	"sync"

	"github.com/eagraf/synchronizer/service"
)

type scheduler interface {
	scheduleWorkers(taskQueue []*Task, workerQueue []*service.WorkersResponse_Worker) *workerSchedule
	scheduleDataServers(jobs []*MapReduceJob, dataServers []*dataServer) *dataServerSchedule
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
	dataServerConnections, err := c.service.AllPeersOfType("Data_Server")
	if err != nil {
		c.service.Log("SchedulingError", "Data Serversunavailable")
	}
	// Map connections into dataserver struct
	dataServers := make([]*dataServer, len(dataServerConnections))
	i = 0
	for _, ds := range dataServerConnections {
		dataServers[i] = &dataServer{ID: ds.Service.ID}
		i++
	}
	// Schedule data servers
	dss := c.scheduler.scheduleDataServers(jobs, dataServers)

	aggregatorConnections, err := c.service.AllPeersOfType("Aggregator")
	if err != nil {
		c.service.Log("SchedulingError", "Aggregators unavailable")
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

	errs := c.sendToDataServers(dss)
	if len(errs) != 0 {
		c.service.Log("SchedulingError", fmt.Sprintf("%d errors sending to Data Servers", len(errs)))
	}

	fmt.Println("hello")
	errs = c.sendToAggregators(as)
	if len(errs) != 0 {
		c.service.Log("SchedulingError", fmt.Sprintf("%d errors sending to Aggregators", len(errs)))
		fmt.Println(errs)
	}

	return res
}

// TODO later
func sendToSelectors() {

}

type dsScheduleRequest = service.DataServerReceiveScheduleRequest
type dsSchedule = service.DataServerReceiveScheduleRequest_Schedule
type dsScheduleJob = service.DataServerReceiveScheduleRequest_Schedule_Job
type dsScheduleWorker = service.DataServerReceiveScheduleRequest_Schedule_Worker

// Concurrently send schedule to each data server. Return list of errors and count.
func (c *Coordinator) sendToDataServers(schedule *dataServerSchedule) []error {
	// Return list of errors
	errs := make([]error, 0)

	// Use waitgroup to block until all requests have completed
	var wg sync.WaitGroup

	for ds, jobs := range schedule.assignments {
		// Make request
		req := &dsScheduleRequest{}
		sched := &dsSchedule{
			Jobs: make([]*dsScheduleJob, len(jobs)),
		}
		req.Schedule = sched

		// Fill in jobs
		for i, job := range jobs {
			sched.Jobs[i] = &dsScheduleJob{
				JobUUID:    c.activeJobs[job].JobUUID,
				JobType:    c.activeJobs[job].JobType,
				TaskSize:   int32(c.activeJobs[job].TaskSize),
				TaskNumber: int32(c.activeJobs[job].TaskNumber),
			}
		}
		// Send to dataserver
		dsConn, err := c.service.GetPeer("Data_Server", ds)
		if err != nil {
			errs = append(errs, err)
		}
		reply := service.DataServerReceiveScheduleResponse{}
		// Make request with callback
		c.service.UniCast(dsConn, service.DataServerReceiveSchedule, req, &reply, func(reply interface{}, err error) {
			// Add each call thread to waitgroup, and then remove when done
			wg.Add(1)
			defer wg.Done()

			if err != nil {
				errs = append(errs, err)
			}
		})
	}
	// Wait for all calls to complete
	wg.Wait()
	return errs
}

type agScheduleRequest = service.AggregatorReceiveScheduleRequest
type agSchedule = service.AggregatorReceiveScheduleRequest_Schedule
type agScheduleJob = service.AggregatorReceiveScheduleRequest_Schedule_Job
type agScheduleWorker = service.AggregatorReceiveScheduleRequest_Schedule_Worker

func (c *Coordinator) sendToAggregators(schedule *aggregatorSchedule) []error {
	// Return list of errors
	errs := make([]error, 0)

	// Use waitgroup to block until all requests have completed
	var wg sync.WaitGroup

	for as, jobs := range schedule.assignments {
		fmt.Println("bb")
		// Make request
		req := &agScheduleRequest{}
		sched := &agSchedule{
			Jobs: make([]*agScheduleJob, len(jobs)),
		}
		req.Schedule = sched

		// Fill in jobs
		for i, job := range jobs {
			sched.Jobs[i] = &agScheduleJob{
				JobUUID:    c.activeJobs[job].JobUUID,
				JobType:    c.activeJobs[job].JobType,
				TaskSize:   int32(c.activeJobs[job].TaskSize),
				TaskNumber: int32(c.activeJobs[job].TaskNumber),
			}
		}
		// Send to dataserver
		agConn, err := c.service.GetPeer("Aggregator", as)
		if err != nil {
			errs = append(errs, err)
		}
		reply := service.AggregatorReceiveScheduleResponse{}
		// Make request with callback
		c.service.UniCast(agConn, service.AggregatorReceiveSchedule, req, &reply, func(reply interface{}, err error) {
			// Add each call thread to waitgroup, and then remove when done
			wg.Add(1)
			defer wg.Done()

			if err != nil {
				errs = append(errs, err)
			}
		})
	}
	// Wait for all calls to complete
	wg.Wait()
	return errs

}
