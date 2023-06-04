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

package service

import (
	"context"
	"fmt"

	coderepobindingv1 "github.com/nautes-labs/api-server/api/coderepobinding/v1"
	"github.com/nautes-labs/api-server/internal/biz"
	"github.com/nautes-labs/api-server/pkg/nodestree"
	resourcev1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
)

type CodeRepoBindingService struct {
	coderepobindingv1.UnimplementedCodeRepoBindingServer
	codeRepoBindingUsecase *biz.CodeRepoBindingUsecase
	resourcesUsecase       *biz.ResourcesUsecase
}

func NewCodeRepoBindingService(codeRepoBindingUsecase *biz.CodeRepoBindingUsecase, resourcesUsecase *biz.ResourcesUsecase) *CodeRepoBindingService {
	return &CodeRepoBindingService{codeRepoBindingUsecase: codeRepoBindingUsecase, resourcesUsecase: resourcesUsecase}
}

func (s *CodeRepoBindingService) CovertCodeRepoBindingValueToReply(codeRepoBinding *resourcev1alpha1.CodeRepoBinding) *coderepobindingv1.GetReply {
	return &coderepobindingv1.GetReply{
		Name:        codeRepoBinding.Name,
		Coderepo:    codeRepoBinding.Spec.CodeRepo,
		Product:     codeRepoBinding.Spec.Product,
		Projects:    codeRepoBinding.Spec.Projects,
		Permissions: codeRepoBinding.Spec.Permissions,
	}
}

func (s *CodeRepoBindingService) GetCodeRepoBinding(ctx context.Context, req *coderepobindingv1.GetRequest) (*coderepobindingv1.GetReply, error) {
	options := &biz.BizOptions{
		ProductName: req.ProductName,
		ResouceName: req.CoderepoBindingName,
	}
	codeRepoBinding, err := s.codeRepoBindingUsecase.GetCodeRepoBinding(ctx, options)
	if err != nil {
		return nil, err
	}

	return s.CovertCodeRepoBindingValueToReply(codeRepoBinding), nil
}

func (s *CodeRepoBindingService) ListCodeRepoBindings(ctx context.Context, req *coderepobindingv1.ListsRequest) (*coderepobindingv1.ListsReply, error) {
	options := &biz.BizOptions{
		ProductName: req.ProductName,
	}

	codeRepoBindings, err := s.codeRepoBindingUsecase.ListCodeRepoBindings(ctx, options)
	if err != nil {
		return nil, err
	}

	var items []*coderepobindingv1.GetReply
	for _, codeRepoBinding := range codeRepoBindings {
		items = append(items, s.CovertCodeRepoBindingValueToReply(codeRepoBinding))
	}

	return &coderepobindingv1.ListsReply{Items: items}, nil
}

func (s *CodeRepoBindingService) SaveCodeRepoBinding(ctx context.Context, req *coderepobindingv1.SaveRequest) (*coderepobindingv1.SaveReply, error) {
	productResourceName, err := s.resourcesUsecase.ConvertGroupToProduct(ctx, req.ProductName)
	if err != nil {
		return nil, err
	}

	codeRepoResourceName, err := s.resourcesUsecase.ConvertRepoNameToCodeRepo(ctx, req.ProductName, req.Body.Coderepo)
	if err != nil {
		return nil, err
	}

	ctx = biz.SetResourceContext(ctx, req.ProductName, biz.SaveMethod, nodestree.CodeRepo, codeRepoResourceName, nodestree.CodeRepoBinding, req.CoderepoBindingName)

	options := &biz.BizOptions{
		ProductName:       req.ProductName,
		ResouceName:       req.CoderepoBindingName,
		InsecureSkipCheck: req.InsecureSkipCheck,
	}

	data := &biz.CodeRepoBindingData{
		Name: req.CoderepoBindingName,
		Spec: resourcev1alpha1.CodeRepoBindingSpec{
			CodeRepo:    codeRepoResourceName,
			Product:     productResourceName,
			Projects:    req.Body.Projects,
			Permissions: req.Body.Permissions,
		},
	}

	if err := s.codeRepoBindingUsecase.SaveCodeRepoBinding(ctx, options, data); err != nil {
		return nil, err
	}

	return &coderepobindingv1.SaveReply{
		Msg: fmt.Sprintf("Successfully saved %v configuration", req.CoderepoBindingName),
	}, nil
}

func (s *CodeRepoBindingService) DeleteCodeRepoBinding(ctx context.Context, req *coderepobindingv1.DeleteRequest) (*coderepobindingv1.DeleteReply, error) {
	ctx = biz.SetResourceContext(ctx, req.ProductName, biz.DeleteMethod, "", "", nodestree.CodeRepoBinding, req.CoderepoBindingName)

	options := &biz.BizOptions{
		ProductName:       req.ProductName,
		ResouceName:       req.CoderepoBindingName,
		InsecureSkipCheck: req.InsecureSkipCheck,
	}

	if err := s.codeRepoBindingUsecase.DeleteCodeRepoBinding(ctx, options); err != nil {
		return nil, err
	}

	return &coderepobindingv1.DeleteReply{
		Msg: fmt.Sprintf("Successfully deleted %v configuration", req.CoderepoBindingName),
	}, nil
}
