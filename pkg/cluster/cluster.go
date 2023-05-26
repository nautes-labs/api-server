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

package cluster

import (
	"fmt"
	"os"
	"regexp"

	utilstrings "github.com/nautes-labs/api-server/util/string"
	resourcev1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
)

const (
	_HostClusterDirectoreyPlaceholder                    = "_HOST_CLUSTER_"
	_RuntimeDirectoryDirectoreyPlaceholder               = "_RUNTIME_"
	_VclusterDirectoryDirectoreyPlaceholder              = "_VCLUSTER_"
	_HostCluster                            ClusterUsage = "HostCluster"
	_PhysicalDeploymentRuntime              ClusterUsage = "PhysicalDeploymentRuntime"
	_VirtualDeploymentRuntime               ClusterUsage = "VirtualDeploymentRuntime"
	_PhysicalProjectPipelineRuntime         ClusterUsage = "PhysicalProjectPipelineRuntime"
	_VirtualProjectPipelineRuntime          ClusterUsage = "VirtualProjectPipelineRuntime"
	_TektonPrefix                                        = "tekton"
	_ArgocdPrefix                                        = "argocd"
	_NIPDomainSuffix                                     = "nip.io"
	_RuntimeSuffix                                       = "runtime"
	_RuntimeProjectSuffix                                = "runtime-project"
	_ArgocdOAuthSuffix                                   = "api/dex/callback"
	_TektonOAuthSuffix                                   = "oauth2/callback"
)

func NewClusterRegistration() ClusterRegistrationOperator {
	return &ClusterRegistration{}
}

// InitializeClusterConfig Initialize the configuration information of the cluster,
// mainly including the host cluster, virtual cluster and runtime configuration information.
func (cr *ClusterRegistration) InitializeClusterConfig(param *ClusterRegistrationParam) error {
	var hostCluster *HostCluster
	var vcluster *Vcluster
	var runtime *Runtime
	var err error
	var resourceHostCluster, resourceVirtualCluster, resourceRuntimeCluster *resourcev1alpha1.Cluster

	if IsHostCluser(param.Cluster) {
		resourceHostCluster = param.Cluster.DeepCopy()
	}

	if IsVirtual(param.Cluster) {
		resourceVirtualCluster = param.Cluster.DeepCopy()
		resourceHostCluster, err = GetHostClusterFromTenantConfigFile(param.TenantConfigRepoLocalPath, resourceVirtualCluster.Spec.HostCluster, param.Configs.Nautes.TenantName)
		if err != nil {
			return err
		}
	}

	if !IsHostCluser(param.Cluster) {
		resourceRuntimeCluster = param.Cluster.DeepCopy()
	}

	hostCluster, err = cr.getHostCluster(param, resourceHostCluster, param.Configs.Nautes.TenantName)
	if err != nil {
		return err
	}

	vcluster, err = cr.getVcluster(param, resourceVirtualCluster, resourceHostCluster)
	if err != nil {
		return err
	}

	runtime, err = cr.getRuntime(param, resourceRuntimeCluster, hostCluster)
	if err != nil {
		return err
	}

	usage, err := getClusterUsage(param.Cluster)
	if err != nil {
		return err
	}

	*cr = ClusterRegistration{
		Cluster:                      param.Cluster,
		ClusterTemplateRepoLocalPath: param.ClusterTemplateRepoLocalPath,
		TenantConfigRepoLocalPath:    param.TenantConfigRepoLocalPath,
		RepoURL:                      param.RepoURL,
		CaBundle:                     param.CaBundle,
		Usage:                        usage,
		HostCluster:                  hostCluster,
		Vcluster:                     vcluster,
		Runtime:                      runtime,
		Traefik:                      param.Traefik,
		NautesConfigs:                param.Configs.Nautes,
		SecretConfigs:                param.Configs.Secret,
		OauthConfigs:                 param.Configs.OAuth,
		GitConfigs:                   param.Configs.Git,
	}

	return nil
}

func (cr *ClusterRegistration) getPojectPipelineItems(param *ClusterRegistrationParam, hostCluster *resourcev1alpha1.Cluster) ([]*ProjectPipelineItem, error) {
	var hostClusterName string
	var httpsNodePort int

	if hostCluster == nil {
		return nil, fmt.Errorf("the host cluster is not empty when get project pipeline configuration")
	}

	hostClusterName = hostCluster.Name
	httpsNodePort, err := cr.GetTraefikNodePortToHostCluster(param.TenantConfigRepoLocalPath, hostClusterName)
	if err != nil {
		return nil, err
	}

	ingressFilePath := concatTektonDashborardFilePath(param.TenantConfigRepoLocalPath, hostClusterName)
	_, err = os.Stat(ingressFilePath)
	if err != nil && os.IsNotExist(err) {
		return nil, nil
	}

	bytes, err := readIngressFileContent(ingressFilePath)
	if err != nil {
		return nil, err
	}

	ingresses, err := parseIngresses(bytes)
	if err != nil {
		return nil, err
	}

	projectPipelineItems := ConvertProjectPipeline(ingresses, hostClusterName, httpsNodePort)

	return projectPipelineItems, nil
}

func readIngressFileContent(path string) (string, error) {
	bytes, err := readFile(path)
	if err != nil {
		return "", err
	}
	if len(bytes) == 0 || string(bytes) == "\n" {
		return "", nil
	}

	return string(bytes), nil
}

func ConvertProjectPipeline(ingresses []Ingress, hostClusterName string, httpsNodePort int) []*ProjectPipelineItem {
	projectPipelineItems := make([]*ProjectPipelineItem, 0)
	for _, ingress := range ingresses {
		re := regexp.MustCompile(`^(.*?)-tekton-dashborard$`)
		match := re.FindStringSubmatch(ingress.Metadata.Name)
		name := match[1]
		projectPipelineItems = append(projectPipelineItems, &ProjectPipelineItem{
			Name: name,
			TektonConfig: &TektonConfig{
				URL:           fmt.Sprintf("https://%s:%d", ingress.Spec.TLS[0].Hosts[0], httpsNodePort),
				Host:          ingress.Spec.Rules[0].Host,
				HttpsNodePort: httpsNodePort,
			},
			HostClusterName: hostClusterName,
		})
	}

	return projectPipelineItems
}

func GetArgocdConfig(cr *ClusterRegistration, cluster *resourcev1alpha1.Cluster, param *ClusterRegistrationParam) (*ArgocdConfig, error) {
	var config = &ArgocdConfig{}
	var err error

	switch {
	case IsVirtual(cluster):
		// Virtual cluster argocd port uses the traefik of the host cluster
		httpsNodePort, err := cr.GetTraefikNodePortToHostCluster(param.TenantConfigRepoLocalPath, cluster.Spec.HostCluster)
		if err != nil {
			return nil, fmt.Errorf("failed to get host cluster %s traefik https NodePort, please check if the host cluster exists", cluster.Spec.HostCluster)
		}

		config.Host, err = GetArgoCDHost(param, cluster.Spec.ApiServer)
		if err != nil {
			return nil, err
		}
		config.URL = fmt.Sprintf("https://%s:%d", config.Host, httpsNodePort)
	case IsPhysical(cluster):
		config.Host, err = GetArgoCDHost(param, cluster.Spec.ApiServer)
		if err != nil {
			return nil, err
		}

		if param.Traefik != nil && param.Traefik.HttpsNodePort != "" {
			config.URL = fmt.Sprintf("https://%s:%s", config.Host, param.Traefik.HttpsNodePort)
		}
	}

	argocdProject := fmt.Sprintf("%s-%s", cluster.Name, _RuntimeProjectSuffix)
	config.Project = argocdProject

	return config, nil
}

func GetTektonConfig(cr *ClusterRegistration, cluster *resourcev1alpha1.Cluster, hostCluster *HostCluster, param *ClusterRegistrationParam) (*TektonConfig, error) {
	config := &TektonConfig{}

	switch {
	case IsVirtualProjectPipelineRuntime(cluster):
		httpsNodePort, err := cr.GetTraefikNodePortToHostCluster(param.TenantConfigRepoLocalPath, cluster.Spec.HostCluster)
		if err != nil {
			return nil, fmt.Errorf("failed to get host cluster %s tarefik https NodePort, please check if the host cluster exists", cluster.Spec.HostCluster)
		}

		oauthURL := fmt.Sprintf("https://auth.%s.%s:%d", hostCluster.Name, hostCluster.PrimaryDomain, httpsNodePort)

		config.OAuthURL = oauthURL
		config.HttpsNodePort = httpsNodePort

		tektonHost, err := GetTektonHost(param, cluster.Spec.ApiServer)
		if err != nil {
			return nil, err
		}
		config.Host = tektonHost
		config.URL = fmt.Sprintf("https://%s:%d", tektonHost, httpsNodePort)
	case IsPhysicalProjectPipelineRuntime(cluster):
		tektonHost, err := GetTektonHost(param, cluster.Spec.ApiServer)
		if err != nil {
			return nil, err
		}
		config.Host = tektonHost

		if param.Traefik != nil && param.Traefik.HttpsNodePort != "" {
			oauthURL := fmt.Sprintf("https://auth.%s.%s:%s/oauth2/callback", cluster.Name, cluster.Spec.PrimaryDomain, param.Traefik.HttpsNodePort)
			config.OAuthURL = oauthURL
			config.URL = fmt.Sprintf("https://%s:%s", tektonHost, param.Traefik.HttpsNodePort)
		}
	}

	return config, nil
}

func (cr *ClusterRegistration) getRuntime(param *ClusterRegistrationParam, cluster *resourcev1alpha1.Cluster, hostCluster *HostCluster) (*Runtime, error) {
	if cluster == nil {
		return nil, nil
	}

	argocdConfig, err := GetArgocdConfig(cr, cluster, param)
	if err != nil {
		return nil, err
	}

	if cluster.Spec.PrimaryDomain == "" {
		if !utilstrings.IsIPPortURL(cluster.Spec.ApiServer) {
			cluster.Spec.PrimaryDomain = cluster.Spec.ApiServer
		} else {
			hostClusterIP, err := utilstrings.ParseUrl(cluster.Spec.ApiServer)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve cluster server, err: %s", err)
			}
			cluster.Spec.PrimaryDomain = fmt.Sprintf("%s.%s", hostClusterIP, _NIPDomainSuffix)
		}
	}

	tektonConfig, err := GetTektonConfig(cr, cluster, hostCluster, param)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		Name:          fmt.Sprintf("%s-%s", cluster.Name, _RuntimeSuffix),
		ClusterName:   cluster.Name,
		Type:          string(cluster.Spec.ClusterType),
		PrimaryDomain: cluster.Spec.PrimaryDomain,
		MountPath:     cluster.Name,
		ApiServer:     cluster.Spec.ApiServer,
		ArgocdConfig:  argocdConfig,
		TektonConfig:  tektonConfig,
	}, nil
}

func (cr *ClusterRegistration) getHostCluster(param *ClusterRegistrationParam, cluster *resourcev1alpha1.Cluster, tenantName string) (*HostCluster, error) {
	if cluster == nil {
		return nil, nil
	}

	primaryDomain := cluster.Spec.PrimaryDomain
	if primaryDomain == "" {
		if !utilstrings.IsIPPortURL(cluster.Spec.ApiServer) {
			primaryDomain = cluster.Spec.ApiServer
		} else {
			hostClusterIP, err := utilstrings.ParseUrl(cluster.Spec.ApiServer)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve cluster server, err: %s", err)
			}
			primaryDomain = fmt.Sprintf("%s.%s", hostClusterIP, _NIPDomainSuffix)
		}
	}

	projectPipelineItems, err := cr.getPojectPipelineItems(param, cluster)
	if err != nil {
		return nil, err
	}

	return &HostCluster{
		Name:                 cluster.Name,
		ApiServer:            cluster.Spec.ApiServer,
		ArgocdProject:        tenantName,
		PrimaryDomain:        primaryDomain,
		ProjectPipelineItems: projectPipelineItems,
	}, nil
}

func (cr *ClusterRegistration) getVcluster(param *ClusterRegistrationParam, cluster, hostCluster *resourcev1alpha1.Cluster) (*Vcluster, error) {
	if cluster == nil {
		return nil, nil
	}

	vcluster := &Vcluster{
		Name:      cluster.Name,
		ApiServer: cluster.Spec.ApiServer,
		Namespace: cluster.Name,
	}

	// Get and set the HTTPS port for vcluster.
	if param.Vcluster != nil && param.Vcluster.HttpsNodePort != "" {
		vcluster.HttpsNodePort = param.Vcluster.HttpsNodePort
	} else {
		if vcluster.ApiServer == "" {
			return nil, fmt.Errorf("the apiserver of vcluster %s is not empty", vcluster.Name)
		}

		port, err := utilstrings.ExtractPortFromURL(vcluster.ApiServer)
		if err != nil {
			return nil, fmt.Errorf("failed to automatically obtain vcluster host, err: %w", err)
		}
		vcluster.HttpsNodePort = port
	}

	vcluster.HostCluster = &HostCluster{
		Name:          hostCluster.Name,
		ApiServer:     hostCluster.Spec.ApiServer,
		ArgocdProject: param.Configs.Nautes.TenantName,
	}

	// Set vcluster tls IP address.
	hostClusterIP, err := utilstrings.ParseUrl(hostCluster.Spec.ApiServer)
	if err != nil {
		return nil, err
	}
	vcluster.TLSSan = hostClusterIP

	return vcluster, nil
}

func GetArgoCDHost(param *ClusterRegistrationParam, apiServer string) (string, error) {
	if param.ArgocdHost != "" {
		return param.ArgocdHost, nil
	}

	clusterIP, err := utilstrings.ParseUrl(apiServer)
	if err != nil {
		return "", fmt.Errorf("argocd host not filled in and automatic parse host cluster IP failed, err: %v", err)
	}

	return generateNipHost(_ArgocdPrefix, param.Cluster.Name, clusterIP), nil
}

func GetTektonHost(param *ClusterRegistrationParam, apiServer string) (string, error) {
	if param.TektonHost != "" {
		return param.TektonHost, nil
	}

	clusterIP, err := utilstrings.ParseUrl(apiServer)
	if err != nil {
		return "", fmt.Errorf("tekton host not filled in and automatic parse host cluster IP failed, err: %s", err)
	}

	return generateNipHost(_TektonPrefix, param.Cluster.Name, clusterIP), nil
}

func (cr *ClusterRegistration) GetArgocdURL() (string, error) {
	if cr.Runtime.ArgocdConfig == nil {
		return "", fmt.Errorf("argocd config is empty")
	}

	url := fmt.Sprintf("%s/%s", cr.Runtime.ArgocdConfig.URL, _ArgocdOAuthSuffix)

	return url, nil
}

func (cr *ClusterRegistration) GetTektonOAuthURL() (string, error) {
	if cr.Runtime.TektonConfig != nil && cr.Runtime.TektonConfig.OAuthURL != "" {
		url := fmt.Sprintf("%s/%s", cr.Runtime.TektonConfig.OAuthURL, _TektonOAuthSuffix)

		return url, nil
	}

	return "", nil
}

const (
	_TraefikAppFile = "production/traefik-app.yaml"
)

func (cr *ClusterRegistration) GetTraefikNodePortToHostCluster(tenantLocalPath, hostClusterName string) (int, error) {
	traefikFilePath := fmt.Sprintf("%s/%s/%s", concatHostClustesrDir(tenantLocalPath), hostClusterName, _TraefikAppFile)
	app, err := parseArgocdApplication(traefikFilePath)
	if err != nil {
		return 0, err
	}
	if app == nil {
		return 0, nil
	}

	httpsNodePort, err := getTraefikHttpsNodePort(app)
	if err != nil {
		return 0, err
	}

	return httpsNodePort, nil
}

func GetHostClusterFromTenantConfigFile(tenantConfigRepoLocalPath, hostClusterName, tenantName string) (*resourcev1alpha1.Cluster, error) {
	clusterFileName := fmt.Sprintf("%s/%s.yaml", concatClustersDir(tenantConfigRepoLocalPath), hostClusterName)
	if _, err := os.Stat(clusterFileName); os.IsNotExist(err) {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("the host cluster %s is invalid or does not exist", hostClusterName)
		}
		return nil, err
	}

	cluster, err := parseCluster(clusterFileName)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func getClusterUsage(cluster *resourcev1alpha1.Cluster) (ClusterUsage, error) {
	switch {
	case IsPhysicalDeploymentRuntime(cluster):
		return _PhysicalDeploymentRuntime, nil
	case IsVirtualDeploymentRuntime(cluster):
		return _VirtualDeploymentRuntime, nil
	case IsPhysicalProjectPipelineRuntime(cluster):
		return _PhysicalProjectPipelineRuntime, nil
	case IsVirtualProjectPipelineRuntime(cluster):
		return _VirtualProjectPipelineRuntime, nil
	case IsHostCluser(cluster):
		return _HostCluster, nil
	default:
		return "", fmt.Errorf("the cluster is null and cannot determine the type")
	}
}
