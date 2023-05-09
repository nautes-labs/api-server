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

package biz

import (
	"context"
	"fmt"
	"strconv"

	errors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	commonv1 "github.com/nautes-labs/api-server/api/common/v1"
	"github.com/nautes-labs/api-server/pkg/nodestree"
	utilkey "github.com/nautes-labs/api-server/util/key"
	resourcev1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
	nautesconfigs "github.com/nautes-labs/pkg/pkg/nautesconfigs"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	_CodeRepoKind    = "CodeRepo"
	_CodeReposSubDir = "code-repos"
)

type CodeRepoUsecase struct {
	log              *log.Helper
	codeRepo         CodeRepo
	secretRepo       Secretrepo
	nodestree        nodestree.NodesTree
	config           *nautesconfigs.Config
	resourcesUsecase *ResourcesUsecase
	client           client.Client
}

type CodeRepoData struct {
	Name string
	Spec resourcev1alpha1.CodeRepoSpec
}

func NewCodeRepoUsecase(logger log.Logger, codeRepo CodeRepo, secretRepo Secretrepo, nodestree nodestree.NodesTree, config *nautesconfigs.Config, resourcesUsecase *ResourcesUsecase, client client.Client) *CodeRepoUsecase {
	codeRepoUsecase := &CodeRepoUsecase{log: log.NewHelper(log.With(logger)), codeRepo: codeRepo, secretRepo: secretRepo, nodestree: nodestree, config: config, resourcesUsecase: resourcesUsecase, client: client}
	nodestree.AppendOperators(codeRepoUsecase)
	return codeRepoUsecase
}

func (c *CodeRepoUsecase) GetCodeRepo(ctx context.Context, codeRepoName, productName string) (*resourcev1alpha1.CodeRepo, *Project, error) {
	pid := fmt.Sprintf("%s/%s", productName, codeRepoName)
	project, err := c.codeRepo.GetCodeRepo(ctx, pid)
	if err != nil {
		return nil, nil, err
	}

	if project != nil {
		resourceName := fmt.Sprintf("%s%d", RepoPrefix, int(project.Id))
		codeRepoName = resourceName
	}

	node, err := c.resourcesUsecase.Get(ctx, nodestree.CodeRepo, productName, c, func(nodes nodestree.Node) (string, error) {
		if project != nil {
			resourceName := fmt.Sprintf("%s%d", RepoPrefix, int(project.Id))
			return resourceName, nil
		}

		resourceName, err := c.getCodeRepoName(nodes, codeRepoName)
		if err != nil {
			return "", nil
		}

		return resourceName, nil
	})
	if err != nil {
		return nil, nil, err
	}

	codeRepo, err := c.nodeToResource(node)
	if err != nil {
		return nil, nil, err
	}

	err = c.convertProductToGroupName(ctx, codeRepo)
	if err != nil {
		return nil, nil, err
	}

	return codeRepo, project, nil
}

func (c *CodeRepoUsecase) convertProductToGroupName(ctx context.Context, codeRepo *resourcev1alpha1.CodeRepo) error {
	if codeRepo.Spec.Product == "" {
		return fmt.Errorf("the product field value of coderepo %s should not be empty", codeRepo.Spec.RepoName)
	}

	groupName, err := c.resourcesUsecase.convertProductToGroupName(ctx, codeRepo.Spec.Product)
	if err != nil {
		return err
	}

	codeRepo.Spec.Product = groupName

	return nil
}

func (c *CodeRepoUsecase) nodeToResource(node *nodestree.Node) (*resourcev1alpha1.CodeRepo, error) {
	r, ok := node.Content.(*resourcev1alpha1.CodeRepo)
	if !ok {
		return nil, fmt.Errorf("failed to get instance when get %s coderepo", node.Name)
	}

	return r, nil
}

type CodeRepoAndProject struct {
	CodeRepo *resourcev1alpha1.CodeRepo
	Project  *Project
}

func (c *CodeRepoUsecase) ListCodeRepos(ctx context.Context, productName string) ([]*CodeRepoAndProject, error) {
	nodes, err := c.resourcesUsecase.List(ctx, productName, c)
	if err != nil {
		return nil, err
	}

	codeRepos, err := c.nodesToLists(*nodes)
	if err != nil {
		return nil, err
	}

	var cps []*CodeRepoAndProject
	for _, codeRepo := range codeRepos {
		err = c.convertProductToGroupName(ctx, codeRepo)
		if err != nil {
			return nil, err
		}
		pid := fmt.Sprintf("%s/%s", productName, codeRepo.Spec.RepoName)
		project, err := c.codeRepo.GetCodeRepo(ctx, pid)
		if err != nil {
			return nil, err
		}
		cps = append(cps, &CodeRepoAndProject{
			CodeRepo: codeRepo,
			Project:  project,
		})
	}

	return cps, nil
}

func (c *CodeRepoUsecase) SaveCodeRepo(ctx context.Context, options *BizOptions, data *CodeRepoData, gitOptions *GitCodeRepoOptions) error {
	group, err := c.codeRepo.GetGroup(ctx, options.ProductName)
	if err != nil {
		return err
	}
	data.Spec.Product = fmt.Sprintf("%s%d", _ProductPrefix, int(group.Id))

	project, err := c.saveRepository(ctx, group, options.ResouceName, gitOptions)
	if err != nil {
		return err
	}
	pid := int(project.Id)
	data.Name = fmt.Sprintf("%s%d", RepoPrefix, int(project.Id))

	resourceOptions := &resourceOptions{
		resourceKind:      nodestree.CodeRepo,
		productName:       options.ProductName,
		insecureSkipCheck: options.InsecureSkipCheck,
		operator:          c,
	}
	err = c.resourcesUsecase.Save(ctx, resourceOptions, data)
	if err != nil {
		return err
	}

	projectReadWriteDeployKey, err := c.SaveDeployKey(ctx, pid, true)
	if err != nil {
		return err
	}

	projectReadOnlyDeployKey, err := c.SaveDeployKey(ctx, pid, false)
	if err != nil {
		return err
	}

	if err := c.RemoveInvalidDeploykey(ctx, pid, projectReadWriteDeployKey, projectReadOnlyDeployKey); err != nil {
		return err
	}

	projectAccessToken, err := c.SaveAccessToken(ctx, pid)
	if err != nil {
		return err
	}

	if err := c.RemoveInvalidProjectAccessTokens(ctx, pid, projectAccessToken); err != nil {
		return err
	}

	return nil
}

func (c *CodeRepoUsecase) saveRepository(ctx context.Context, group *Group, resourceName string, gitOptions *GitCodeRepoOptions) (*Project, error) {
	pid := fmt.Sprintf("%s/%s", group.Path, resourceName)
	project, err := c.codeRepo.GetCodeRepo(ctx, pid)
	e := errors.FromError(err)
	if err != nil && e.Code != 404 {
		return nil, err
	}

	if err != nil && e.Code == 404 {
		project, err = c.codeRepo.CreateCodeRepo(ctx, int(group.Id), gitOptions)
		if err != nil {
			return nil, err
		}
	} else {
		project, err = c.codeRepo.UpdateCodeRepo(ctx, int(project.Id), gitOptions)
		if err != nil {
			return nil, err
		}
	}

	return project, nil
}

func (c *CodeRepoUsecase) GetDeployKeyFromSecretRepo(ctx context.Context, pid int, permission string) (*DeployKeySecretData, error) {
	repoName := fmt.Sprintf("%s%d", RepoPrefix, pid)
	gitType := c.config.Git.GitType
	secretPath := fmt.Sprintf("%s/%s/%s/%s", gitType, repoName, DefaultUser, permission)
	secretOptions := &SecretOptions{
		SecretPath:   secretPath,
		SecretEngine: SecretsGitEngine,
		SecretKey:    SecretsDeployKey,
	}

	deployKeySecretData, err := c.secretRepo.GetDeployKey(ctx, secretOptions)
	if err != nil {
		return nil, err
	}

	return deployKeySecretData, nil
}

func (c *CodeRepoUsecase) GetAccessTokenFromSecretRepo(ctx context.Context, pid int) (*AccessTokenSecretData, error) {
	repoName := fmt.Sprintf("%s%d", RepoPrefix, pid)
	gitType := c.config.Git.GitType
	secretPath := fmt.Sprintf("%s/%s/%s/%s", gitType, repoName, DefaultUser, AccessTokenName)
	secretOptions := &SecretOptions{
		SecretPath:   secretPath,
		SecretEngine: SecretsGitEngine,
		SecretKey:    SecretsAccessToken,
	}

	accessTokenSecretData, err := c.secretRepo.GetProjectAccessToken(ctx, secretOptions)
	if err != nil {
		return nil, err
	}

	return accessTokenSecretData, nil
}

// SaveDeployKey is to store the deploykey to git and the key repository separately
// It takes in an argument pid of type int and an argument canPush of type bool,
// and returns a project deploykey and error.
func (c *CodeRepoUsecase) SaveDeployKey(ctx context.Context, pid int, canPush bool) (*ProjectDeployKey, error) {
	permission := string(ReadOnly)
	if canPush {
		permission = string(ReadWrite)
	}

	secretData, err := c.GetDeployKeyFromSecretRepo(ctx, pid, permission)
	if err != nil {
		ok := commonv1.IsSecretNotFound(err)
		if !ok {
			return nil, err
		}

		projectDeployKey, err := c.saveDeployKeyToGitAndSecretRepo(ctx, pid, canPush, permission)
		if err != nil {
			return nil, err
		}

		return projectDeployKey, nil
	}

	projectDeployKey, err := c.codeRepo.GetDeployKey(ctx, pid, secretData.ID)
	if err != nil {
		ok := commonv1.IsDeploykeyNotFound(err)
		if !ok {
			return nil, err
		}

		projectDeployKey, err = c.saveDeployKeyToGitAndSecretRepo(ctx, pid, canPush, permission)
		if err != nil {
			return nil, err
		}

		return projectDeployKey, nil
	}

	if projectDeployKey.Key == secretData.Fingerprint {
		return projectDeployKey, nil
	}

	projectDeployKey, err = c.saveDeployKeyToGitAndSecretRepo(ctx, pid, canPush, permission)
	if err != nil {
		return nil, err
	}

	return projectDeployKey, nil
}

func (c *CodeRepoUsecase) RemoveInvalidDeploykey(ctx context.Context, pid int, projectDeployKeys ...*ProjectDeployKey) error {
	validProjectDeployKeys := make(map[int]bool, len(projectDeployKeys))
	for _, deployKey := range projectDeployKeys {
		validProjectDeployKeys[deployKey.ID] = true
	}

	deployKeysToDelete, err := c.FilterProjectDeployKeysByIDList(ctx, pid, validProjectDeployKeys)
	if err != nil {
		return err
	}

	for _, key := range deployKeysToDelete {
		for _, projectDeployKey := range projectDeployKeys {
			if key.Title == projectDeployKey.Title {
				err := c.codeRepo.DeleteDeployKey(ctx, pid, key.ID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *CodeRepoUsecase) FilterProjectDeployKeysByIDList(ctx context.Context, pid int, idList map[int]bool) ([]*ProjectDeployKey, error) {
	keys, err := c.GetAllDeployKeys(ctx, pid)
	if err != nil {
		return nil, err
	}

	var filteredDeployKeys []*ProjectDeployKey
	for _, key := range keys {
		if !idList[key.ID] {
			filteredDeployKeys = append(filteredDeployKeys, key)
		}
	}

	return filteredDeployKeys, nil
}

func (c *CodeRepoUsecase) SaveAccessToken(ctx context.Context, pid int) (*ProjectAccessToken, error) {
	accessTokenSecretData, err := c.GetAccessTokenFromSecretRepo(ctx, pid)
	if err != nil {
		ok := commonv1.IsSecretNotFound(err)
		if !ok {
			return nil, err
		}

		projectAccessToken, err := c.createAccessTokenToGitAndSecretRepo(ctx, pid)
		if err != nil {
			return nil, err
		}

		return projectAccessToken, nil
	}

	projectAccessToken, err := c.codeRepo.GetProjectAccessToken(ctx, pid, accessTokenSecretData.ID)
	if err != nil {
		ok := commonv1.IsAccesstokenNotFound(err)
		if !ok {
			return nil, err
		}

		projectAccessToken, err := c.createAccessTokenToGitAndSecretRepo(ctx, pid)
		if err != nil {
			return nil, err
		}

		return projectAccessToken, nil
	}

	return projectAccessToken, nil
}

func (c *CodeRepoUsecase) RemoveInvalidProjectAccessTokens(ctx context.Context, pid int, projectAccessTokens ...*ProjectAccessToken) error {
	validAccessTokens := make(map[int]bool, len(projectAccessTokens))
	for _, projectAccessToken := range projectAccessTokens {
		validAccessTokens[projectAccessToken.ID] = true
	}

	accessTokensToDelete, err := c.FilterAccessTokensByIDList(ctx, pid, validAccessTokens)
	if err != nil {
		return err
	}

	for _, token := range accessTokensToDelete {
		for _, validToken := range projectAccessTokens {
			if token.Name == validToken.Name {
				err := c.codeRepo.DeleteProjectAccessToken(ctx, pid, token.ID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *CodeRepoUsecase) FilterAccessTokensByIDList(ctx context.Context, pid int, idList map[int]bool) ([]*ProjectAccessToken, error) {
	allAccessTokens, err := c.GetAllAccessTokens(ctx, pid)
	if err != nil {
		return nil, err
	}

	var filteredAccessTokens []*ProjectAccessToken
	for _, token := range allAccessTokens {
		if !idList[token.ID] {
			filteredAccessTokens = append(filteredAccessTokens, token)
		}
	}

	return filteredAccessTokens, nil
}

func (c *CodeRepoUsecase) getAllItems(ctx context.Context, pid int, opts *ListOptions, getAllFunc func(ctx context.Context, pid int, opts *ListOptions) ([]interface{}, error)) ([]interface{}, error) {
	allItems := []interface{}{}
	var err error
	for {
		var items []interface{}
		items, err = getAllFunc(ctx, pid, opts)
		if err != nil {
			return nil, err
		}
		opts.Page += 1
		allItems = append(allItems, items...)
		if len(items) != opts.PerPage {
			break
		}
	}
	return allItems, err
}

func (c *CodeRepoUsecase) GetAllAccessTokens(ctx context.Context, pid int) ([]*ProjectAccessToken, error) {
	opts := &ListOptions{
		Page:    1,
		PerPage: 10,
	}
	items, err := c.getAllItems(ctx, pid, opts, c.ListAccessTokens)
	if err != nil {
		return nil, err
	}
	accessTokens := make([]*ProjectAccessToken, len(items))
	for i, item := range items {
		accessToken, ok := item.(*ProjectAccessToken)
		if !ok {
			return nil, fmt.Errorf("unexpected item type: %T", item)
		}
		accessTokens[i] = accessToken
	}

	return accessTokens, nil
}

func (c *CodeRepoUsecase) ListAccessTokens(ctx context.Context, pid int, opts *ListOptions) ([]interface{}, error) {
	tokens, err := c.codeRepo.ListAccessTokens(ctx, pid, opts)
	if err != nil {
		return nil, err
	}

	var interfaceSlice []interface{}
	for _, v := range tokens {
		interfaceSlice = append(interfaceSlice, interface{}(v))
	}

	return interfaceSlice, nil
}

func (c *CodeRepoUsecase) GetAllDeployKeys(ctx context.Context, pid int) ([]*ProjectDeployKey, error) {
	opts := &ListOptions{
		Page:    1,
		PerPage: 10,
	}
	items, err := c.getAllItems(ctx, pid, opts, c.ListDeployKeys)
	if err != nil {
		return nil, err
	}
	allDeployKeys := make([]*ProjectDeployKey, len(items))
	for i, item := range items {
		accessToken, ok := item.(*ProjectDeployKey)
		if !ok {
			return nil, fmt.Errorf("unexpected item type: %T", item)
		}
		allDeployKeys[i] = accessToken
	}

	return allDeployKeys, nil
}

func (c *CodeRepoUsecase) ListDeployKeys(ctx context.Context, pid int, opts *ListOptions) ([]interface{}, error) {
	keys, err := c.codeRepo.ListDeployKeys(ctx, pid, opts)
	if err != nil {
		return nil, err
	}

	var interfaceSlice []interface{}
	for _, v := range keys {
		interfaceSlice = append(interfaceSlice, interface{}(v))
	}

	return interfaceSlice, nil
}

func (c *CodeRepoUsecase) saveDeployKeyToGitAndSecretRepo(ctx context.Context, pid int, canPush bool, permission string) (*ProjectDeployKey, error) {
	title := fmt.Sprintf("%s%d-%s", RepoPrefix, pid, permission)

	publicKey, privateKey, err := utilkey.GenerateKeyPair(c.config.Git.DefaultDeployKeyType, title)
	if err != nil {
		return nil, err
	}

	projectDeployKey, err := c.codeRepo.SaveDeployKey(ctx, pid, title, canPush, publicKey)
	if err != nil {
		return nil, err
	}

	extendKVs := make(map[string]string)
	extendKVs[Fingerprint] = projectDeployKey.Key
	extendKVs[DeployKeyID] = strconv.Itoa(projectDeployKey.ID)
	err = c.secretRepo.SaveDeployKey(ctx, convertRepoName(pid), string(privateKey), DefaultUser, permission, extendKVs)
	if err != nil {
		return nil, err
	}

	return projectDeployKey, nil
}

func (c *CodeRepoUsecase) createAccessTokenToGitAndSecretRepo(ctx context.Context, pid int) (*ProjectAccessToken, error) {
	name := AccessTokenName
	scopes := []string{string(APIPermission)}
	accessLevel := AccessLevelValue(Developer)
	projectToken, err := c.codeRepo.CreateProjectAccessToken(ctx, pid, &CreateProjectAccessTokenOptions{
		Name:        &name,
		Scopes:      &scopes,
		AccessLevel: &accessLevel,
	})
	if err != nil {
		return nil, err
	}

	extendKVs := make(map[string]string)
	extendKVs[AccessTokenID] = strconv.Itoa(projectToken.ID)
	err = c.secretRepo.SaveProjectAccessToken(ctx, convertRepoName(pid), string(projectToken.Token), DefaultUser, name, extendKVs)
	if err != nil {
		return nil, err
	}

	return projectToken, nil
}

func (c *CodeRepoUsecase) CreateNode(path string, data interface{}) (*nodestree.Node, error) {
	val, ok := data.(*CodeRepoData)
	if !ok {
		return nil, fmt.Errorf("failed to creating specify node, the path: %s", path)
	}

	if len(val.Spec.Webhook.Events) == 0 {
		val.Spec.Webhook.Events = make([]string, 0)
	}

	codeRepo := &resourcev1alpha1.CodeRepo{
		TypeMeta: v1.TypeMeta{
			Kind:       nodestree.CodeRepo,
			APIVersion: resourcev1alpha1.GroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name: val.Name,
		},
		Spec: val.Spec,
	}

	resourceDirectory := fmt.Sprintf("%s/%s", path, "code-repos")
	resourcePath := fmt.Sprintf("%s/%s/%s.yaml", resourceDirectory, val.Name, val.Name)

	codeRepoProvider, err := getCodeRepoProvider(c.client, c.config.Nautes.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get the infomation of the codeRepoProvider list when creating node")
	}

	if codeRepoProvider != nil {
		codeRepo.Spec.CodeRepoProvider = codeRepoProvider.Name
	}

	return &nodestree.Node{
		Name:    val.Name,
		Path:    resourcePath,
		Content: codeRepo,
		Kind:    nodestree.CodeRepo,
		Level:   4,
	}, nil
}

func (c *CodeRepoUsecase) UpdateNode(resourceNode *nodestree.Node, data interface{}) (*nodestree.Node, error) {
	val, ok := data.(*CodeRepoData)
	if !ok {
		return nil, fmt.Errorf("failed to get conderepo %s data when updating node", resourceNode.Name)
	}

	if len(val.Spec.Webhook.Events) == 0 {
		val.Spec.Webhook.Events = make([]string, 0)
	}

	codeRepo, ok := resourceNode.Content.(*resourcev1alpha1.CodeRepo)
	if !ok {
		return nil, fmt.Errorf("failed to get coderepo %s when updating node", resourceNode.Name)
	}

	codeRepo.Spec = val.Spec
	codeRepoProvider, err := getCodeRepoProvider(c.client, c.config.Nautes.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get the infomation of the codeRepoProvider list when creating node")
	}

	if codeRepoProvider != nil {
		codeRepo.Spec.CodeRepoProvider = codeRepoProvider.Name
	}

	resourceNode.Content = codeRepo
	return resourceNode, nil
}

func (c *CodeRepoUsecase) CheckReference(options nodestree.CompareOptions, node *nodestree.Node, k8sClient client.Client) (bool, error) {
	if node.Kind != nodestree.CodeRepo {
		return false, nil
	}

	codeRepo, ok := node.Content.(*resourcev1alpha1.CodeRepo)
	if !ok {
		return true, fmt.Errorf("wrong type found for %s node when checking CodeRepo type", node.Name)
	}

	err := nodestree.CheckResourceSubdirectory(&options.Nodes, node)
	if err != nil {
		return true, err
	}

	productName := codeRepo.Spec.Product
	if productName != options.ProductName {
		return true, fmt.Errorf("the product name of resource %s does not match the current product name, expected product is %s, but now is %s", codeRepo.Spec.RepoName, options.ProductName, productName)
	}

	gitType := c.config.Git.GitType
	if gitType == "" {
		return true, fmt.Errorf("git type cannot be empty")
	}

	err = nodestree.CheckGitHooks(string(gitType), codeRepo.Spec.Webhook.Events)
	if err != nil {
		return true, err
	}

	if codeRepo.Spec.Project != "" {
		ok = nodestree.IsResourceExist(options, codeRepo.Spec.Project, nodestree.Project)
		if !ok {
			return true, fmt.Errorf(_ResourceDoesNotExistOrUnavailable, _ProjectKind,
				codeRepo.Spec.Project, _CodeRepoKind, codeRepo.Spec.RepoName, _CodeReposSubDir+"/"+codeRepo.ObjectMeta.Name)
		}
	}

	tenantAdminNamespace := c.config.Nautes.Namespace
	if tenantAdminNamespace == "" {
		return true, fmt.Errorf("tenant admin namspace cannot be empty")
	}

	objKey := client.ObjectKey{
		Namespace: tenantAdminNamespace,
		Name:      codeRepo.Spec.CodeRepoProvider,
	}

	err = k8sClient.Get(context.TODO(), objKey, &resourcev1alpha1.CodeRepoProvider{})
	if err != nil {
		return true, fmt.Errorf("codeRepoProvider %s is an invalid resource, err: %s", codeRepo.Spec.CodeRepoProvider, err)
	}

	return true, nil
}

func (c *CodeRepoUsecase) CreateResource(kind string) interface{} {
	if kind != nodestree.CodeRepo {
		return nil
	}

	return &resourcev1alpha1.CodeRepo{}
}

func (c *CodeRepoUsecase) DeleteCodeRepo(ctx context.Context, options *BizOptions) error {
	group, err := c.codeRepo.GetGroup(ctx, options.ProductName)
	if err != nil {
		return err
	}

	projectPath := fmt.Sprintf("%s/%s", group.Path, options.ResouceName)
	project, err := c.codeRepo.GetCodeRepo(ctx, projectPath)
	e := errors.FromError(err)
	if err != nil && e.Code != 404 {
		return err
	}

	resourceOptions := &resourceOptions{
		resourceKind:      nodestree.CodeRepo,
		productName:       options.ProductName,
		insecureSkipCheck: options.InsecureSkipCheck,
		operator:          c,
	}
	err = c.resourcesUsecase.Delete(ctx, resourceOptions, func(nodes nodestree.Node) (string, error) {
		if project != nil {
			resourceName := fmt.Sprintf("%s%d", RepoPrefix, int(project.Id))
			return resourceName, nil
		}

		resourceName, err := c.getCodeRepoName(nodes, options.ResouceName)
		if err != nil {
			return "", err
		}

		return resourceName, nil
	})
	if err != nil {
		return err
	}

	if err == nil {
		err = c.codeRepo.DeleteCodeRepo(ctx, int(project.Id))
		if err != nil {
			return err
		}

		err = c.secretRepo.DeleteSecret(ctx, int(project.Id))
		if err != nil {
			return err
		}
	}

	return nil
}

type ListMatchOptions func(*resourcev1alpha1.CodeRepo) bool

func (e *CodeRepoUsecase) nodesToLists(nodes nodestree.Node, options ...ListMatchOptions) ([]*resourcev1alpha1.CodeRepo, error) {
	var codeReposDir *nodestree.Node
	var resources []*resourcev1alpha1.CodeRepo
	var filteredRepos []*resourcev1alpha1.CodeRepo

	for _, childNode := range nodes.Children {
		if childNode.Name == _CodeReposSubDir {
			codeReposDir = childNode
			break
		}
	}

	if codeReposDir == nil {
		return nil, fmt.Errorf("the %s directory is not exist", _CodeReposSubDir)
	}

	for _, subNode := range codeReposDir.Children {
		for _, node := range subNode.Children {
			if node.Kind != nodestree.CodeRepo {
				continue
			}

			resource, err := e.nodeToResource(node)
			if err != nil {
				return nil, err
			}

			matching := false
			for _, fn := range options {
				if fn(resource) {
					filteredRepos = append(filteredRepos, resource)
					matching = true
					break
				}
			}

			if !matching {
				resources = append(resources, resource)
			}
		}
	}

	if len(filteredRepos) > 0 {
		return filteredRepos, nil
	}

	return resources, nil
}

func (c *CodeRepoUsecase) getCodeRepoName(nodes nodestree.Node, codeRepoName string) (string, error) {
	var resourceName string
	var codeReposDir *nodestree.Node

	for _, node := range nodes.Children {
		if node.Name == _CodeReposSubDir {
			codeReposDir = node
			break
		}
	}

	if codeReposDir == nil {
		return "", ErrorResourceNoFound
	}

	for _, subNode := range codeReposDir.Children {
		for _, node := range subNode.Children {
			if node.Kind != nodestree.CodeRepo {
				continue
			}

			name := nodestree.GetResourceValue(node.Content, "Spec", "RepoName")
			if name == codeRepoName {
				resourceName = node.Name
				break
			}

		}
	}

	if resourceName == "" {
		return "", ErrorResourceNoFound
	}

	return resourceName, nil
}

func getCodeRepoProvider(k8sClient client.Client, namespace string) (*resourcev1alpha1.CodeRepoProvider, error) {
	providers := &resourcev1alpha1.CodeRepoProviderList{}
	err := k8sClient.List(context.TODO(), providers, &client.ListOptions{
		Namespace: namespace,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get codeRepoProvider list, err: %w", err)
	}

	if providers.Items != nil && len(providers.Items) > 0 {
		return &providers.Items[0], nil
	}

	return nil, nil
}

func convertRepoName(id int) string {
	return fmt.Sprintf("repo-%d", id)
}
