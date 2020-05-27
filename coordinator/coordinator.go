package coordinator

import (
	"fmt"
	"time"

	"github.com/eagraf/synchronizer/service"
)

type Scheduler interface {
	Schedule()
}

type Coordinator struct {
	service   *service.Service
	interval  time.Duration
	scheduler Scheduler
	round     uint
}

func newCoordinator(si service.ServiceInitiator) (*Coordinator, error) {
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

	// Start the interval timer
	go c.schedulingInterval()

	return c, nil
}

func (c *Coordinator) schedulingInterval() {
	for {
		c.round++
		c.schedule()
		time.Sleep(c.interval)
	}
}

func (c *Coordinator) schedule() {
	fmt.Println(c.round)
}
