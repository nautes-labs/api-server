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

package validate

import (
	"context"
	"fmt"

	"github.com/nautes-labs/api-server/pkg/nodestree"
	resourcev1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type validateClient struct {
	client               client.Client
	nodes                *nodestree.Node
	nodestree            nodestree.NodesTree
	tenantAdminNamespace string
}

func NewValidateClient(client client.Client, nodestree nodestree.NodesTree, nodes *nodestree.Node, tenantAdminNamespace string) resourcev1alpha1.ValidateClient {
	return &validateClient{client: client, nodestree: nodestree, nodes: nodes, tenantAdminNamespace: tenantAdminNamespace}
}

func (v *validateClient) GetCodeRepo(ctx context.Context, repoName string) (*resourcev1alpha1.CodeRepo, error) {
	node := v.nodestree.GetNode(v.nodes, nodestree.CodeRepo, repoName)
	if node == nil {
		return nil, fmt.Errorf("node %s not found", repoName)
	}
	codeRepo, ok := node.Content.(*resourcev1alpha1.CodeRepo)
	if !ok {
		return nil, fmt.Errorf("the node %s content type is wrong", node.Name)
	}
	return codeRepo, nil
}

func (v *validateClient) GetEnvironment(ctx context.Context, productName, name string) (*resourcev1alpha1.Environment, error) {
	node := v.nodestree.GetNode(v.nodes, nodestree.Environment, name)
	if node == nil {
		return nil, fmt.Errorf("the environment %s is not found", name)
	}
	env, ok := node.Content.(*resourcev1alpha1.Environment)
	if !ok {
		return nil, fmt.Errorf("the node %s content type is wrong", node.Name)
	}

	return env, nil
}

func (v *validateClient) GetCluster(ctx context.Context, name string) (*resourcev1alpha1.Cluster, error) {
	objKey := client.ObjectKey{
		Namespace: v.tenantAdminNamespace,
		Name:      name,
	}

	cluster := &resourcev1alpha1.Cluster{}
	err := v.client.Get(ctx, objKey, cluster)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (v *validateClient) ListCodeRepoBinding(ctx context.Context, productName, repoName string) (*resourcev1alpha1.CodeRepoBindingList, error) {
	nodes := nodestree.ListsResourceNodes(*v.nodes, nodestree.CodeRepoBinding)

	list := &resourcev1alpha1.CodeRepoBindingList{}
	for _, node := range nodes {
		codeRepoBinding, ok := node.Content.(*resourcev1alpha1.CodeRepoBinding)
		if !ok {
			continue
		}
		if codeRepoBinding.Spec.CodeRepo == repoName {
			list.Items = append(list.Items, *codeRepoBinding)
		}
	}

	return list, nil
}

func (v *validateClient) ListDeploymentRuntime(ctx context.Context, productName string) (*resourcev1alpha1.DeploymentRuntimeList, error) {
	nodes := nodestree.ListsResourceNodes(*v.nodes, nodestree.DeploymentRuntime)

	list := &resourcev1alpha1.DeploymentRuntimeList{}
	for _, node := range nodes {
		runtime, ok := node.Content.(*resourcev1alpha1.DeploymentRuntime)
		if !ok {
			continue
		}
		runtime.Namespace = productName
		list.Items = append(list.Items, *runtime)
	}

	return list, nil
}

func (v *validateClient) ListProjectPipelineRuntime(ctx context.Context, productName string) (*resourcev1alpha1.ProjectPipelineRuntimeList, error) {
	nodes := nodestree.ListsResourceNodes(*v.nodes, nodestree.ProjectPipelineRuntime)

	list := &resourcev1alpha1.ProjectPipelineRuntimeList{}
	for _, node := range nodes {
		runtime, ok := node.Content.(*resourcev1alpha1.ProjectPipelineRuntime)
		if !ok {
			continue
		}
		runtime.Namespace = productName
		list.Items = append(list.Items, *runtime)
	}

	return list, nil
}
