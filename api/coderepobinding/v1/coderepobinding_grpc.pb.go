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
// - protoc             v3.6.1
// source: coderepobinding/v1/coderepobinding.proto

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
	CodeRepoBinding_GetCodeRepoBinding_FullMethodName    = "/api.coderepobinding.v1.CodeRepoBinding/GetCodeRepoBinding"
	CodeRepoBinding_ListCodeRepoBindings_FullMethodName  = "/api.coderepobinding.v1.CodeRepoBinding/ListCodeRepoBindings"
	CodeRepoBinding_SaveCodeRepoBinding_FullMethodName   = "/api.coderepobinding.v1.CodeRepoBinding/SaveCodeRepoBinding"
	CodeRepoBinding_DeleteCodeRepoBinding_FullMethodName = "/api.coderepobinding.v1.CodeRepoBinding/DeleteCodeRepoBinding"
)

// CodeRepoBindingClient is the client API for CodeRepoBinding service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CodeRepoBindingClient interface {
	GetCodeRepoBinding(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error)
	ListCodeRepoBindings(ctx context.Context, in *ListsRequest, opts ...grpc.CallOption) (*ListsReply, error)
	SaveCodeRepoBinding(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveReply, error)
	DeleteCodeRepoBinding(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error)
}

type codeRepoBindingClient struct {
	cc grpc.ClientConnInterface
}

func NewCodeRepoBindingClient(cc grpc.ClientConnInterface) CodeRepoBindingClient {
	return &codeRepoBindingClient{cc}
}

func (c *codeRepoBindingClient) GetCodeRepoBinding(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error) {
	out := new(GetReply)
	err := c.cc.Invoke(ctx, CodeRepoBinding_GetCodeRepoBinding_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *codeRepoBindingClient) ListCodeRepoBindings(ctx context.Context, in *ListsRequest, opts ...grpc.CallOption) (*ListsReply, error) {
	out := new(ListsReply)
	err := c.cc.Invoke(ctx, CodeRepoBinding_ListCodeRepoBindings_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *codeRepoBindingClient) SaveCodeRepoBinding(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveReply, error) {
	out := new(SaveReply)
	err := c.cc.Invoke(ctx, CodeRepoBinding_SaveCodeRepoBinding_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *codeRepoBindingClient) DeleteCodeRepoBinding(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error) {
	out := new(DeleteReply)
	err := c.cc.Invoke(ctx, CodeRepoBinding_DeleteCodeRepoBinding_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CodeRepoBindingServer is the server API for CodeRepoBinding service.
// All implementations must embed UnimplementedCodeRepoBindingServer
// for forward compatibility
type CodeRepoBindingServer interface {
	GetCodeRepoBinding(context.Context, *GetRequest) (*GetReply, error)
	ListCodeRepoBindings(context.Context, *ListsRequest) (*ListsReply, error)
	SaveCodeRepoBinding(context.Context, *SaveRequest) (*SaveReply, error)
	DeleteCodeRepoBinding(context.Context, *DeleteRequest) (*DeleteReply, error)
	mustEmbedUnimplementedCodeRepoBindingServer()
}

// UnimplementedCodeRepoBindingServer must be embedded to have forward compatible implementations.
type UnimplementedCodeRepoBindingServer struct {
}

func (UnimplementedCodeRepoBindingServer) GetCodeRepoBinding(context.Context, *GetRequest) (*GetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCodeRepoBinding not implemented")
}
func (UnimplementedCodeRepoBindingServer) ListCodeRepoBindings(context.Context, *ListsRequest) (*ListsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCodeRepoBindings not implemented")
}
func (UnimplementedCodeRepoBindingServer) SaveCodeRepoBinding(context.Context, *SaveRequest) (*SaveReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveCodeRepoBinding not implemented")
}
func (UnimplementedCodeRepoBindingServer) DeleteCodeRepoBinding(context.Context, *DeleteRequest) (*DeleteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCodeRepoBinding not implemented")
}
func (UnimplementedCodeRepoBindingServer) mustEmbedUnimplementedCodeRepoBindingServer() {}

// UnsafeCodeRepoBindingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CodeRepoBindingServer will
// result in compilation errors.
type UnsafeCodeRepoBindingServer interface {
	mustEmbedUnimplementedCodeRepoBindingServer()
}

func RegisterCodeRepoBindingServer(s grpc.ServiceRegistrar, srv CodeRepoBindingServer) {
	s.RegisterService(&CodeRepoBinding_ServiceDesc, srv)
}

func _CodeRepoBinding_GetCodeRepoBinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CodeRepoBindingServer).GetCodeRepoBinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CodeRepoBinding_GetCodeRepoBinding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CodeRepoBindingServer).GetCodeRepoBinding(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CodeRepoBinding_ListCodeRepoBindings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CodeRepoBindingServer).ListCodeRepoBindings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CodeRepoBinding_ListCodeRepoBindings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CodeRepoBindingServer).ListCodeRepoBindings(ctx, req.(*ListsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CodeRepoBinding_SaveCodeRepoBinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CodeRepoBindingServer).SaveCodeRepoBinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CodeRepoBinding_SaveCodeRepoBinding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CodeRepoBindingServer).SaveCodeRepoBinding(ctx, req.(*SaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CodeRepoBinding_DeleteCodeRepoBinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CodeRepoBindingServer).DeleteCodeRepoBinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CodeRepoBinding_DeleteCodeRepoBinding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CodeRepoBindingServer).DeleteCodeRepoBinding(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CodeRepoBinding_ServiceDesc is the grpc.ServiceDesc for CodeRepoBinding service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CodeRepoBinding_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.coderepobinding.v1.CodeRepoBinding",
	HandlerType: (*CodeRepoBindingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCodeRepoBinding",
			Handler:    _CodeRepoBinding_GetCodeRepoBinding_Handler,
		},
		{
			MethodName: "ListCodeRepoBindings",
			Handler:    _CodeRepoBinding_ListCodeRepoBindings_Handler,
		},
		{
			MethodName: "SaveCodeRepoBinding",
			Handler:    _CodeRepoBinding_SaveCodeRepoBinding_Handler,
		},
		{
			MethodName: "DeleteCodeRepoBinding",
			Handler:    _CodeRepoBinding_DeleteCodeRepoBinding_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "coderepobinding/v1/coderepobinding.proto",
}
