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
// source: cluster/v1/cluster.proto

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
	Cluster_GetCluster_FullMethodName    = "/api.cluster.v1.Cluster/GetCluster"
	Cluster_ListClusters_FullMethodName  = "/api.cluster.v1.Cluster/ListClusters"
	Cluster_SaveCluster_FullMethodName   = "/api.cluster.v1.Cluster/SaveCluster"
	Cluster_DeleteCluster_FullMethodName = "/api.cluster.v1.Cluster/DeleteCluster"
)

// ClusterClient is the client API for Cluster service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClusterClient interface {
	GetCluster(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error)
	ListClusters(ctx context.Context, in *ListsRequest, opts ...grpc.CallOption) (*ListsReply, error)
	SaveCluster(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveReply, error)
	DeleteCluster(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error)
}

type clusterClient struct {
	cc grpc.ClientConnInterface
}

func NewClusterClient(cc grpc.ClientConnInterface) ClusterClient {
	return &clusterClient{cc}
}

func (c *clusterClient) GetCluster(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error) {
	out := new(GetReply)
	err := c.cc.Invoke(ctx, Cluster_GetCluster_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) ListClusters(ctx context.Context, in *ListsRequest, opts ...grpc.CallOption) (*ListsReply, error) {
	out := new(ListsReply)
	err := c.cc.Invoke(ctx, Cluster_ListClusters_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) SaveCluster(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveReply, error) {
	out := new(SaveReply)
	err := c.cc.Invoke(ctx, Cluster_SaveCluster_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) DeleteCluster(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error) {
	out := new(DeleteReply)
	err := c.cc.Invoke(ctx, Cluster_DeleteCluster_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterServer is the server API for Cluster service.
// All implementations must embed UnimplementedClusterServer
// for forward compatibility
type ClusterServer interface {
	GetCluster(context.Context, *GetRequest) (*GetReply, error)
	ListClusters(context.Context, *ListsRequest) (*ListsReply, error)
	SaveCluster(context.Context, *SaveRequest) (*SaveReply, error)
	DeleteCluster(context.Context, *DeleteRequest) (*DeleteReply, error)
	mustEmbedUnimplementedClusterServer()
}

// UnimplementedClusterServer must be embedded to have forward compatible implementations.
type UnimplementedClusterServer struct {
}

func (UnimplementedClusterServer) GetCluster(context.Context, *GetRequest) (*GetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCluster not implemented")
}
func (UnimplementedClusterServer) ListClusters(context.Context, *ListsRequest) (*ListsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListClusters not implemented")
}
func (UnimplementedClusterServer) SaveCluster(context.Context, *SaveRequest) (*SaveReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveCluster not implemented")
}
func (UnimplementedClusterServer) DeleteCluster(context.Context, *DeleteRequest) (*DeleteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCluster not implemented")
}
func (UnimplementedClusterServer) mustEmbedUnimplementedClusterServer() {}

// UnsafeClusterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClusterServer will
// result in compilation errors.
type UnsafeClusterServer interface {
	mustEmbedUnimplementedClusterServer()
}

func RegisterClusterServer(s grpc.ServiceRegistrar, srv ClusterServer) {
	s.RegisterService(&Cluster_ServiceDesc, srv)
}

func _Cluster_GetCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).GetCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cluster_GetCluster_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).GetCluster(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_ListClusters_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).ListClusters(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cluster_ListClusters_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).ListClusters(ctx, req.(*ListsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_SaveCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).SaveCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cluster_SaveCluster_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).SaveCluster(ctx, req.(*SaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_DeleteCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).DeleteCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cluster_DeleteCluster_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).DeleteCluster(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Cluster_ServiceDesc is the grpc.ServiceDesc for Cluster service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cluster_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.cluster.v1.Cluster",
	HandlerType: (*ClusterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCluster",
			Handler:    _Cluster_GetCluster_Handler,
		},
		{
			MethodName: "ListClusters",
			Handler:    _Cluster_ListClusters_Handler,
		},
		{
			MethodName: "SaveCluster",
			Handler:    _Cluster_SaveCluster_Handler,
		},
		{
			MethodName: "DeleteCluster",
			Handler:    _Cluster_DeleteCluster_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cluster/v1/cluster.proto",
}
