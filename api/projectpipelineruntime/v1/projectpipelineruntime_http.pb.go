// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.6.2
// - protoc             v3.6.1
// source: projectpipelineruntime/v1/projectpipelineruntime.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationProjectPipelineRuntimeDeleteProjectPipelineRuntime = "/api.projectpipelineruntime.v1.ProjectPipelineRuntime/DeleteProjectPipelineRuntime"
const OperationProjectPipelineRuntimeGetProjectPipelineRuntime = "/api.projectpipelineruntime.v1.ProjectPipelineRuntime/GetProjectPipelineRuntime"
const OperationProjectPipelineRuntimeListProjectPipelineRuntimes = "/api.projectpipelineruntime.v1.ProjectPipelineRuntime/ListProjectPipelineRuntimes"
const OperationProjectPipelineRuntimeSaveProjectPipelineRuntime = "/api.projectpipelineruntime.v1.ProjectPipelineRuntime/SaveProjectPipelineRuntime"

type ProjectPipelineRuntimeHTTPServer interface {
	DeleteProjectPipelineRuntime(context.Context, *DeleteRequest) (*DeleteReply, error)
	GetProjectPipelineRuntime(context.Context, *GetRequest) (*GetReply, error)
	ListProjectPipelineRuntimes(context.Context, *ListsRequest) (*ListsReply, error)
	SaveProjectPipelineRuntime(context.Context, *SaveRequest) (*SaveReply, error)
}

func RegisterProjectPipelineRuntimeHTTPServer(s *http.Server, srv ProjectPipelineRuntimeHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/products/{product_name}/projectpipelineruntimes/{project_pipeline_runtime_name}", _ProjectPipelineRuntime_GetProjectPipelineRuntime0_HTTP_Handler(srv))
	r.GET("/api/v1/products/{product_name}/projectpipelineruntimes", _ProjectPipelineRuntime_ListProjectPipelineRuntimes0_HTTP_Handler(srv))
	r.POST("/api/v1/products/{product_name}/projectpipelineruntimes/{project_pipeline_runtime_name}", _ProjectPipelineRuntime_SaveProjectPipelineRuntime0_HTTP_Handler(srv))
	r.DELETE("/api/v1/products/{product_name}/projectpipelineruntimes/{project_pipeline_runtime_name}", _ProjectPipelineRuntime_DeleteProjectPipelineRuntime0_HTTP_Handler(srv))
}

func _ProjectPipelineRuntime_GetProjectPipelineRuntime0_HTTP_Handler(srv ProjectPipelineRuntimeHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProjectPipelineRuntimeGetProjectPipelineRuntime)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetProjectPipelineRuntime(ctx, req.(*GetRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetReply)
		return ctx.Result(200, reply)
	}
}

func _ProjectPipelineRuntime_ListProjectPipelineRuntimes0_HTTP_Handler(srv ProjectPipelineRuntimeHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListsRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProjectPipelineRuntimeListProjectPipelineRuntimes)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListProjectPipelineRuntimes(ctx, req.(*ListsRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListsReply)
		return ctx.Result(200, reply)
	}
}

func _ProjectPipelineRuntime_SaveProjectPipelineRuntime0_HTTP_Handler(srv ProjectPipelineRuntimeHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in SaveRequest
		if err := ctx.Bind(&in.Body); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProjectPipelineRuntimeSaveProjectPipelineRuntime)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SaveProjectPipelineRuntime(ctx, req.(*SaveRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SaveReply)
		return ctx.Result(200, reply)
	}
}

func _ProjectPipelineRuntime_DeleteProjectPipelineRuntime0_HTTP_Handler(srv ProjectPipelineRuntimeHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProjectPipelineRuntimeDeleteProjectPipelineRuntime)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteProjectPipelineRuntime(ctx, req.(*DeleteRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteReply)
		return ctx.Result(200, reply)
	}
}

type ProjectPipelineRuntimeHTTPClient interface {
	DeleteProjectPipelineRuntime(ctx context.Context, req *DeleteRequest, opts ...http.CallOption) (rsp *DeleteReply, err error)
	GetProjectPipelineRuntime(ctx context.Context, req *GetRequest, opts ...http.CallOption) (rsp *GetReply, err error)
	ListProjectPipelineRuntimes(ctx context.Context, req *ListsRequest, opts ...http.CallOption) (rsp *ListsReply, err error)
	SaveProjectPipelineRuntime(ctx context.Context, req *SaveRequest, opts ...http.CallOption) (rsp *SaveReply, err error)
}

type ProjectPipelineRuntimeHTTPClientImpl struct {
	cc *http.Client
}

func NewProjectPipelineRuntimeHTTPClient(client *http.Client) ProjectPipelineRuntimeHTTPClient {
	return &ProjectPipelineRuntimeHTTPClientImpl{client}
}

func (c *ProjectPipelineRuntimeHTTPClientImpl) DeleteProjectPipelineRuntime(ctx context.Context, in *DeleteRequest, opts ...http.CallOption) (*DeleteReply, error) {
	var out DeleteReply
	pattern := "/api/v1/products/{product_name}/projectpipelineruntimes/{project_pipeline_runtime_name}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProjectPipelineRuntimeDeleteProjectPipelineRuntime))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *ProjectPipelineRuntimeHTTPClientImpl) GetProjectPipelineRuntime(ctx context.Context, in *GetRequest, opts ...http.CallOption) (*GetReply, error) {
	var out GetReply
	pattern := "/api/v1/products/{product_name}/projectpipelineruntimes/{project_pipeline_runtime_name}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProjectPipelineRuntimeGetProjectPipelineRuntime))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *ProjectPipelineRuntimeHTTPClientImpl) ListProjectPipelineRuntimes(ctx context.Context, in *ListsRequest, opts ...http.CallOption) (*ListsReply, error) {
	var out ListsReply
	pattern := "/api/v1/products/{product_name}/projectpipelineruntimes"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProjectPipelineRuntimeListProjectPipelineRuntimes))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *ProjectPipelineRuntimeHTTPClientImpl) SaveProjectPipelineRuntime(ctx context.Context, in *SaveRequest, opts ...http.CallOption) (*SaveReply, error) {
	var out SaveReply
	pattern := "/api/v1/products/{product_name}/projectpipelineruntimes/{project_pipeline_runtime_name}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationProjectPipelineRuntimeSaveProjectPipelineRuntime))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in.Body, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
