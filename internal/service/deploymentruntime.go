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

	deploymentruntimev1 "github.com/nautes-labs/api-server/api/deploymentruntime/v1"
	"github.com/nautes-labs/api-server/internal/biz"
	"github.com/nautes-labs/api-server/pkg/nodestree"
	"github.com/nautes-labs/api-server/pkg/selector"
	"github.com/nautes-labs/pkg/api/v1alpha1"
	resourcev1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
)

var (
	deploymentRuntimeFilterFieldRules = map[string]map[string]selector.FieldSelector{
		FieldProjectsRef: {
			selector.EqualOperator: selector.NewStringSelector(_ProjectsRef, selector.In),
		},
		FeldManifestSourceCodeRepo: {
			selector.EqualOperator: selector.NewStringSelector(_ManifestSource, selector.In),
		},
		FieldDestination: {
			selector.EqualOperator: selector.NewStringSelector(_Destination, selector.In),
		},
	}
)

type DeploymentruntimeService struct {
	deploymentruntimev1.UnimplementedDeploymentruntimeServer
	deploymentRuntime *biz.DeploymentRuntimeUsecase
}

func NewDeploymentruntimeService(deploymentRuntime *biz.DeploymentRuntimeUsecase) *DeploymentruntimeService {
	return &DeploymentruntimeService{deploymentRuntime: deploymentRuntime}
}

func (s *DeploymentruntimeService) CovertDeploymentRuntimeValueToReply(runtime *resourcev1alpha1.DeploymentRuntime) *deploymentruntimev1.GetReply {
	return &deploymentruntimev1.GetReply{
		Product:     runtime.Spec.Product,
		Name:        runtime.Name,
		Destination: runtime.Spec.Destination,
		ProjectsRef: runtime.Spec.ProjectsRef,
		ManifestSource: &deploymentruntimev1.ManifestSource{
			CodeRepo:       runtime.Spec.ManifestSource.CodeRepo,
			TargetRevision: runtime.Spec.ManifestSource.TargetRevision,
			Path:           runtime.Spec.ManifestSource.Path,
		},
	}
}

func (s *DeploymentruntimeService) GetDeploymentRuntime(ctx context.Context, req *deploymentruntimev1.GetRequest) (*deploymentruntimev1.GetReply, error) {
	runtime, err := s.deploymentRuntime.GetDeploymentRuntime(ctx, req.DeploymentruntimeName, req.ProductName)
	if err != nil {
		return nil, err
	}

	return s.CovertDeploymentRuntimeValueToReply(runtime), nil
}

func (s *DeploymentruntimeService) ListDeploymentRuntimes(ctx context.Context, req *deploymentruntimev1.ListsRequest) (*deploymentruntimev1.ListsReply, error) {
	nodes, err := s.deploymentRuntime.ListDeploymentRuntimes(ctx, req.ProductName)
	if err != nil {
		return nil, err
	}

	var items []*deploymentruntimev1.GetReply
	for _, node := range nodes {
		runtime, ok := node.Content.(*resourcev1alpha1.DeploymentRuntime)
		if !ok {
			continue
		}

		err := s.deploymentRuntime.ConvertRuntime(ctx, runtime)
		if err != nil {
			return nil, err
		}
		node.Content = runtime

		passed, err := selector.Match(req.FieldSelector, node.Content, deploymentRuntimeFilterFieldRules)
		if err != nil {
			return nil, err
		}
		if !passed {
			continue
		}

		items = append(items, s.CovertDeploymentRuntimeValueToReply(runtime))
	}

	return &deploymentruntimev1.ListsReply{
		Items: items,
	}, nil
}

func (s *DeploymentruntimeService) SaveDeploymentRuntime(ctx context.Context, req *deploymentruntimev1.SaveRequest) (*deploymentruntimev1.SaveReply, error) {
	ctx = biz.SetResourceContext(ctx, req.ProductName, biz.SaveMethod, "", "", nodestree.DeploymentRuntime, req.DeploymentruntimeName)

	data := &biz.DeploymentRuntimeData{
		Name: req.DeploymentruntimeName,
		Spec: v1alpha1.DeploymentRuntimeSpec{
			Product:     req.ProductName,
			ProjectsRef: req.Body.ProjectsRef,
			Destination: req.Body.Destination,
			ManifestSource: resourcev1alpha1.ManifestSource{
				CodeRepo:       req.Body.ManifestSource.CodeRepo,
				TargetRevision: req.Body.ManifestSource.TargetRevision,
				Path:           req.Body.ManifestSource.Path,
			},
		},
	}
	options := &biz.BizOptions{
		ResouceName:       req.DeploymentruntimeName,
		ProductName:       req.ProductName,
		InsecureSkipCheck: req.InsecureSkipCheck,
	}
	err := s.deploymentRuntime.SaveDeploymentRuntime(ctx, options, data)
	if err != nil {
		return nil, err
	}

	return &deploymentruntimev1.SaveReply{
		Msg: fmt.Sprintf("Successfully saved %s configuration", req.DeploymentruntimeName),
	}, nil
}

func (s *DeploymentruntimeService) DeleteDeploymentRuntime(ctx context.Context, req *deploymentruntimev1.DeleteRequest) (*deploymentruntimev1.DeleteReply, error) {
	ctx = biz.SetResourceContext(ctx, req.ProductName, biz.DeleteMethod, "", "", nodestree.DeploymentRuntime, req.DeploymentruntimeName)

	options := &biz.BizOptions{
		ResouceName:       req.DeploymentruntimeName,
		ProductName:       req.ProductName,
		InsecureSkipCheck: req.InsecureSkipCheck,
	}
	err := s.deploymentRuntime.DeleteDeploymentRuntime(ctx, options)
	if err != nil {
		return nil, err
	}

	return &deploymentruntimev1.DeleteReply{
		Msg: fmt.Sprintf("Successfully deleted %s configuration", req.DeploymentruntimeName),
	}, nil
}
