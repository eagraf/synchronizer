package aggregator

import (
	"context"

	"github.com/eagraf/synchronizer/service"
)

// RPCService for Aggregator
type RPCService struct {
	aggregator *Aggregator
}

// ReceiveSchedule handles RPC call from coordinator giving the aggregator a schedule
func (RPCService) ReceiveSchedule(ctx context.Context, req *service.AggregatorReceiveScheduleRequest) (*service.AggregatorReceiveScheduleResponse, error) {
	return &service.AggregatorReceiveScheduleResponse{}, nil
}
