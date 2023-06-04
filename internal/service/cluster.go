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

	clusterv1 "github.com/nautes-labs/api-server/api/cluster/v1"
	"github.com/nautes-labs/api-server/internal/biz"
	ClusterRegistration "github.com/nautes-labs/api-server/pkg/cluster"
	registercluster "github.com/nautes-labs/api-server/pkg/cluster"
	"github.com/nautes-labs/api-server/pkg/nodestree"
	resourcev1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
	nautesconfigs "github.com/nautes-labs/pkg/pkg/nautesconfigs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterService struct {
	clusterv1.UnimplementedClusterServer
	cluster *biz.ClusterUsecase
	configs *nautesconfigs.Config
}

func NewClusterService(cluster *biz.ClusterUsecase, configs *nautesconfigs.Config) *ClusterService {
	return &ClusterService{cluster: cluster, configs: configs}
}

func (s *ClusterService) SaveCluster(ctx context.Context, req *clusterv1.SaveRequest) (*clusterv1.SaveReply, error) {
	err := s.Validate(req)
	if err != nil {
		return nil, err
	}

	ctx = biz.SetResourceContext(ctx, "", biz.SaveMethod, "", "", nodestree.Cluster, req.ClusterName)

	cluster := &resourcev1alpha1.Cluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: resourcev1alpha1.GroupVersion.String(),
			Kind:       nodestree.Cluster,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.ClusterName,
			Namespace: s.configs.Nautes.Namespace,
		},
		Spec: resourcev1alpha1.ClusterSpec{
			ApiServer:     req.Body.ApiServer,
			ClusterType:   resourcev1alpha1.ClusterType(req.Body.ClusterType),
			ClusterKind:   resourcev1alpha1.ClusterKind(req.Body.ClusterKind),
			Usage:         resourcev1alpha1.ClusterUsage(req.Body.Usage),
			HostCluster:   req.Body.HostCluster,
			PrimaryDomain: req.Body.PrimaryDomain,
			WorkerType:    resourcev1alpha1.ClusterWorkType(req.Body.WorkerType),
		},
	}

	var vcluster *registercluster.Vcluster
	if ok := registercluster.IsVirtualDeploymentRuntime(cluster); ok {
		if req.Body.Vcluster != nil {
			vcluster = &registercluster.Vcluster{
				HttpsNodePort: req.Body.Vcluster.HttpsNodePort,
			}
		}
	}

	traefik := &ClusterRegistration.Traefik{}
	if req.Body.ClusterType == string(resourcev1alpha1.CLUSTER_TYPE_PHYSICAL) {
		traefik.HttpNodePort = req.Body.Traefik.HttpNodePort
		traefik.HttpsNodePort = req.Body.Traefik.HttpsNodePort
	}

	param := &ClusterRegistration.ClusterRegistrationParam{
		Cluster:    cluster,
		ArgocdHost: req.Body.ArgocdHost,
		TektonHost: req.Body.TektonHost,
		Traefik:    traefik,
		Vcluster:   vcluster,
	}

	if err := s.cluster.SaveCluster(ctx, param, req.Body.Kubeconfig); err != nil {
		return nil, err
	}

	return &clusterv1.SaveReply{
		Msg: fmt.Sprintf("Successfully saved %s cluster", req.ClusterName),
	}, nil
}

func (s *ClusterService) DeleteCluster(ctx context.Context, req *clusterv1.DeleteRequest) (*clusterv1.DeleteReply, error) {
	ctx = biz.SetResourceContext(ctx, "", biz.DeleteMethod, "", "", nodestree.Cluster, req.ClusterName)

	err := s.cluster.DeleteCluster(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	return &clusterv1.DeleteReply{
		Msg: fmt.Sprintf("Successfully deleted %s cluster", req.ClusterName),
	}, nil
}

func (s *ClusterService) Validate(req *clusterv1.SaveRequest) error {
	if req.Body.Usage == string(resourcev1alpha1.CLUSTER_USAGE_WORKER) &&
		req.Body.ClusterType == string(resourcev1alpha1.CLUSTER_TYPE_VIRTUAL) &&
		req.Body.HostCluster == "" {
		return fmt.Errorf("the 'Host Cluster' for virtual cluster is required")
	}

	if req.Body.Usage == string(resourcev1alpha1.CLUSTER_USAGE_WORKER) &&
		req.Body.WorkerType == "" {
		return fmt.Errorf("when the cluster usage is 'worker', the 'WorkerType' is required")
	}

	if req.Body.ClusterType == string(resourcev1alpha1.CLUSTER_TYPE_PHYSICAL) &&
		req.Body.Traefik == nil {
		return fmt.Errorf("traefik parameter is required when cluster type is 'Host Cluster' or 'Physical Runtime'")
	}

	return nil
}
