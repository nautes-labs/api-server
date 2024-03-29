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
// source: coderepo/v1/coderepo.proto

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
	CodeRepo_GetCodeRepo_FullMethodName    = "/api.coderepo.v1.CodeRepo/GetCodeRepo"
	CodeRepo_ListCodeRepos_FullMethodName  = "/api.coderepo.v1.CodeRepo/ListCodeRepos"
	CodeRepo_SaveCodeRepo_FullMethodName   = "/api.coderepo.v1.CodeRepo/SaveCodeRepo"
	CodeRepo_DeleteCodeRepo_FullMethodName = "/api.coderepo.v1.CodeRepo/DeleteCodeRepo"
)

// CodeRepoClient is the client API for CodeRepo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CodeRepoClient interface {
	GetCodeRepo(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error)
	ListCodeRepos(ctx context.Context, in *ListsRequest, opts ...grpc.CallOption) (*ListsReply, error)
	SaveCodeRepo(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveReply, error)
	DeleteCodeRepo(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error)
}

type codeRepoClient struct {
	cc grpc.ClientConnInterface
}

func NewCodeRepoClient(cc grpc.ClientConnInterface) CodeRepoClient {
	return &codeRepoClient{cc}
}

func (c *codeRepoClient) GetCodeRepo(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error) {
	out := new(GetReply)
	err := c.cc.Invoke(ctx, CodeRepo_GetCodeRepo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *codeRepoClient) ListCodeRepos(ctx context.Context, in *ListsRequest, opts ...grpc.CallOption) (*ListsReply, error) {
	out := new(ListsReply)
	err := c.cc.Invoke(ctx, CodeRepo_ListCodeRepos_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *codeRepoClient) SaveCodeRepo(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveReply, error) {
	out := new(SaveReply)
	err := c.cc.Invoke(ctx, CodeRepo_SaveCodeRepo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *codeRepoClient) DeleteCodeRepo(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error) {
	out := new(DeleteReply)
	err := c.cc.Invoke(ctx, CodeRepo_DeleteCodeRepo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CodeRepoServer is the server API for CodeRepo service.
// All implementations must embed UnimplementedCodeRepoServer
// for forward compatibility
type CodeRepoServer interface {
	GetCodeRepo(context.Context, *GetRequest) (*GetReply, error)
	ListCodeRepos(context.Context, *ListsRequest) (*ListsReply, error)
	SaveCodeRepo(context.Context, *SaveRequest) (*SaveReply, error)
	DeleteCodeRepo(context.Context, *DeleteRequest) (*DeleteReply, error)
	mustEmbedUnimplementedCodeRepoServer()
}

// UnimplementedCodeRepoServer must be embedded to have forward compatible implementations.
type UnimplementedCodeRepoServer struct {
}

func (UnimplementedCodeRepoServer) GetCodeRepo(context.Context, *GetRequest) (*GetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCodeRepo not implemented")
}
func (UnimplementedCodeRepoServer) ListCodeRepos(context.Context, *ListsRequest) (*ListsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCodeRepos not implemented")
}
func (UnimplementedCodeRepoServer) SaveCodeRepo(context.Context, *SaveRequest) (*SaveReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveCodeRepo not implemented")
}
func (UnimplementedCodeRepoServer) DeleteCodeRepo(context.Context, *DeleteRequest) (*DeleteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCodeRepo not implemented")
}
func (UnimplementedCodeRepoServer) mustEmbedUnimplementedCodeRepoServer() {}

// UnsafeCodeRepoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CodeRepoServer will
// result in compilation errors.
type UnsafeCodeRepoServer interface {
	mustEmbedUnimplementedCodeRepoServer()
}

func RegisterCodeRepoServer(s grpc.ServiceRegistrar, srv CodeRepoServer) {
	s.RegisterService(&CodeRepo_ServiceDesc, srv)
}

func _CodeRepo_GetCodeRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CodeRepoServer).GetCodeRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CodeRepo_GetCodeRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CodeRepoServer).GetCodeRepo(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CodeRepo_ListCodeRepos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CodeRepoServer).ListCodeRepos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CodeRepo_ListCodeRepos_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CodeRepoServer).ListCodeRepos(ctx, req.(*ListsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CodeRepo_SaveCodeRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CodeRepoServer).SaveCodeRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CodeRepo_SaveCodeRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CodeRepoServer).SaveCodeRepo(ctx, req.(*SaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CodeRepo_DeleteCodeRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CodeRepoServer).DeleteCodeRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CodeRepo_DeleteCodeRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CodeRepoServer).DeleteCodeRepo(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CodeRepo_ServiceDesc is the grpc.ServiceDesc for CodeRepo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CodeRepo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.coderepo.v1.CodeRepo",
	HandlerType: (*CodeRepoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCodeRepo",
			Handler:    _CodeRepo_GetCodeRepo_Handler,
		},
		{
			MethodName: "ListCodeRepos",
			Handler:    _CodeRepo_ListCodeRepos_Handler,
		},
		{
			MethodName: "SaveCodeRepo",
			Handler:    _CodeRepo_SaveCodeRepo_Handler,
		},
		{
			MethodName: "DeleteCodeRepo",
			Handler:    _CodeRepo_DeleteCodeRepo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "coderepo/v1/coderepo.proto",
}
