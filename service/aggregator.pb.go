// Code generated by protoc-gen-go. DO NOT EDIT.
// source: aggregator.proto

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
type AggregatorReceiveScheduleRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AggregatorReceiveScheduleRequest) Reset()         { *m = AggregatorReceiveScheduleRequest{} }
func (m *AggregatorReceiveScheduleRequest) String() string { return proto.CompactTextString(m) }
func (*AggregatorReceiveScheduleRequest) ProtoMessage()    {}
func (*AggregatorReceiveScheduleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_60785b04c84bec7e, []int{0}
}

func (m *AggregatorReceiveScheduleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AggregatorReceiveScheduleRequest.Unmarshal(m, b)
}
func (m *AggregatorReceiveScheduleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AggregatorReceiveScheduleRequest.Marshal(b, m, deterministic)
}
func (m *AggregatorReceiveScheduleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AggregatorReceiveScheduleRequest.Merge(m, src)
}
func (m *AggregatorReceiveScheduleRequest) XXX_Size() int {
	return xxx_messageInfo_AggregatorReceiveScheduleRequest.Size(m)
}
func (m *AggregatorReceiveScheduleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AggregatorReceiveScheduleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AggregatorReceiveScheduleRequest proto.InternalMessageInfo

type AggregatorReceiveScheduleRequest_Schedule struct {
	Jobs                 []*AggregatorReceiveScheduleRequest_Schedule_Job    `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty"`
	Workers              []*AggregatorReceiveScheduleRequest_Schedule_Worker `protobuf:"bytes,2,rep,name=workers,proto3" json:"workers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                            `json:"-"`
	XXX_unrecognized     []byte                                              `json:"-"`
	XXX_sizecache        int32                                               `json:"-"`
}

func (m *AggregatorReceiveScheduleRequest_Schedule) Reset() {
	*m = AggregatorReceiveScheduleRequest_Schedule{}
}
func (m *AggregatorReceiveScheduleRequest_Schedule) String() string {
	return proto.CompactTextString(m)
}
func (*AggregatorReceiveScheduleRequest_Schedule) ProtoMessage() {}
func (*AggregatorReceiveScheduleRequest_Schedule) Descriptor() ([]byte, []int) {
	return fileDescriptor_60785b04c84bec7e, []int{0, 0}
}

func (m *AggregatorReceiveScheduleRequest_Schedule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule.Unmarshal(m, b)
}
func (m *AggregatorReceiveScheduleRequest_Schedule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule.Marshal(b, m, deterministic)
}
func (m *AggregatorReceiveScheduleRequest_Schedule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule.Merge(m, src)
}
func (m *AggregatorReceiveScheduleRequest_Schedule) XXX_Size() int {
	return xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule.Size(m)
}
func (m *AggregatorReceiveScheduleRequest_Schedule) XXX_DiscardUnknown() {
	xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule.DiscardUnknown(m)
}

var xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule proto.InternalMessageInfo

func (m *AggregatorReceiveScheduleRequest_Schedule) GetJobs() []*AggregatorReceiveScheduleRequest_Schedule_Job {
	if m != nil {
		return m.Jobs
	}
	return nil
}

func (m *AggregatorReceiveScheduleRequest_Schedule) GetWorkers() []*AggregatorReceiveScheduleRequest_Schedule_Worker {
	if m != nil {
		return m.Workers
	}
	return nil
}

type AggregatorReceiveScheduleRequest_Schedule_Job struct {
	JobUUID              string   `protobuf:"bytes,1,opt,name=JobUUID,proto3" json:"JobUUID,omitempty"`
	JobType              string   `protobuf:"bytes,2,opt,name=JobType,proto3" json:"JobType,omitempty"`
	TaskSize             int32    `protobuf:"varint,3,opt,name=TaskSize,proto3" json:"TaskSize,omitempty"`
	TaskNumber           int32    `protobuf:"varint,4,opt,name=TaskNumber,proto3" json:"TaskNumber,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AggregatorReceiveScheduleRequest_Schedule_Job) Reset() {
	*m = AggregatorReceiveScheduleRequest_Schedule_Job{}
}
func (m *AggregatorReceiveScheduleRequest_Schedule_Job) String() string {
	return proto.CompactTextString(m)
}
func (*AggregatorReceiveScheduleRequest_Schedule_Job) ProtoMessage() {}
func (*AggregatorReceiveScheduleRequest_Schedule_Job) Descriptor() ([]byte, []int) {
	return fileDescriptor_60785b04c84bec7e, []int{0, 0, 0}
}

func (m *AggregatorReceiveScheduleRequest_Schedule_Job) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Job.Unmarshal(m, b)
}
func (m *AggregatorReceiveScheduleRequest_Schedule_Job) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Job.Marshal(b, m, deterministic)
}
func (m *AggregatorReceiveScheduleRequest_Schedule_Job) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Job.Merge(m, src)
}
func (m *AggregatorReceiveScheduleRequest_Schedule_Job) XXX_Size() int {
	return xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Job.Size(m)
}
func (m *AggregatorReceiveScheduleRequest_Schedule_Job) XXX_DiscardUnknown() {
	xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Job.DiscardUnknown(m)
}

var xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Job proto.InternalMessageInfo

func (m *AggregatorReceiveScheduleRequest_Schedule_Job) GetJobUUID() string {
	if m != nil {
		return m.JobUUID
	}
	return ""
}

func (m *AggregatorReceiveScheduleRequest_Schedule_Job) GetJobType() string {
	if m != nil {
		return m.JobType
	}
	return ""
}

func (m *AggregatorReceiveScheduleRequest_Schedule_Job) GetTaskSize() int32 {
	if m != nil {
		return m.TaskSize
	}
	return 0
}

func (m *AggregatorReceiveScheduleRequest_Schedule_Job) GetTaskNumber() int32 {
	if m != nil {
		return m.TaskNumber
	}
	return 0
}

type AggregatorReceiveScheduleRequest_Schedule_Worker struct {
	WorkerUUID           string                                                   `protobuf:"bytes,1,opt,name=workerUUID,proto3" json:"workerUUID,omitempty"`
	DeviceType           string                                                   `protobuf:"bytes,2,opt,name=deviceType,proto3" json:"deviceType,omitempty"`
	Tasks                []*AggregatorReceiveScheduleRequest_Schedule_Worker_Task `protobuf:"bytes,3,rep,name=tasks,proto3" json:"tasks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                 `json:"-"`
	XXX_unrecognized     []byte                                                   `json:"-"`
	XXX_sizecache        int32                                                    `json:"-"`
}

func (m *AggregatorReceiveScheduleRequest_Schedule_Worker) Reset() {
	*m = AggregatorReceiveScheduleRequest_Schedule_Worker{}
}
func (m *AggregatorReceiveScheduleRequest_Schedule_Worker) String() string {
	return proto.CompactTextString(m)
}
func (*AggregatorReceiveScheduleRequest_Schedule_Worker) ProtoMessage() {}
func (*AggregatorReceiveScheduleRequest_Schedule_Worker) Descriptor() ([]byte, []int) {
	return fileDescriptor_60785b04c84bec7e, []int{0, 0, 1}
}

func (m *AggregatorReceiveScheduleRequest_Schedule_Worker) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Worker.Unmarshal(m, b)
}
func (m *AggregatorReceiveScheduleRequest_Schedule_Worker) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Worker.Marshal(b, m, deterministic)
}
func (m *AggregatorReceiveScheduleRequest_Schedule_Worker) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Worker.Merge(m, src)
}
func (m *AggregatorReceiveScheduleRequest_Schedule_Worker) XXX_Size() int {
	return xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Worker.Size(m)
}
func (m *AggregatorReceiveScheduleRequest_Schedule_Worker) XXX_DiscardUnknown() {
	xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Worker.DiscardUnknown(m)
}

var xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Worker proto.InternalMessageInfo

func (m *AggregatorReceiveScheduleRequest_Schedule_Worker) GetWorkerUUID() string {
	if m != nil {
		return m.WorkerUUID
	}
	return ""
}

func (m *AggregatorReceiveScheduleRequest_Schedule_Worker) GetDeviceType() string {
	if m != nil {
		return m.DeviceType
	}
	return ""
}

func (m *AggregatorReceiveScheduleRequest_Schedule_Worker) GetTasks() []*AggregatorReceiveScheduleRequest_Schedule_Worker_Task {
	if m != nil {
		return m.Tasks
	}
	return nil
}

type AggregatorReceiveScheduleRequest_Schedule_Worker_Task struct {
	JobUUID              string   `protobuf:"bytes,1,opt,name=JobUUID,proto3" json:"JobUUID,omitempty"`
	TaskIndex            int32    `protobuf:"varint,2,opt,name=TaskIndex,proto3" json:"TaskIndex,omitempty"`
	TaskSize             int32    `protobuf:"varint,3,opt,name=TaskSize,proto3" json:"TaskSize,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AggregatorReceiveScheduleRequest_Schedule_Worker_Task) Reset() {
	*m = AggregatorReceiveScheduleRequest_Schedule_Worker_Task{}
}
func (m *AggregatorReceiveScheduleRequest_Schedule_Worker_Task) String() string {
	return proto.CompactTextString(m)
}
func (*AggregatorReceiveScheduleRequest_Schedule_Worker_Task) ProtoMessage() {}
func (*AggregatorReceiveScheduleRequest_Schedule_Worker_Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_60785b04c84bec7e, []int{0, 0, 1, 0}
}

func (m *AggregatorReceiveScheduleRequest_Schedule_Worker_Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Worker_Task.Unmarshal(m, b)
}
func (m *AggregatorReceiveScheduleRequest_Schedule_Worker_Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Worker_Task.Marshal(b, m, deterministic)
}
func (m *AggregatorReceiveScheduleRequest_Schedule_Worker_Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Worker_Task.Merge(m, src)
}
func (m *AggregatorReceiveScheduleRequest_Schedule_Worker_Task) XXX_Size() int {
	return xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Worker_Task.Size(m)
}
func (m *AggregatorReceiveScheduleRequest_Schedule_Worker_Task) XXX_DiscardUnknown() {
	xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Worker_Task.DiscardUnknown(m)
}

var xxx_messageInfo_AggregatorReceiveScheduleRequest_Schedule_Worker_Task proto.InternalMessageInfo

func (m *AggregatorReceiveScheduleRequest_Schedule_Worker_Task) GetJobUUID() string {
	if m != nil {
		return m.JobUUID
	}
	return ""
}

func (m *AggregatorReceiveScheduleRequest_Schedule_Worker_Task) GetTaskIndex() int32 {
	if m != nil {
		return m.TaskIndex
	}
	return 0
}

func (m *AggregatorReceiveScheduleRequest_Schedule_Worker_Task) GetTaskSize() int32 {
	if m != nil {
		return m.TaskSize
	}
	return 0
}

// Empty response for now
type AggregatorReceiveScheduleResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AggregatorReceiveScheduleResponse) Reset()         { *m = AggregatorReceiveScheduleResponse{} }
func (m *AggregatorReceiveScheduleResponse) String() string { return proto.CompactTextString(m) }
func (*AggregatorReceiveScheduleResponse) ProtoMessage()    {}
func (*AggregatorReceiveScheduleResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_60785b04c84bec7e, []int{1}
}

func (m *AggregatorReceiveScheduleResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AggregatorReceiveScheduleResponse.Unmarshal(m, b)
}
func (m *AggregatorReceiveScheduleResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AggregatorReceiveScheduleResponse.Marshal(b, m, deterministic)
}
func (m *AggregatorReceiveScheduleResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AggregatorReceiveScheduleResponse.Merge(m, src)
}
func (m *AggregatorReceiveScheduleResponse) XXX_Size() int {
	return xxx_messageInfo_AggregatorReceiveScheduleResponse.Size(m)
}
func (m *AggregatorReceiveScheduleResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AggregatorReceiveScheduleResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AggregatorReceiveScheduleResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*AggregatorReceiveScheduleRequest)(nil), "aggregatorservice.AggregatorReceiveScheduleRequest")
	proto.RegisterType((*AggregatorReceiveScheduleRequest_Schedule)(nil), "aggregatorservice.AggregatorReceiveScheduleRequest.Schedule")
	proto.RegisterType((*AggregatorReceiveScheduleRequest_Schedule_Job)(nil), "aggregatorservice.AggregatorReceiveScheduleRequest.Schedule.Job")
	proto.RegisterType((*AggregatorReceiveScheduleRequest_Schedule_Worker)(nil), "aggregatorservice.AggregatorReceiveScheduleRequest.Schedule.Worker")
	proto.RegisterType((*AggregatorReceiveScheduleRequest_Schedule_Worker_Task)(nil), "aggregatorservice.AggregatorReceiveScheduleRequest.Schedule.Worker.Task")
	proto.RegisterType((*AggregatorReceiveScheduleResponse)(nil), "aggregatorservice.AggregatorReceiveScheduleResponse")
}

func init() {
	proto.RegisterFile("aggregator.proto", fileDescriptor_60785b04c84bec7e)
}

var fileDescriptor_60785b04c84bec7e = []byte{
	// 395 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0xcf, 0xae, 0xd2, 0x40,
	0x14, 0xc6, 0x2d, 0x2d, 0xff, 0x8e, 0x0b, 0x65, 0x56, 0x93, 0xc6, 0x98, 0x8a, 0x09, 0x61, 0x35,
	0x24, 0xe0, 0x03, 0x28, 0xba, 0x10, 0x16, 0xc6, 0x0c, 0x10, 0x13, 0x12, 0x4d, 0xda, 0x72, 0x2c,
	0x15, 0xe9, 0xe0, 0x4c, 0x8b, 0xc2, 0xc2, 0x95, 0x4b, 0x5f, 0xe9, 0xbe, 0xcc, 0x7d, 0x92, 0x9b,
	0xe9, 0x1f, 0xda, 0xdc, 0x9b, 0x70, 0x73, 0xff, 0xec, 0xf8, 0xbe, 0x73, 0xbe, 0x1f, 0xdf, 0x9c,
	0xa4, 0xf0, 0xdc, 0x0d, 0x02, 0x89, 0x81, 0x1b, 0x0b, 0xc9, 0x76, 0x52, 0xc4, 0x82, 0x74, 0x4a,
	0x47, 0xa1, 0xdc, 0x87, 0x3e, 0x76, 0x2f, 0x2d, 0x70, 0xde, 0x9d, 0x5c, 0x8e, 0x3e, 0x86, 0x7b,
	0x9c, 0xf9, 0x6b, 0x5c, 0x25, 0x3f, 0x91, 0xe3, 0xaf, 0x04, 0x55, 0x6c, 0x5f, 0x58, 0xd0, 0x2a,
	0x3c, 0x32, 0x07, 0xeb, 0x87, 0xf0, 0x14, 0x35, 0x1c, 0xb3, 0xff, 0x74, 0xf8, 0x96, 0xdd, 0x60,
	0xb2, 0xdb, 0x78, 0xac, 0xd0, 0x6c, 0x2a, 0x3c, 0x9e, 0xd2, 0xc8, 0x57, 0x68, 0xfe, 0x16, 0x72,
	0x83, 0x52, 0xd1, 0x5a, 0x0a, 0x7e, 0xff, 0x20, 0xf0, 0x97, 0x94, 0xc5, 0x0b, 0xa6, 0x9d, 0x80,
	0x39, 0x15, 0x1e, 0xa1, 0xd0, 0x9c, 0x0a, 0x6f, 0xb1, 0x98, 0x7c, 0xa0, 0x86, 0x63, 0xf4, 0xdb,
	0xbc, 0x90, 0xf9, 0x64, 0x7e, 0xd8, 0x21, 0xad, 0x9d, 0x26, 0x5a, 0x12, 0x1b, 0x5a, 0x73, 0x57,
	0x6d, 0x66, 0xe1, 0x11, 0xa9, 0xe9, 0x18, 0xfd, 0x3a, 0x3f, 0x69, 0xf2, 0x12, 0x40, 0xff, 0xfe,
	0x94, 0x6c, 0x3d, 0x94, 0xd4, 0x4a, 0xa7, 0x15, 0xc7, 0xfe, 0x57, 0x83, 0x46, 0x56, 0x45, 0xaf,
	0x66, 0x65, 0x2a, 0xff, 0x5e, 0x71, 0xf4, 0x7c, 0x85, 0xfa, 0x95, 0x95, 0x0e, 0x15, 0x87, 0x7c,
	0x83, 0x7a, 0xec, 0xaa, 0x8d, 0xa2, 0x66, 0x7a, 0x9e, 0x8f, 0x8f, 0x70, 0x1e, 0xa6, 0x9b, 0xf2,
	0x0c, 0x6b, 0x2f, 0xc1, 0xd2, 0xf2, 0xcc, 0x89, 0x5e, 0x40, 0x5b, 0x6f, 0x4c, 0xa2, 0x15, 0xfe,
	0x49, 0x0b, 0xd6, 0x79, 0x69, 0x9c, 0x3b, 0x53, 0xf7, 0x35, 0xbc, 0x3a, 0xd3, 0x4d, 0xed, 0x44,
	0xa4, 0x70, 0xf8, 0xdf, 0x00, 0x28, 0xb7, 0xc8, 0x5f, 0x78, 0x76, 0x6d, 0x93, 0x8c, 0xee, 0xf1,
	0x66, 0xfb, 0xcd, 0xdd, 0x42, 0x59, 0x99, 0xee, 0x93, 0x31, 0x42, 0xcf, 0x17, 0x5b, 0x16, 0x84,
	0xf1, 0x3a, 0xf1, 0x18, 0xba, 0x81, 0x74, 0xbf, 0x33, 0x75, 0x88, 0xfc, 0xb5, 0x14, 0x51, 0x78,
	0x44, 0xc9, 0x72, 0xda, 0xb8, 0x53, 0xe2, 0x66, 0x99, 0xf5, 0xd9, 0x58, 0xf6, 0xf2, 0xa0, 0x2f,
	0xb6, 0x83, 0x2c, 0x3c, 0xa8, 0x86, 0x07, 0x79, 0xd8, 0x6b, 0xa4, 0x5f, 0xe6, 0xe8, 0x2a, 0x00,
	0x00, 0xff, 0xff, 0xba, 0x5c, 0xfa, 0xae, 0xad, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AggregatorClient is the client API for Aggregator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AggregatorClient interface {
	ReceiveSchedule(ctx context.Context, in *AggregatorReceiveScheduleRequest, opts ...grpc.CallOption) (*AggregatorReceiveScheduleResponse, error)
}

type aggregatorClient struct {
	cc grpc.ClientConnInterface
}

func NewAggregatorClient(cc grpc.ClientConnInterface) AggregatorClient {
	return &aggregatorClient{cc}
}

func (c *aggregatorClient) ReceiveSchedule(ctx context.Context, in *AggregatorReceiveScheduleRequest, opts ...grpc.CallOption) (*AggregatorReceiveScheduleResponse, error) {
	out := new(AggregatorReceiveScheduleResponse)
	err := c.cc.Invoke(ctx, "/aggregatorservice.Aggregator/ReceiveSchedule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AggregatorServer is the server API for Aggregator service.
type AggregatorServer interface {
	ReceiveSchedule(context.Context, *AggregatorReceiveScheduleRequest) (*AggregatorReceiveScheduleResponse, error)
}

// UnimplementedAggregatorServer can be embedded to have forward compatible implementations.
type UnimplementedAggregatorServer struct {
}

func (*UnimplementedAggregatorServer) ReceiveSchedule(ctx context.Context, req *AggregatorReceiveScheduleRequest) (*AggregatorReceiveScheduleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReceiveSchedule not implemented")
}

func RegisterAggregatorServer(s *grpc.Server, srv AggregatorServer) {
	s.RegisterService(&_Aggregator_serviceDesc, srv)
}

func _Aggregator_ReceiveSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AggregatorReceiveScheduleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AggregatorServer).ReceiveSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/aggregatorservice.Aggregator/ReceiveSchedule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AggregatorServer).ReceiveSchedule(ctx, req.(*AggregatorReceiveScheduleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Aggregator_serviceDesc = grpc.ServiceDesc{
	ServiceName: "aggregatorservice.Aggregator",
	HandlerType: (*AggregatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReceiveSchedule",
			Handler:    _Aggregator_ReceiveSchedule_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "aggregator.proto",
}
