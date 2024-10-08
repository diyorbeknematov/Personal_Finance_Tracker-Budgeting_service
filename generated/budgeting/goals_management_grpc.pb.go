// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.12.4
// source: budgeting_service/goals_management.proto

package budgeting

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	GoalsManagemenService_CreateGoal_FullMethodName = "/goals_management.GoalsManagemenService/CreateGoal"
	GoalsManagemenService_GetGoals_FullMethodName   = "/goals_management.GoalsManagemenService/GetGoals"
	GoalsManagemenService_GetGoal_FullMethodName    = "/goals_management.GoalsManagemenService/GetGoal"
	GoalsManagemenService_UpdateGoal_FullMethodName = "/goals_management.GoalsManagemenService/UpdateGoal"
	GoalsManagemenService_DeleteGoal_FullMethodName = "/goals_management.GoalsManagemenService/DeleteGoal"
)

// GoalsManagemenServiceClient is the client API for GoalsManagemenService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GoalsManagemenServiceClient interface {
	// Moliyaviy maqsadlarni boshqarish
	CreateGoal(ctx context.Context, in *CreateGoalReq, opts ...grpc.CallOption) (*CreateGoalResp, error)
	GetGoals(ctx context.Context, in *GetGoalsReq, opts ...grpc.CallOption) (*GetGoalsResp, error)
	GetGoal(ctx context.Context, in *GetGoalReq, opts ...grpc.CallOption) (*GetGoalResp, error)
	UpdateGoal(ctx context.Context, in *UpdateGoalReq, opts ...grpc.CallOption) (*UpdateGoalResp, error)
	DeleteGoal(ctx context.Context, in *DeleteGoalReq, opts ...grpc.CallOption) (*DeleteGoalResp, error)
}

type goalsManagemenServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGoalsManagemenServiceClient(cc grpc.ClientConnInterface) GoalsManagemenServiceClient {
	return &goalsManagemenServiceClient{cc}
}

func (c *goalsManagemenServiceClient) CreateGoal(ctx context.Context, in *CreateGoalReq, opts ...grpc.CallOption) (*CreateGoalResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateGoalResp)
	err := c.cc.Invoke(ctx, GoalsManagemenService_CreateGoal_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goalsManagemenServiceClient) GetGoals(ctx context.Context, in *GetGoalsReq, opts ...grpc.CallOption) (*GetGoalsResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetGoalsResp)
	err := c.cc.Invoke(ctx, GoalsManagemenService_GetGoals_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goalsManagemenServiceClient) GetGoal(ctx context.Context, in *GetGoalReq, opts ...grpc.CallOption) (*GetGoalResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetGoalResp)
	err := c.cc.Invoke(ctx, GoalsManagemenService_GetGoal_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goalsManagemenServiceClient) UpdateGoal(ctx context.Context, in *UpdateGoalReq, opts ...grpc.CallOption) (*UpdateGoalResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateGoalResp)
	err := c.cc.Invoke(ctx, GoalsManagemenService_UpdateGoal_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goalsManagemenServiceClient) DeleteGoal(ctx context.Context, in *DeleteGoalReq, opts ...grpc.CallOption) (*DeleteGoalResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteGoalResp)
	err := c.cc.Invoke(ctx, GoalsManagemenService_DeleteGoal_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GoalsManagemenServiceServer is the server API for GoalsManagemenService service.
// All implementations must embed UnimplementedGoalsManagemenServiceServer
// for forward compatibility
type GoalsManagemenServiceServer interface {
	// Moliyaviy maqsadlarni boshqarish
	CreateGoal(context.Context, *CreateGoalReq) (*CreateGoalResp, error)
	GetGoals(context.Context, *GetGoalsReq) (*GetGoalsResp, error)
	GetGoal(context.Context, *GetGoalReq) (*GetGoalResp, error)
	UpdateGoal(context.Context, *UpdateGoalReq) (*UpdateGoalResp, error)
	DeleteGoal(context.Context, *DeleteGoalReq) (*DeleteGoalResp, error)
	mustEmbedUnimplementedGoalsManagemenServiceServer()
}

// UnimplementedGoalsManagemenServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGoalsManagemenServiceServer struct {
}

func (UnimplementedGoalsManagemenServiceServer) CreateGoal(context.Context, *CreateGoalReq) (*CreateGoalResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGoal not implemented")
}
func (UnimplementedGoalsManagemenServiceServer) GetGoals(context.Context, *GetGoalsReq) (*GetGoalsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGoals not implemented")
}
func (UnimplementedGoalsManagemenServiceServer) GetGoal(context.Context, *GetGoalReq) (*GetGoalResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGoal not implemented")
}
func (UnimplementedGoalsManagemenServiceServer) UpdateGoal(context.Context, *UpdateGoalReq) (*UpdateGoalResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateGoal not implemented")
}
func (UnimplementedGoalsManagemenServiceServer) DeleteGoal(context.Context, *DeleteGoalReq) (*DeleteGoalResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteGoal not implemented")
}
func (UnimplementedGoalsManagemenServiceServer) mustEmbedUnimplementedGoalsManagemenServiceServer() {}

// UnsafeGoalsManagemenServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GoalsManagemenServiceServer will
// result in compilation errors.
type UnsafeGoalsManagemenServiceServer interface {
	mustEmbedUnimplementedGoalsManagemenServiceServer()
}

func RegisterGoalsManagemenServiceServer(s grpc.ServiceRegistrar, srv GoalsManagemenServiceServer) {
	s.RegisterService(&GoalsManagemenService_ServiceDesc, srv)
}

func _GoalsManagemenService_CreateGoal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGoalReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoalsManagemenServiceServer).CreateGoal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GoalsManagemenService_CreateGoal_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoalsManagemenServiceServer).CreateGoal(ctx, req.(*CreateGoalReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoalsManagemenService_GetGoals_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGoalsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoalsManagemenServiceServer).GetGoals(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GoalsManagemenService_GetGoals_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoalsManagemenServiceServer).GetGoals(ctx, req.(*GetGoalsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoalsManagemenService_GetGoal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGoalReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoalsManagemenServiceServer).GetGoal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GoalsManagemenService_GetGoal_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoalsManagemenServiceServer).GetGoal(ctx, req.(*GetGoalReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoalsManagemenService_UpdateGoal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateGoalReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoalsManagemenServiceServer).UpdateGoal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GoalsManagemenService_UpdateGoal_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoalsManagemenServiceServer).UpdateGoal(ctx, req.(*UpdateGoalReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoalsManagemenService_DeleteGoal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteGoalReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoalsManagemenServiceServer).DeleteGoal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GoalsManagemenService_DeleteGoal_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoalsManagemenServiceServer).DeleteGoal(ctx, req.(*DeleteGoalReq))
	}
	return interceptor(ctx, in, info, handler)
}

// GoalsManagemenService_ServiceDesc is the grpc.ServiceDesc for GoalsManagemenService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GoalsManagemenService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "goals_management.GoalsManagemenService",
	HandlerType: (*GoalsManagemenServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateGoal",
			Handler:    _GoalsManagemenService_CreateGoal_Handler,
		},
		{
			MethodName: "GetGoals",
			Handler:    _GoalsManagemenService_GetGoals_Handler,
		},
		{
			MethodName: "GetGoal",
			Handler:    _GoalsManagemenService_GetGoal_Handler,
		},
		{
			MethodName: "UpdateGoal",
			Handler:    _GoalsManagemenService_UpdateGoal_Handler,
		},
		{
			MethodName: "DeleteGoal",
			Handler:    _GoalsManagemenService_DeleteGoal_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "budgeting_service/goals_management.proto",
}
