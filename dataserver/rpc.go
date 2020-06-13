package dataserver

import (
	"context"

	"github.com/eagraf/synchronizer/service"
)

// RPCService for DataServer
type RPCService struct {
	dataServer *DataServer
}

// DataServerReceiveSchedule handles RPC call from coordinator giving the dataserver a schedule
func (RPCService) DataServerReceiveSchedule(ctx context.Context, req *service.DataServerReceiveScheduleRequest) (*service.DataServerReceiveScheduleResponse, error) {
	//return nil, status.Errorf(codes.Unimplemented, "method ReceiveSchedule not implemented")
	return &service.DataServerReceiveScheduleResponse{}, nil
}
