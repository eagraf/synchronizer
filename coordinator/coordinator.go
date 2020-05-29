package coordinator

import (
	"time"

	"github.com/eagraf/synchronizer/service"
)

// Coordinator service struct
type Coordinator struct {
	service   *service.Service
	interval  time.Duration
	scheduler Scheduler
	round     uint
	workers   []*service.WorkersResponse_Worker // Helpful for tests
}

// NewCoordinator creates a new coordinator service
func NewCoordinator(si service.ServiceInitiator) (*Coordinator, error) {
	// Create new Coordinator
	var c *Coordinator = &Coordinator{
		interval: 10 * time.Second,
		round:    0,
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

func (c *Coordinator) schedulingInterval() {
	for {
		c.round++

		// Step 1: collect workers from selectors
		//workerSet := make([]*service.WorkersResponse_Worker, 0)

		selectors, err := c.service.AllPeersOfType("Selector")
		if err != nil {
			// There are no selectors therefore nothing to schedule
			return
		}

		replys := make([]interface{}, len(selectors))
		for i := range replys {
			replys[i] = &service.WorkersResponse{}
		}
		req := &service.WorkersRequest{
			All:    true,
			Number: 0,
		}
		responses, errs := c.service.MultiCast(selectors, service.SelectorGetWorkers, req, replys)
		if len(errs) != 0 {

		}

		// Aggregate the responses
		workers := make([]*service.WorkersResponse_Worker, 0)
		for _, r := range responses {
			workers = append(workers, r.(*service.WorkersResponse).Workers...)
		}
		c.workers = workers

		// Step 2: produce a schedule
		c.schedule()

		// Step 3: notify all relevent services about the schedule

		// Step 4: wait for the next scheduling period
		time.Sleep(c.interval)
	}
}
