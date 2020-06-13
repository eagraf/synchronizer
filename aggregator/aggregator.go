package aggregator

import (
	"github.com/eagraf/synchronizer/messenger"
	"github.com/eagraf/synchronizer/service"
)

// TODO move these into messenger package
const (
	MessageInitiateDataTransfer string = "InitiateDataTransfer"
)

// Aggregator type contains state for this service
type Aggregator struct {
	workers    map[string]*Worker // TODO maybe use a sorted tree based structure -> when workers are transferred to coordinator, coord. just has to merge workers together from each selector
	messenger  *messenger.Messenger
	rpcHandler RPCService
	service    *service.Service
}

// Worker models a connected worker
type Worker struct {
	UUID string
}

// NewAggregator creates a new data server object
func NewAggregator(si service.ServiceInitiator) (*Aggregator, error) {
	a := &Aggregator{
		workers:   make(map[string]*Worker),
		messenger: messenger.NewMessenger(),
	}
	// Setup service
	rpcHandler := RPCService{aggregator: a}
	apiHandler := registerRoutes(a)
	service, err := si.StartService("Aggregator", rpcHandler, apiHandler)
	if err != nil {
		return nil, err
	}

	a.service = service
	return a, nil
}
