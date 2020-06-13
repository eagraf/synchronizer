package dataserver

import (
	"github.com/eagraf/synchronizer/messenger"
	"github.com/eagraf/synchronizer/service"
)

// TODO move these into messenger package
const (
	MessageInitiateDataTransfer string = "InitiateDataTransfer"
)

// DataServer type contains state for this service
type DataServer struct {
	workers    map[string]*Worker // TODO maybe use a sorted tree based structure -> when workers are transferred to coordinator, coord. just has to merge workers together from each selector
	messenger  *messenger.Messenger
	rpcHandler RPCService
	service    *service.Service
}

// Worker models a connected worker
type Worker struct {
	UUID string
}

// NewDataServer creates a new data server object
func NewDataServer(si service.ServiceInitiator) (*DataServer, error) {
	ds := &DataServer{
		workers:   make(map[string]*Worker),
		messenger: messenger.NewMessenger(),
	}
	// Setup service
	rpcHandler := RPCService{dataServer: ds}
	apiHandler := registerRoutes(ds)
	service, err := si.StartService("Data_Server", rpcHandler, apiHandler)
	if err != nil {
		return nil, err
	}

	ds.service = service
	return ds, nil
}
