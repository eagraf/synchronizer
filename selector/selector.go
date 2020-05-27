package selector

import (
	"time"

	"github.com/eagraf/synchronizer/messenger"
	"github.com/eagraf/synchronizer/service"
	"github.com/google/uuid"
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

	rpcHandler := RPCService{}
	apiHandler := registerRoutes(s)

	service, err := si.StartService("Selector", rpcHandler, apiHandler)
	if err != nil {
		return nil, err
	}
	s.service = service

	// Return selector
	return s, nil
}

func (s *Selector) sendRegistrationResponse() {

}

func (s *Selector) sendHealthCheck(workerUUID string) error {
	// Send the health check message
	mb := new(messenger.MessageBuilder)
	requestID := uuid.New().String()
	m, err := mb.NewMessage(MessageHealthCheck, requestID).Done()
	if err != nil {
		return err
	}
	s.messenger.Send(workerUUID, m)

	// Timeout waits in a new thread
	go func() {
		time.Sleep(HealthCheckTimeout)

		// Check if timeout was successful
		rt := s.messenger.GetRequestRoundTrip(requestID)
		if rt != nil && rt.Response != nil {
			s.workers[workerUUID].Healthy = true
		}
		s.workers[workerUUID].Healthy = false
	}()
	return nil
}

func (s *Selector) sendHandoff() {

}

func (s *Selector) getWorker(workerUUID string) (*Worker, bool) {
	res, ok := s.workers[workerUUID]
	return res, ok
}
