package coordinator

import (
	"fmt"
	"sync"
	"time"

	"github.com/eagraf/synchronizer/service"
)

// Coordinator service struct
type Coordinator struct {
	service    *service.Service
	interval   time.Duration
	scheduler  Scheduler
	round      uint
	workers    []*service.WorkersResponse_Worker // Helpful for tests
	activeJobs map[string]*MapReduceJob
	taskQueue  []*Task
	jobMutex   sync.Mutex
}

// NewCoordinator creates a new coordinator service
func NewCoordinator(si service.ServiceInitiator) (*Coordinator, error) {
	// Create new Coordinator
	var c *Coordinator = &Coordinator{
		interval:   10 * time.Second,
		round:      0,
		activeJobs: make(map[string]*MapReduceJob),
		taskQueue:  make([]*Task, 0, 1<<20), // Give a rather large initial capacity
		jobMutex:   sync.Mutex{},
	}

	// Setup service
	apiHandler := registerRoutes(c)
	service, err := si.StartService("Coordinator", nil, apiHandler)
	if err != nil {
		return nil, err
	}
	c.service = service

	// TODO this should probably be moved
	// Link up to other services
	si.ConnectService(c.service)

	// Start the interval timer
	go c.schedulingInterval()

	return c, nil
}

// Workers gets the list of workers scheduled in this interval
func (c *Coordinator) Workers() []*service.WorkersResponse_Worker {
	return c.workers
}

func (c *Coordinator) schedulingInterval() {
	for {
		c.round++

		// Step 1: collect workers from selectors
		workers, err := c.getWorkers()
		if err != nil {
			// Ignore errors for now
			// TODO some sort of error handling. Check health of selectors?
		}
		c.workers = workers

		// Step 2: produce a schedule
		c.schedule()

		// Step 3: notify all relevent services about the schedule

		// Step 4: wait for the next scheduling period
		time.Sleep(c.interval)
	}
}

func (c *Coordinator) getWorkers() ([]*service.WorkersResponse_Worker, error) {
	selectors, err := c.service.AllPeersOfType("Selector")
	if err != nil {
		// There are no selectors therefore nothing to schedule
		return []*service.WorkersResponse_Worker{}, nil
	}

	replys := make([]interface{}, len(selectors))
	for i := range replys {
		replys[i] = &service.WorkersResponse{}
	}
	req := &service.WorkersRequest{
		All: true,
	}
	// Make the multicast request
	responses, errs := c.service.MultiCast(selectors, service.SelectorGetWorkers, req, replys)
	err = nil
	if len(errs) != 0 {
		err = fmt.Errorf("There were %d errors returned by GetWorkers MultiCast, the first is %s", len(errs), errs[0].Error())
	}

	// Aggregate the responses
	workers := make([]*service.WorkersResponse_Worker, 0)
	for _, r := range responses {
		if wr, ok := r.(*service.WorkersResponse); ok == true {
			// Only include workers that are not errors
			workers = append(workers, wr.Workers...)
		}
	}
	return workers, err
}
