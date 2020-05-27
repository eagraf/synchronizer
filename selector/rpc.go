package selector

import (
	"context"

	"github.com/eagraf/synchronizer/service"
)

// RPCService is an implementation of service.SelectorServer
type RPCService struct{}

// GetWorkers reutrns a list of available workers to the coordinator
func (rs RPCService) GetWorkers(ctx context.Context, req *service.WorkersRequest) (*service.WorkersResponse, error) {
	return &service.WorkersResponse{}, nil
}

// Handoff performs handoff of a worker to data server/aggregator
func (rs RPCService) Handoff(ctx context.Context, req *service.HandoffRequest) (*service.HandoffResponse, error) {
	return &service.HandoffResponse{}, nil
}
