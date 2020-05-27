package selector

import (
	"context"

	"github.com/eagraf/synchronizer/service"
)

// RPCService is an implementation of service.SelectorServer
type RPCService struct {
	selector *Selector
}

// GetWorkers reutrns a list of available workers to the coordinator
func (rs RPCService) GetWorkers(ctx context.Context, req *service.WorkersRequest) (*service.WorkersResponse, error) {
	// Convert workers into a slice of service.WorkersResponse_Worker
	workers := make([]*service.WorkersResponse_Worker, 0, len(rs.selector.workers))
	for _, w := range rs.selector.workers {
		wrw := &service.WorkersResponse_Worker{
			WorkerUUID: w.UUID,
		}
		workers = append(workers, wrw)
	}
	return &service.WorkersResponse{Workers: workers}, nil
}

// Handoff performs handoff of a worker to data server/aggregator
func (rs RPCService) Handoff(ctx context.Context, req *service.HandoffRequest) (*service.HandoffResponse, error) {
	return &service.HandoffResponse{}, nil
}
