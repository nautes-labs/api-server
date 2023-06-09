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

package data

import (
	"context"

	"github.com/nautes-labs/api-server/internal/biz"
	"github.com/nautes-labs/api-server/pkg/kubernetes"
	resourcev1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Kubernetes struct {
	Client client.Client
}

func NewKubernetes() (biz.Kubernetes, error) {
	client, err := kubernetes.NewKubernetes()
	if err != nil {
		return nil, err
	}
	return &Kubernetes{Client: client}, nil
}

func (k *Kubernetes) ListCodeRepoBindings(ctx context.Context) (*resourcev1alpha1.CodeRepoBindingList, error) {
	if k.Client == nil {
		_, err := NewKubernetes()
		if err != nil {
			return nil, err
		}
	}

	lists := &resourcev1alpha1.CodeRepoBindingList{}
	if err := k.Client.List(ctx, lists); err != nil {
		return nil, err
	}

	return lists, nil
}

func (k *Kubernetes) ListCodeRepo(ctx context.Context) (*resourcev1alpha1.CodeRepoList, error) {
	if k.Client == nil {
		_, err := NewKubernetes()
		if err != nil {
			return nil, err
		}
	}

	lists := &resourcev1alpha1.CodeRepoList{}
	if err := k.Client.List(ctx, lists); err != nil {
		return nil, err
	}

	return lists, nil
}
