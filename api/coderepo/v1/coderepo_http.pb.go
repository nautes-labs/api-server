// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.6.3
// - protoc             v3.12.4
// source: coderepo/v1/coderepo.proto

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

const OperationCodeRepoDeleteCodeRepo = "/api.coderepo.v1.CodeRepo/DeleteCodeRepo"
const OperationCodeRepoGetCodeRepo = "/api.coderepo.v1.CodeRepo/GetCodeRepo"
const OperationCodeRepoListCodeRepos = "/api.coderepo.v1.CodeRepo/ListCodeRepos"
const OperationCodeRepoSaveCodeRepo = "/api.coderepo.v1.CodeRepo/SaveCodeRepo"

type CodeRepoHTTPServer interface {
	DeleteCodeRepo(context.Context, *DeleteRequest) (*DeleteReply, error)
	GetCodeRepo(context.Context, *GetRequest) (*GetReply, error)
	ListCodeRepos(context.Context, *ListsRequest) (*ListsReply, error)
	SaveCodeRepo(context.Context, *SaveRequest) (*SaveReply, error)
}

func RegisterCodeRepoHTTPServer(s *http.Server, srv CodeRepoHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/products/{product_name}/coderepos/{coderepoName}", _CodeRepo_GetCodeRepo0_HTTP_Handler(srv))
	r.GET("/api/v1/products/{product_name}/coderepos", _CodeRepo_ListCodeRepos0_HTTP_Handler(srv))
	r.POST("/api/v1/products/{product_name}/coderepos/{coderepoName}", _CodeRepo_SaveCodeRepo0_HTTP_Handler(srv))
	r.DELETE("/api/v1/products/{product_name}/coderepos/{coderepoName}", _CodeRepo_DeleteCodeRepo0_HTTP_Handler(srv))
}

func _CodeRepo_GetCodeRepo0_HTTP_Handler(srv CodeRepoHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCodeRepoGetCodeRepo)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetCodeRepo(ctx, req.(*GetRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetReply)
		return ctx.Result(200, reply)
	}
}

func _CodeRepo_ListCodeRepos0_HTTP_Handler(srv CodeRepoHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListsRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCodeRepoListCodeRepos)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListCodeRepos(ctx, req.(*ListsRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListsReply)
		return ctx.Result(200, reply)
	}
}

func _CodeRepo_SaveCodeRepo0_HTTP_Handler(srv CodeRepoHTTPServer) func(ctx http.Context) error {
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
		http.SetOperation(ctx, OperationCodeRepoSaveCodeRepo)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SaveCodeRepo(ctx, req.(*SaveRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SaveReply)
		return ctx.Result(200, reply)
	}
}

func _CodeRepo_DeleteCodeRepo0_HTTP_Handler(srv CodeRepoHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCodeRepoDeleteCodeRepo)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteCodeRepo(ctx, req.(*DeleteRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteReply)
		return ctx.Result(200, reply)
	}
}

type CodeRepoHTTPClient interface {
	DeleteCodeRepo(ctx context.Context, req *DeleteRequest, opts ...http.CallOption) (rsp *DeleteReply, err error)
	GetCodeRepo(ctx context.Context, req *GetRequest, opts ...http.CallOption) (rsp *GetReply, err error)
	ListCodeRepos(ctx context.Context, req *ListsRequest, opts ...http.CallOption) (rsp *ListsReply, err error)
	SaveCodeRepo(ctx context.Context, req *SaveRequest, opts ...http.CallOption) (rsp *SaveReply, err error)
}

type CodeRepoHTTPClientImpl struct {
	cc *http.Client
}

func NewCodeRepoHTTPClient(client *http.Client) CodeRepoHTTPClient {
	return &CodeRepoHTTPClientImpl{client}
}

func (c *CodeRepoHTTPClientImpl) DeleteCodeRepo(ctx context.Context, in *DeleteRequest, opts ...http.CallOption) (*DeleteReply, error) {
	var out DeleteReply
	pattern := "/api/v1/products/{product_name}/coderepos/{coderepoName}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationCodeRepoDeleteCodeRepo))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *CodeRepoHTTPClientImpl) GetCodeRepo(ctx context.Context, in *GetRequest, opts ...http.CallOption) (*GetReply, error) {
	var out GetReply
	pattern := "/api/v1/products/{product_name}/coderepos/{coderepoName}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationCodeRepoGetCodeRepo))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *CodeRepoHTTPClientImpl) ListCodeRepos(ctx context.Context, in *ListsRequest, opts ...http.CallOption) (*ListsReply, error) {
	var out ListsReply
	pattern := "/api/v1/products/{product_name}/coderepos"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationCodeRepoListCodeRepos))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *CodeRepoHTTPClientImpl) SaveCodeRepo(ctx context.Context, in *SaveRequest, opts ...http.CallOption) (*SaveReply, error) {
	var out SaveReply
	pattern := "/api/v1/products/{product_name}/coderepos/{coderepoName}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationCodeRepoSaveCodeRepo))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in.Body, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
