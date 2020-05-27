package selector

import (
	"time"

	"github.com/eagraf/synchronizer/messenger"
	"github.com/eagraf/synchronizer/service"
)

const (
	// MessageRegistrationResponse is a message type identifier for selector registration responses
	MessageRegistrationResponse string = "selector_registration_response"
	// MessageHealthCheck is a message type identifier for selector health checks
	MessageHealthCheck string = "selector_health_check"
	// MessageHandoff is a message type identifier for selector handoffs
	MessageHandoff string = "selector_handoff"
	// HealthCheckTimeout Timeout length (TODO exponential backoff)
	HealthCheckTimeout = 5 * time.Second
)

// Selector service struct
type Selector struct {
	workers    map[string]*Worker // TODO maybe use a sorted tree based structure -> when workers are transferred to coordinator, coord. just has to merge workers together from each selector
	messenger  *messenger.Messenger
	rpcHandler RPCService
	service    *service.Service
}

// Worker struct for mobile device
type Worker struct {
	UUID         string
	Healthy      bool
	Disconnected bool
}

func newSelector(si service.ServiceInitiator) (*Selector, error) {
	// Initialize selector
	var s *Selector = &Selector{
		workers:   make(map[string]*Worker),
		messenger: messenger.NewMessenger(),
	}

	rpcHandler := RPCService{
		selector: s,
	}
	apiHandler := registerRoutes(s)

	service, err := si.StartService("Selector", rpcHandler, apiHandler)
	if err != nil {
		return nil, err
	}
	s.service = service

	// Return selector
	return s, nil
}

func (s *Selector) getWorker(workerUUID string) (*Worker, bool) {
	res, ok := s.workers[workerUUID]
	return res, ok
}
