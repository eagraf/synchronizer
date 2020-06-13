package dataserver

import (
	"context"

	"github.com/eagraf/synchronizer/service"
)

// RPCService for DataServer
type RPCService struct {
	dataServer *DataServer
}

// ReceiveSchedule handles RPC call from coordinator giving the dataserver a schedule
func (RPCService) ReceiveSchedule(ctx context.Context, req *service.ReceiveScheduleRequest) (*service.ReceiveScheduleResponse, error) {
	//return nil, status.Errorf(codes.Unimplemented, "method ReceiveSchedule not implemented")
	return &service.ReceiveScheduleResponse{}, nil
}
