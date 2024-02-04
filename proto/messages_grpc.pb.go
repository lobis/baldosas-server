// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.2
// source: messages.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	PositionsService_GetPositions_FullMethodName = "/main.PositionsService/GetPositions"
)

// PositionsServiceClient is the client API for PositionsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PositionsServiceClient interface {
	GetPositions(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Positions, error)
}

type positionsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPositionsServiceClient(cc grpc.ClientConnInterface) PositionsServiceClient {
	return &positionsServiceClient{cc}
}

func (c *positionsServiceClient) GetPositions(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Positions, error) {
	out := new(Positions)
	err := c.cc.Invoke(ctx, PositionsService_GetPositions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PositionsServiceServer is the server API for PositionsService service.
// All implementations must embed UnimplementedPositionsServiceServer
// for forward compatibility
type PositionsServiceServer interface {
	GetPositions(context.Context, *Empty) (*Positions, error)
	mustEmbedUnimplementedPositionsServiceServer()
}

// UnimplementedPositionsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPositionsServiceServer struct {
}

func (UnimplementedPositionsServiceServer) GetPositions(context.Context, *Empty) (*Positions, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPositions not implemented")
}
func (UnimplementedPositionsServiceServer) mustEmbedUnimplementedPositionsServiceServer() {}

// UnsafePositionsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PositionsServiceServer will
// result in compilation errors.
type UnsafePositionsServiceServer interface {
	mustEmbedUnimplementedPositionsServiceServer()
}

func RegisterPositionsServiceServer(s grpc.ServiceRegistrar, srv PositionsServiceServer) {
	s.RegisterService(&PositionsService_ServiceDesc, srv)
}

func _PositionsService_GetPositions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PositionsServiceServer).GetPositions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PositionsService_GetPositions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PositionsServiceServer).GetPositions(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// PositionsService_ServiceDesc is the grpc.ServiceDesc for PositionsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PositionsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.PositionsService",
	HandlerType: (*PositionsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPositions",
			Handler:    _PositionsService_GetPositions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "messages.proto",
}

const (
	SensorService_GetSensorStatusUpdates_FullMethodName = "/main.SensorService/GetSensorStatusUpdates"
)

// SensorServiceClient is the client API for SensorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SensorServiceClient interface {
	GetSensorStatusUpdates(ctx context.Context, in *Empty, opts ...grpc.CallOption) (SensorService_GetSensorStatusUpdatesClient, error)
}

type sensorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSensorServiceClient(cc grpc.ClientConnInterface) SensorServiceClient {
	return &sensorServiceClient{cc}
}

func (c *sensorServiceClient) GetSensorStatusUpdates(ctx context.Context, in *Empty, opts ...grpc.CallOption) (SensorService_GetSensorStatusUpdatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &SensorService_ServiceDesc.Streams[0], SensorService_GetSensorStatusUpdates_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &sensorServiceGetSensorStatusUpdatesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SensorService_GetSensorStatusUpdatesClient interface {
	Recv() (*SensorStatus, error)
	grpc.ClientStream
}

type sensorServiceGetSensorStatusUpdatesClient struct {
	grpc.ClientStream
}

func (x *sensorServiceGetSensorStatusUpdatesClient) Recv() (*SensorStatus, error) {
	m := new(SensorStatus)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SensorServiceServer is the server API for SensorService service.
// All implementations must embed UnimplementedSensorServiceServer
// for forward compatibility
type SensorServiceServer interface {
	GetSensorStatusUpdates(*Empty, SensorService_GetSensorStatusUpdatesServer) error
	mustEmbedUnimplementedSensorServiceServer()
}

// UnimplementedSensorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSensorServiceServer struct {
}

func (UnimplementedSensorServiceServer) GetSensorStatusUpdates(*Empty, SensorService_GetSensorStatusUpdatesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetSensorStatusUpdates not implemented")
}
func (UnimplementedSensorServiceServer) mustEmbedUnimplementedSensorServiceServer() {}

// UnsafeSensorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SensorServiceServer will
// result in compilation errors.
type UnsafeSensorServiceServer interface {
	mustEmbedUnimplementedSensorServiceServer()
}

func RegisterSensorServiceServer(s grpc.ServiceRegistrar, srv SensorServiceServer) {
	s.RegisterService(&SensorService_ServiceDesc, srv)
}

func _SensorService_GetSensorStatusUpdates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SensorServiceServer).GetSensorStatusUpdates(m, &sensorServiceGetSensorStatusUpdatesServer{stream})
}

type SensorService_GetSensorStatusUpdatesServer interface {
	Send(*SensorStatus) error
	grpc.ServerStream
}

type sensorServiceGetSensorStatusUpdatesServer struct {
	grpc.ServerStream
}

func (x *sensorServiceGetSensorStatusUpdatesServer) Send(m *SensorStatus) error {
	return x.ServerStream.SendMsg(m)
}

// SensorService_ServiceDesc is the grpc.ServiceDesc for SensorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SensorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.SensorService",
	HandlerType: (*SensorServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetSensorStatusUpdates",
			Handler:       _SensorService_GetSensorStatusUpdates_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "messages.proto",
}

const (
	LightService_GetLightStatusUpdates_FullMethodName = "/main.LightService/GetLightStatusUpdates"
)

// LightServiceClient is the client API for LightService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LightServiceClient interface {
	GetLightStatusUpdates(ctx context.Context, in *Empty, opts ...grpc.CallOption) (LightService_GetLightStatusUpdatesClient, error)
}

type lightServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLightServiceClient(cc grpc.ClientConnInterface) LightServiceClient {
	return &lightServiceClient{cc}
}

func (c *lightServiceClient) GetLightStatusUpdates(ctx context.Context, in *Empty, opts ...grpc.CallOption) (LightService_GetLightStatusUpdatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &LightService_ServiceDesc.Streams[0], LightService_GetLightStatusUpdates_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &lightServiceGetLightStatusUpdatesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type LightService_GetLightStatusUpdatesClient interface {
	Recv() (*LightStatus, error)
	grpc.ClientStream
}

type lightServiceGetLightStatusUpdatesClient struct {
	grpc.ClientStream
}

func (x *lightServiceGetLightStatusUpdatesClient) Recv() (*LightStatus, error) {
	m := new(LightStatus)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// LightServiceServer is the server API for LightService service.
// All implementations must embed UnimplementedLightServiceServer
// for forward compatibility
type LightServiceServer interface {
	GetLightStatusUpdates(*Empty, LightService_GetLightStatusUpdatesServer) error
	mustEmbedUnimplementedLightServiceServer()
}

// UnimplementedLightServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLightServiceServer struct {
}

func (UnimplementedLightServiceServer) GetLightStatusUpdates(*Empty, LightService_GetLightStatusUpdatesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetLightStatusUpdates not implemented")
}
func (UnimplementedLightServiceServer) mustEmbedUnimplementedLightServiceServer() {}

// UnsafeLightServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LightServiceServer will
// result in compilation errors.
type UnsafeLightServiceServer interface {
	mustEmbedUnimplementedLightServiceServer()
}

func RegisterLightServiceServer(s grpc.ServiceRegistrar, srv LightServiceServer) {
	s.RegisterService(&LightService_ServiceDesc, srv)
}

func _LightService_GetLightStatusUpdates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(LightServiceServer).GetLightStatusUpdates(m, &lightServiceGetLightStatusUpdatesServer{stream})
}

type LightService_GetLightStatusUpdatesServer interface {
	Send(*LightStatus) error
	grpc.ServerStream
}

type lightServiceGetLightStatusUpdatesServer struct {
	grpc.ServerStream
}

func (x *lightServiceGetLightStatusUpdatesServer) Send(m *LightStatus) error {
	return x.ServerStream.SendMsg(m)
}

// LightService_ServiceDesc is the grpc.ServiceDesc for LightService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LightService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.LightService",
	HandlerType: (*LightServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetLightStatusUpdates",
			Handler:       _LightService_GetLightStatusUpdates_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "messages.proto",
}

const (
	SetLightsService_SetLights_FullMethodName = "/main.SetLightsService/SetLights"
)

// SetLightsServiceClient is the client API for SetLightsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SetLightsServiceClient interface {
	SetLights(ctx context.Context, in *LightsStatus, opts ...grpc.CallOption) (*Empty, error)
}

type setLightsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSetLightsServiceClient(cc grpc.ClientConnInterface) SetLightsServiceClient {
	return &setLightsServiceClient{cc}
}

func (c *setLightsServiceClient) SetLights(ctx context.Context, in *LightsStatus, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, SetLightsService_SetLights_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SetLightsServiceServer is the server API for SetLightsService service.
// All implementations must embed UnimplementedSetLightsServiceServer
// for forward compatibility
type SetLightsServiceServer interface {
	SetLights(context.Context, *LightsStatus) (*Empty, error)
	mustEmbedUnimplementedSetLightsServiceServer()
}

// UnimplementedSetLightsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSetLightsServiceServer struct {
}

func (UnimplementedSetLightsServiceServer) SetLights(context.Context, *LightsStatus) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetLights not implemented")
}
func (UnimplementedSetLightsServiceServer) mustEmbedUnimplementedSetLightsServiceServer() {}

// UnsafeSetLightsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SetLightsServiceServer will
// result in compilation errors.
type UnsafeSetLightsServiceServer interface {
	mustEmbedUnimplementedSetLightsServiceServer()
}

func RegisterSetLightsServiceServer(s grpc.ServiceRegistrar, srv SetLightsServiceServer) {
	s.RegisterService(&SetLightsService_ServiceDesc, srv)
}

func _SetLightsService_SetLights_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LightsStatus)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SetLightsServiceServer).SetLights(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SetLightsService_SetLights_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SetLightsServiceServer).SetLights(ctx, req.(*LightsStatus))
	}
	return interceptor(ctx, in, info, handler)
}

// SetLightsService_ServiceDesc is the grpc.ServiceDesc for SetLightsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SetLightsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.SetLightsService",
	HandlerType: (*SetLightsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetLights",
			Handler:    _SetLightsService_SetLights_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "messages.proto",
}

const (
	SetLightsStreamService_SetLightsStream_FullMethodName = "/main.SetLightsStreamService/SetLightsStream"
)

// SetLightsStreamServiceClient is the client API for SetLightsStreamService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SetLightsStreamServiceClient interface {
	SetLightsStream(ctx context.Context, opts ...grpc.CallOption) (SetLightsStreamService_SetLightsStreamClient, error)
}

type setLightsStreamServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSetLightsStreamServiceClient(cc grpc.ClientConnInterface) SetLightsStreamServiceClient {
	return &setLightsStreamServiceClient{cc}
}

func (c *setLightsStreamServiceClient) SetLightsStream(ctx context.Context, opts ...grpc.CallOption) (SetLightsStreamService_SetLightsStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &SetLightsStreamService_ServiceDesc.Streams[0], SetLightsStreamService_SetLightsStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &setLightsStreamServiceSetLightsStreamClient{stream}
	return x, nil
}

type SetLightsStreamService_SetLightsStreamClient interface {
	Send(*LightsStatus) error
	CloseAndRecv() (*Empty, error)
	grpc.ClientStream
}

type setLightsStreamServiceSetLightsStreamClient struct {
	grpc.ClientStream
}

func (x *setLightsStreamServiceSetLightsStreamClient) Send(m *LightsStatus) error {
	return x.ClientStream.SendMsg(m)
}

func (x *setLightsStreamServiceSetLightsStreamClient) CloseAndRecv() (*Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SetLightsStreamServiceServer is the server API for SetLightsStreamService service.
// All implementations must embed UnimplementedSetLightsStreamServiceServer
// for forward compatibility
type SetLightsStreamServiceServer interface {
	SetLightsStream(SetLightsStreamService_SetLightsStreamServer) error
	mustEmbedUnimplementedSetLightsStreamServiceServer()
}

// UnimplementedSetLightsStreamServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSetLightsStreamServiceServer struct {
}

func (UnimplementedSetLightsStreamServiceServer) SetLightsStream(SetLightsStreamService_SetLightsStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method SetLightsStream not implemented")
}
func (UnimplementedSetLightsStreamServiceServer) mustEmbedUnimplementedSetLightsStreamServiceServer() {
}

// UnsafeSetLightsStreamServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SetLightsStreamServiceServer will
// result in compilation errors.
type UnsafeSetLightsStreamServiceServer interface {
	mustEmbedUnimplementedSetLightsStreamServiceServer()
}

func RegisterSetLightsStreamServiceServer(s grpc.ServiceRegistrar, srv SetLightsStreamServiceServer) {
	s.RegisterService(&SetLightsStreamService_ServiceDesc, srv)
}

func _SetLightsStreamService_SetLightsStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SetLightsStreamServiceServer).SetLightsStream(&setLightsStreamServiceSetLightsStreamServer{stream})
}

type SetLightsStreamService_SetLightsStreamServer interface {
	SendAndClose(*Empty) error
	Recv() (*LightsStatus, error)
	grpc.ServerStream
}

type setLightsStreamServiceSetLightsStreamServer struct {
	grpc.ServerStream
}

func (x *setLightsStreamServiceSetLightsStreamServer) SendAndClose(m *Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *setLightsStreamServiceSetLightsStreamServer) Recv() (*LightsStatus, error) {
	m := new(LightsStatus)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SetLightsStreamService_ServiceDesc is the grpc.ServiceDesc for SetLightsStreamService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SetLightsStreamService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.SetLightsStreamService",
	HandlerType: (*SetLightsStreamServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SetLightsStream",
			Handler:       _SetLightsStreamService_SetLightsStream_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "messages.proto",
}
