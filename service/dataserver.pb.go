// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dataserver.proto

package service

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Contains info on jobs and workers assigned to this data server
type DataServerReceiveScheduleRequest struct {
	Schedule             *DataServerReceiveScheduleRequest_Schedule `protobuf:"bytes,1,opt,name=schedule,proto3" json:"schedule,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                   `json:"-"`
	XXX_unrecognized     []byte                                     `json:"-"`
	XXX_sizecache        int32                                      `json:"-"`
}

func (m *DataServerReceiveScheduleRequest) Reset()         { *m = DataServerReceiveScheduleRequest{} }
func (m *DataServerReceiveScheduleRequest) String() string { return proto.CompactTextString(m) }
func (*DataServerReceiveScheduleRequest) ProtoMessage()    {}
func (*DataServerReceiveScheduleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d248291538510d87, []int{0}
}

func (m *DataServerReceiveScheduleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataServerReceiveScheduleRequest.Unmarshal(m, b)
}
func (m *DataServerReceiveScheduleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataServerReceiveScheduleRequest.Marshal(b, m, deterministic)
}
func (m *DataServerReceiveScheduleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataServerReceiveScheduleRequest.Merge(m, src)
}
func (m *DataServerReceiveScheduleRequest) XXX_Size() int {
	return xxx_messageInfo_DataServerReceiveScheduleRequest.Size(m)
}
func (m *DataServerReceiveScheduleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DataServerReceiveScheduleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DataServerReceiveScheduleRequest proto.InternalMessageInfo

func (m *DataServerReceiveScheduleRequest) GetSchedule() *DataServerReceiveScheduleRequest_Schedule {
	if m != nil {
		return m.Schedule
	}
	return nil
}

type DataServerReceiveScheduleRequest_Schedule struct {
	Jobs                 []*DataServerReceiveScheduleRequest_Schedule_Job    `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty"`
	Workers              []*DataServerReceiveScheduleRequest_Schedule_Worker `protobuf:"bytes,2,rep,name=workers,proto3" json:"workers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                            `json:"-"`
	XXX_unrecognized     []byte                                              `json:"-"`
	XXX_sizecache        int32                                               `json:"-"`
}

func (m *DataServerReceiveScheduleRequest_Schedule) Reset() {
	*m = DataServerReceiveScheduleRequest_Schedule{}
}
func (m *DataServerReceiveScheduleRequest_Schedule) String() string {
	return proto.CompactTextString(m)
}
func (*DataServerReceiveScheduleRequest_Schedule) ProtoMessage() {}
func (*DataServerReceiveScheduleRequest_Schedule) Descriptor() ([]byte, []int) {
	return fileDescriptor_d248291538510d87, []int{0, 0}
}

func (m *DataServerReceiveScheduleRequest_Schedule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule.Unmarshal(m, b)
}
func (m *DataServerReceiveScheduleRequest_Schedule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule.Marshal(b, m, deterministic)
}
func (m *DataServerReceiveScheduleRequest_Schedule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule.Merge(m, src)
}
func (m *DataServerReceiveScheduleRequest_Schedule) XXX_Size() int {
	return xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule.Size(m)
}
func (m *DataServerReceiveScheduleRequest_Schedule) XXX_DiscardUnknown() {
	xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule.DiscardUnknown(m)
}

var xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule proto.InternalMessageInfo

func (m *DataServerReceiveScheduleRequest_Schedule) GetJobs() []*DataServerReceiveScheduleRequest_Schedule_Job {
	if m != nil {
		return m.Jobs
	}
	return nil
}

func (m *DataServerReceiveScheduleRequest_Schedule) GetWorkers() []*DataServerReceiveScheduleRequest_Schedule_Worker {
	if m != nil {
		return m.Workers
	}
	return nil
}

type DataServerReceiveScheduleRequest_Schedule_Job struct {
	JobUUID              string   `protobuf:"bytes,1,opt,name=JobUUID,proto3" json:"JobUUID,omitempty"`
	JobType              string   `protobuf:"bytes,2,opt,name=JobType,proto3" json:"JobType,omitempty"`
	TaskSize             int32    `protobuf:"varint,3,opt,name=TaskSize,proto3" json:"TaskSize,omitempty"`
	TaskNumber           int32    `protobuf:"varint,4,opt,name=TaskNumber,proto3" json:"TaskNumber,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DataServerReceiveScheduleRequest_Schedule_Job) Reset() {
	*m = DataServerReceiveScheduleRequest_Schedule_Job{}
}
func (m *DataServerReceiveScheduleRequest_Schedule_Job) String() string {
	return proto.CompactTextString(m)
}
func (*DataServerReceiveScheduleRequest_Schedule_Job) ProtoMessage() {}
func (*DataServerReceiveScheduleRequest_Schedule_Job) Descriptor() ([]byte, []int) {
	return fileDescriptor_d248291538510d87, []int{0, 0, 0}
}

func (m *DataServerReceiveScheduleRequest_Schedule_Job) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Job.Unmarshal(m, b)
}
func (m *DataServerReceiveScheduleRequest_Schedule_Job) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Job.Marshal(b, m, deterministic)
}
func (m *DataServerReceiveScheduleRequest_Schedule_Job) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Job.Merge(m, src)
}
func (m *DataServerReceiveScheduleRequest_Schedule_Job) XXX_Size() int {
	return xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Job.Size(m)
}
func (m *DataServerReceiveScheduleRequest_Schedule_Job) XXX_DiscardUnknown() {
	xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Job.DiscardUnknown(m)
}

var xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Job proto.InternalMessageInfo

func (m *DataServerReceiveScheduleRequest_Schedule_Job) GetJobUUID() string {
	if m != nil {
		return m.JobUUID
	}
	return ""
}

func (m *DataServerReceiveScheduleRequest_Schedule_Job) GetJobType() string {
	if m != nil {
		return m.JobType
	}
	return ""
}

func (m *DataServerReceiveScheduleRequest_Schedule_Job) GetTaskSize() int32 {
	if m != nil {
		return m.TaskSize
	}
	return 0
}

func (m *DataServerReceiveScheduleRequest_Schedule_Job) GetTaskNumber() int32 {
	if m != nil {
		return m.TaskNumber
	}
	return 0
}

type DataServerReceiveScheduleRequest_Schedule_Worker struct {
	WorkerUUID           string                                                   `protobuf:"bytes,1,opt,name=workerUUID,proto3" json:"workerUUID,omitempty"`
	DeviceType           string                                                   `protobuf:"bytes,2,opt,name=deviceType,proto3" json:"deviceType,omitempty"`
	Tasks                []*DataServerReceiveScheduleRequest_Schedule_Worker_Task `protobuf:"bytes,3,rep,name=tasks,proto3" json:"tasks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                 `json:"-"`
	XXX_unrecognized     []byte                                                   `json:"-"`
	XXX_sizecache        int32                                                    `json:"-"`
}

func (m *DataServerReceiveScheduleRequest_Schedule_Worker) Reset() {
	*m = DataServerReceiveScheduleRequest_Schedule_Worker{}
}
func (m *DataServerReceiveScheduleRequest_Schedule_Worker) String() string {
	return proto.CompactTextString(m)
}
func (*DataServerReceiveScheduleRequest_Schedule_Worker) ProtoMessage() {}
func (*DataServerReceiveScheduleRequest_Schedule_Worker) Descriptor() ([]byte, []int) {
	return fileDescriptor_d248291538510d87, []int{0, 0, 1}
}

func (m *DataServerReceiveScheduleRequest_Schedule_Worker) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Worker.Unmarshal(m, b)
}
func (m *DataServerReceiveScheduleRequest_Schedule_Worker) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Worker.Marshal(b, m, deterministic)
}
func (m *DataServerReceiveScheduleRequest_Schedule_Worker) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Worker.Merge(m, src)
}
func (m *DataServerReceiveScheduleRequest_Schedule_Worker) XXX_Size() int {
	return xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Worker.Size(m)
}
func (m *DataServerReceiveScheduleRequest_Schedule_Worker) XXX_DiscardUnknown() {
	xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Worker.DiscardUnknown(m)
}

var xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Worker proto.InternalMessageInfo

func (m *DataServerReceiveScheduleRequest_Schedule_Worker) GetWorkerUUID() string {
	if m != nil {
		return m.WorkerUUID
	}
	return ""
}

func (m *DataServerReceiveScheduleRequest_Schedule_Worker) GetDeviceType() string {
	if m != nil {
		return m.DeviceType
	}
	return ""
}

func (m *DataServerReceiveScheduleRequest_Schedule_Worker) GetTasks() []*DataServerReceiveScheduleRequest_Schedule_Worker_Task {
	if m != nil {
		return m.Tasks
	}
	return nil
}

type DataServerReceiveScheduleRequest_Schedule_Worker_Task struct {
	JobUUID              string   `protobuf:"bytes,1,opt,name=JobUUID,proto3" json:"JobUUID,omitempty"`
	TaskIndex            int32    `protobuf:"varint,2,opt,name=TaskIndex,proto3" json:"TaskIndex,omitempty"`
	TaskSize             int32    `protobuf:"varint,3,opt,name=TaskSize,proto3" json:"TaskSize,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DataServerReceiveScheduleRequest_Schedule_Worker_Task) Reset() {
	*m = DataServerReceiveScheduleRequest_Schedule_Worker_Task{}
}
func (m *DataServerReceiveScheduleRequest_Schedule_Worker_Task) String() string {
	return proto.CompactTextString(m)
}
func (*DataServerReceiveScheduleRequest_Schedule_Worker_Task) ProtoMessage() {}
func (*DataServerReceiveScheduleRequest_Schedule_Worker_Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_d248291538510d87, []int{0, 0, 1, 0}
}

func (m *DataServerReceiveScheduleRequest_Schedule_Worker_Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Worker_Task.Unmarshal(m, b)
}
func (m *DataServerReceiveScheduleRequest_Schedule_Worker_Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Worker_Task.Marshal(b, m, deterministic)
}
func (m *DataServerReceiveScheduleRequest_Schedule_Worker_Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Worker_Task.Merge(m, src)
}
func (m *DataServerReceiveScheduleRequest_Schedule_Worker_Task) XXX_Size() int {
	return xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Worker_Task.Size(m)
}
func (m *DataServerReceiveScheduleRequest_Schedule_Worker_Task) XXX_DiscardUnknown() {
	xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Worker_Task.DiscardUnknown(m)
}

var xxx_messageInfo_DataServerReceiveScheduleRequest_Schedule_Worker_Task proto.InternalMessageInfo

func (m *DataServerReceiveScheduleRequest_Schedule_Worker_Task) GetJobUUID() string {
	if m != nil {
		return m.JobUUID
	}
	return ""
}

func (m *DataServerReceiveScheduleRequest_Schedule_Worker_Task) GetTaskIndex() int32 {
	if m != nil {
		return m.TaskIndex
	}
	return 0
}

func (m *DataServerReceiveScheduleRequest_Schedule_Worker_Task) GetTaskSize() int32 {
	if m != nil {
		return m.TaskSize
	}
	return 0
}

// Empty response for now
type DataServerReceiveScheduleResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DataServerReceiveScheduleResponse) Reset()         { *m = DataServerReceiveScheduleResponse{} }
func (m *DataServerReceiveScheduleResponse) String() string { return proto.CompactTextString(m) }
func (*DataServerReceiveScheduleResponse) ProtoMessage()    {}
func (*DataServerReceiveScheduleResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d248291538510d87, []int{1}
}

func (m *DataServerReceiveScheduleResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataServerReceiveScheduleResponse.Unmarshal(m, b)
}
func (m *DataServerReceiveScheduleResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataServerReceiveScheduleResponse.Marshal(b, m, deterministic)
}
func (m *DataServerReceiveScheduleResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataServerReceiveScheduleResponse.Merge(m, src)
}
func (m *DataServerReceiveScheduleResponse) XXX_Size() int {
	return xxx_messageInfo_DataServerReceiveScheduleResponse.Size(m)
}
func (m *DataServerReceiveScheduleResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DataServerReceiveScheduleResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DataServerReceiveScheduleResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*DataServerReceiveScheduleRequest)(nil), "dataserverservice.DataServerReceiveScheduleRequest")
	proto.RegisterType((*DataServerReceiveScheduleRequest_Schedule)(nil), "dataserverservice.DataServerReceiveScheduleRequest.Schedule")
	proto.RegisterType((*DataServerReceiveScheduleRequest_Schedule_Job)(nil), "dataserverservice.DataServerReceiveScheduleRequest.Schedule.Job")
	proto.RegisterType((*DataServerReceiveScheduleRequest_Schedule_Worker)(nil), "dataserverservice.DataServerReceiveScheduleRequest.Schedule.Worker")
	proto.RegisterType((*DataServerReceiveScheduleRequest_Schedule_Worker_Task)(nil), "dataserverservice.DataServerReceiveScheduleRequest.Schedule.Worker.Task")
	proto.RegisterType((*DataServerReceiveScheduleResponse)(nil), "dataserverservice.DataServerReceiveScheduleResponse")
}

func init() {
	proto.RegisterFile("dataserver.proto", fileDescriptor_d248291538510d87)
}

var fileDescriptor_d248291538510d87 = []byte{
	// 413 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0xc1, 0x6e, 0xd3, 0x40,
	0x14, 0xc4, 0xb1, 0xdd, 0xa6, 0xaf, 0x17, 0xba, 0xa7, 0xc5, 0x42, 0xc8, 0x04, 0xa9, 0xca, 0x69,
	0x23, 0xa5, 0x1c, 0x39, 0xa0, 0xd2, 0x03, 0xcd, 0x01, 0xa1, 0x75, 0x2a, 0x50, 0x25, 0x90, 0xd6,
	0xf6, 0x23, 0x31, 0x21, 0xde, 0xb0, 0x6b, 0x07, 0x92, 0x33, 0x07, 0xbe, 0x82, 0x2f, 0x81, 0x7f,
	0x43, 0xbb, 0x76, 0x1c, 0x4b, 0x11, 0x41, 0x90, 0x5e, 0xa2, 0xcc, 0xbc, 0x37, 0xb3, 0xf3, 0xe6,
	0x60, 0xb8, 0x9f, 0x8a, 0x42, 0x68, 0x54, 0x4b, 0x54, 0x6c, 0xa1, 0x64, 0x21, 0xc9, 0xd9, 0x96,
	0x31, 0xbf, 0x59, 0x82, 0xbd, 0x9f, 0x3e, 0x84, 0x57, 0xa2, 0x10, 0x91, 0x65, 0x39, 0x26, 0x98,
	0x2d, 0x31, 0x4a, 0xa6, 0x98, 0x96, 0x9f, 0x90, 0xe3, 0xe7, 0x12, 0x75, 0x41, 0xde, 0x42, 0x57,
	0xd7, 0x14, 0x75, 0x42, 0xa7, 0x7f, 0x3a, 0x7c, 0xc6, 0x76, 0xac, 0xd8, 0xdf, 0x6c, 0x58, 0x83,
	0x1b, 0xb7, 0xe0, 0x97, 0x07, 0xdd, 0x0d, 0x4d, 0xc6, 0xe0, 0x7d, 0x94, 0xb1, 0xa6, 0x4e, 0xe8,
	0xf6, 0x4f, 0x87, 0xcf, 0x0f, 0x79, 0x82, 0x8d, 0x64, 0xcc, 0xad, 0x1b, 0x79, 0x07, 0xc7, 0x5f,
	0xa4, 0x9a, 0xa1, 0xd2, 0xb4, 0x63, 0x8d, 0x5f, 0x1c, 0x64, 0xfc, 0xc6, 0x7a, 0xf1, 0x8d, 0x67,
	0x50, 0x82, 0x3b, 0x92, 0x31, 0xa1, 0x70, 0x3c, 0x92, 0xf1, 0xcd, 0xcd, 0xf5, 0x95, 0x6d, 0xe8,
	0x84, 0x6f, 0x60, 0x3d, 0x19, 0xaf, 0x16, 0x48, 0x3b, 0xcd, 0xc4, 0x40, 0x12, 0x40, 0x77, 0x2c,
	0xf4, 0x2c, 0xca, 0xd6, 0x48, 0xdd, 0xd0, 0xe9, 0xfb, 0xbc, 0xc1, 0xe4, 0x11, 0x80, 0xf9, 0xff,
	0xaa, 0x9c, 0xc7, 0xa8, 0xa8, 0x67, 0xa7, 0x2d, 0x26, 0xf8, 0xd6, 0x81, 0xa3, 0x2a, 0x8a, 0x59,
	0xad, 0xc2, 0xb4, 0x5e, 0x6f, 0x31, 0x66, 0x9e, 0xa2, 0xb9, 0xb2, 0x95, 0xa1, 0xc5, 0x90, 0xf7,
	0xe0, 0x17, 0x42, 0xcf, 0x34, 0x75, 0x6d, 0x3d, 0x2f, 0xef, 0xa0, 0x1e, 0x66, 0x92, 0xf2, 0xca,
	0x36, 0xb8, 0x05, 0xcf, 0xc0, 0x3d, 0x15, 0x3d, 0x84, 0x13, 0xb3, 0x71, 0x9d, 0xa7, 0xf8, 0xd5,
	0x06, 0xf4, 0xf9, 0x96, 0xd8, 0x57, 0x53, 0xef, 0x09, 0x3c, 0xde, 0x93, 0x4d, 0x2f, 0x64, 0xae,
	0x71, 0xf8, 0xc3, 0x81, 0xb3, 0xed, 0x56, 0x54, 0xdd, 0x44, 0xbe, 0x3b, 0xf0, 0xe0, 0x8f, 0x5a,
	0x72, 0xf1, 0x1f, 0x2d, 0x04, 0x4f, 0xff, 0x4d, 0x54, 0xc5, 0xeb, 0xdd, 0xbb, 0x44, 0x38, 0x4f,
	0xe4, 0x9c, 0x4d, 0xb2, 0x62, 0x5a, 0xc6, 0x0c, 0xc5, 0x44, 0x89, 0x0f, 0x4c, 0xaf, 0xf2, 0x64,
	0xaa, 0x64, 0x9e, 0xad, 0x51, 0xb1, 0xda, 0xed, 0x72, 0xf7, 0x8e, 0xd7, 0xce, 0xed, 0x79, 0x2d,
	0x4c, 0xe4, 0x7c, 0x50, 0x89, 0x07, 0x6d, 0xf1, 0xa0, 0x16, 0xc7, 0x47, 0xf6, 0x2b, 0x70, 0xf1,
	0x3b, 0x00, 0x00, 0xff, 0xff, 0x3c, 0x65, 0x5e, 0x9b, 0x19, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// DataServerServiceClient is the client API for DataServerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DataServerServiceClient interface {
	DataServerReceiveSchedule(ctx context.Context, in *DataServerReceiveScheduleRequest, opts ...grpc.CallOption) (*DataServerReceiveScheduleResponse, error)
}

type dataServerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDataServerServiceClient(cc grpc.ClientConnInterface) DataServerServiceClient {
	return &dataServerServiceClient{cc}
}

func (c *dataServerServiceClient) DataServerReceiveSchedule(ctx context.Context, in *DataServerReceiveScheduleRequest, opts ...grpc.CallOption) (*DataServerReceiveScheduleResponse, error) {
	out := new(DataServerReceiveScheduleResponse)
	err := c.cc.Invoke(ctx, "/dataserverservice.DataServerService/DataServerReceiveSchedule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataServerServiceServer is the server API for DataServerService service.
type DataServerServiceServer interface {
	DataServerReceiveSchedule(context.Context, *DataServerReceiveScheduleRequest) (*DataServerReceiveScheduleResponse, error)
}

// UnimplementedDataServerServiceServer can be embedded to have forward compatible implementations.
type UnimplementedDataServerServiceServer struct {
}

func (*UnimplementedDataServerServiceServer) DataServerReceiveSchedule(ctx context.Context, req *DataServerReceiveScheduleRequest) (*DataServerReceiveScheduleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DataServerReceiveSchedule not implemented")
}

func RegisterDataServerServiceServer(s *grpc.Server, srv DataServerServiceServer) {
	s.RegisterService(&_DataServerService_serviceDesc, srv)
}

func _DataServerService_DataServerReceiveSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DataServerReceiveScheduleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServerServiceServer).DataServerReceiveSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dataserverservice.DataServerService/DataServerReceiveSchedule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServerServiceServer).DataServerReceiveSchedule(ctx, req.(*DataServerReceiveScheduleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DataServerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dataserverservice.DataServerService",
	HandlerType: (*DataServerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DataServerReceiveSchedule",
			Handler:    _DataServerService_DataServerReceiveSchedule_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dataserver.proto",
}