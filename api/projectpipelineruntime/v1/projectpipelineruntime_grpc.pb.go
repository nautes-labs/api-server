// Copyright 2023 Nautes Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: projectpipelineruntime/v1/projectpipelineruntime.proto

package v1

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
	ProjectPipelineRuntime_GetProjectPipelineRuntime_FullMethodName    = "/api.projectpipelineruntime.v1.ProjectPipelineRuntime/GetProjectPipelineRuntime"
	ProjectPipelineRuntime_ListProjectPipelineRuntimes_FullMethodName  = "/api.projectpipelineruntime.v1.ProjectPipelineRuntime/ListProjectPipelineRuntimes"
	ProjectPipelineRuntime_SaveProjectPipelineRuntime_FullMethodName   = "/api.projectpipelineruntime.v1.ProjectPipelineRuntime/SaveProjectPipelineRuntime"
	ProjectPipelineRuntime_DeleteProjectPipelineRuntime_FullMethodName = "/api.projectpipelineruntime.v1.ProjectPipelineRuntime/DeleteProjectPipelineRuntime"
)

// ProjectPipelineRuntimeClient is the client API for ProjectPipelineRuntime service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProjectPipelineRuntimeClient interface {
	GetProjectPipelineRuntime(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error)
	ListProjectPipelineRuntimes(ctx context.Context, in *ListsRequest, opts ...grpc.CallOption) (*ListsReply, error)
	SaveProjectPipelineRuntime(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveReply, error)
	DeleteProjectPipelineRuntime(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error)
}

type projectPipelineRuntimeClient struct {
	cc grpc.ClientConnInterface
}

func NewProjectPipelineRuntimeClient(cc grpc.ClientConnInterface) ProjectPipelineRuntimeClient {
	return &projectPipelineRuntimeClient{cc}
}

func (c *projectPipelineRuntimeClient) GetProjectPipelineRuntime(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error) {
	out := new(GetReply)
	err := c.cc.Invoke(ctx, ProjectPipelineRuntime_GetProjectPipelineRuntime_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectPipelineRuntimeClient) ListProjectPipelineRuntimes(ctx context.Context, in *ListsRequest, opts ...grpc.CallOption) (*ListsReply, error) {
	out := new(ListsReply)
	err := c.cc.Invoke(ctx, ProjectPipelineRuntime_ListProjectPipelineRuntimes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectPipelineRuntimeClient) SaveProjectPipelineRuntime(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveReply, error) {
	out := new(SaveReply)
	err := c.cc.Invoke(ctx, ProjectPipelineRuntime_SaveProjectPipelineRuntime_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectPipelineRuntimeClient) DeleteProjectPipelineRuntime(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error) {
	out := new(DeleteReply)
	err := c.cc.Invoke(ctx, ProjectPipelineRuntime_DeleteProjectPipelineRuntime_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProjectPipelineRuntimeServer is the server API for ProjectPipelineRuntime service.
// All implementations must embed UnimplementedProjectPipelineRuntimeServer
// for forward compatibility
type ProjectPipelineRuntimeServer interface {
	GetProjectPipelineRuntime(context.Context, *GetRequest) (*GetReply, error)
	ListProjectPipelineRuntimes(context.Context, *ListsRequest) (*ListsReply, error)
	SaveProjectPipelineRuntime(context.Context, *SaveRequest) (*SaveReply, error)
	DeleteProjectPipelineRuntime(context.Context, *DeleteRequest) (*DeleteReply, error)
	mustEmbedUnimplementedProjectPipelineRuntimeServer()
}

// UnimplementedProjectPipelineRuntimeServer must be embedded to have forward compatible implementations.
type UnimplementedProjectPipelineRuntimeServer struct {
}

func (UnimplementedProjectPipelineRuntimeServer) GetProjectPipelineRuntime(context.Context, *GetRequest) (*GetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProjectPipelineRuntime not implemented")
}
func (UnimplementedProjectPipelineRuntimeServer) ListProjectPipelineRuntimes(context.Context, *ListsRequest) (*ListsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProjectPipelineRuntimes not implemented")
}
func (UnimplementedProjectPipelineRuntimeServer) SaveProjectPipelineRuntime(context.Context, *SaveRequest) (*SaveReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveProjectPipelineRuntime not implemented")
}
func (UnimplementedProjectPipelineRuntimeServer) DeleteProjectPipelineRuntime(context.Context, *DeleteRequest) (*DeleteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProjectPipelineRuntime not implemented")
}
func (UnimplementedProjectPipelineRuntimeServer) mustEmbedUnimplementedProjectPipelineRuntimeServer() {
}

// UnsafeProjectPipelineRuntimeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProjectPipelineRuntimeServer will
// result in compilation errors.
type UnsafeProjectPipelineRuntimeServer interface {
	mustEmbedUnimplementedProjectPipelineRuntimeServer()
}

func RegisterProjectPipelineRuntimeServer(s grpc.ServiceRegistrar, srv ProjectPipelineRuntimeServer) {
	s.RegisterService(&ProjectPipelineRuntime_ServiceDesc, srv)
}

func _ProjectPipelineRuntime_GetProjectPipelineRuntime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectPipelineRuntimeServer).GetProjectPipelineRuntime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectPipelineRuntime_GetProjectPipelineRuntime_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectPipelineRuntimeServer).GetProjectPipelineRuntime(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectPipelineRuntime_ListProjectPipelineRuntimes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectPipelineRuntimeServer).ListProjectPipelineRuntimes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectPipelineRuntime_ListProjectPipelineRuntimes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectPipelineRuntimeServer).ListProjectPipelineRuntimes(ctx, req.(*ListsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectPipelineRuntime_SaveProjectPipelineRuntime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectPipelineRuntimeServer).SaveProjectPipelineRuntime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectPipelineRuntime_SaveProjectPipelineRuntime_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectPipelineRuntimeServer).SaveProjectPipelineRuntime(ctx, req.(*SaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectPipelineRuntime_DeleteProjectPipelineRuntime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectPipelineRuntimeServer).DeleteProjectPipelineRuntime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectPipelineRuntime_DeleteProjectPipelineRuntime_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectPipelineRuntimeServer).DeleteProjectPipelineRuntime(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProjectPipelineRuntime_ServiceDesc is the grpc.ServiceDesc for ProjectPipelineRuntime service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProjectPipelineRuntime_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.projectpipelineruntime.v1.ProjectPipelineRuntime",
	HandlerType: (*ProjectPipelineRuntimeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetProjectPipelineRuntime",
			Handler:    _ProjectPipelineRuntime_GetProjectPipelineRuntime_Handler,
		},
		{
			MethodName: "ListProjectPipelineRuntimes",
			Handler:    _ProjectPipelineRuntime_ListProjectPipelineRuntimes_Handler,
		},
		{
			MethodName: "SaveProjectPipelineRuntime",
			Handler:    _ProjectPipelineRuntime_SaveProjectPipelineRuntime_Handler,
		},
		{
			MethodName: "DeleteProjectPipelineRuntime",
			Handler:    _ProjectPipelineRuntime_DeleteProjectPipelineRuntime_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "projectpipelineruntime/v1/projectpipelineruntime.proto",
}
