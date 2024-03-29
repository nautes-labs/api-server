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
	utilstring "github.com/nautes-labs/api-server/util/string"
	resourcev1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
	nautesconfigs "github.com/nautes-labs/pkg/pkg/nautesconfigs"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CodeRepoUsecase struct {
	log                    *log.Helper
	codeRepo               CodeRepo
	secretRepo             Secretrepo
	nodestree              nodestree.NodesTree
	config                 *nautesconfigs.Config
	resourcesUsecase       *ResourcesUsecase
	codeRepoBindingUsecase *CodeRepoBindingUsecase
	client                 client.Client
	wg                     sync.WaitGroup
}

type CodeRepoData struct {
	Name string
	Spec resourcev1alpha1.CodeRepoSpec
}

func NewCodeRepoUsecase(logger log.Logger, codeRepo CodeRepo, secretRepo Secretrepo, nodestree nodestree.NodesTree, config *nautesconfigs.Config, resourcesUsecase *ResourcesUsecase, codeRepoBindingUsecase *CodeRepoBindingUsecase, client client.Client) *CodeRepoUsecase {
	codeRepoUsecase := &CodeRepoUsecase{
		log:                    log.NewHelper(log.With(logger)),
		codeRepo:               codeRepo,
		secretRepo:             secretRepo,
		nodestree:              nodestree,
		config:                 config,
		resourcesUsecase:       resourcesUsecase,
		codeRepoBindingUsecase: codeRepoBindingUsecase,
		client:                 client,
	}
	nodestree.AppendOperators(codeRepoUsecase)
	return codeRepoUsecase
}

type ProjectCodeRepo struct {
	CodeRepo *resourcev1alpha1.CodeRepo
	Project  *Project
}

func (c *CodeRepoUsecase) GetCodeRepo(ctx context.Context, codeRepoName, productName string) (*ProjectCodeRepo, error) {
	pid := fmt.Sprintf("%s/%s", productName, codeRepoName)
	project, err := c.codeRepo.GetCodeRepo(ctx, pid)
	if err != nil {
		return nil, err
	}

	if project != nil {
		resourceName := fmt.Sprintf("%s%d", RepoPrefix, int(project.ID))
		codeRepoName = resourceName
	}

	node, err := c.resourcesUsecase.Get(ctx, nodestree.CodeRepo, productName, c, func(nodes nodestree.Node) (string, error) {
		if project != nil {
			resourceName := fmt.Sprintf("%s%d", RepoPrefix, int(project.ID))
			return resourceName, nil
		}

		return "", fmt.Errorf("failed to get resource name")
	})

	if err != nil {
		return nil, err
	}

	codeRepo, ok := node.Content.(*resourcev1alpha1.CodeRepo)
	if !ok {
		return nil, fmt.Errorf("failed to get codeRepo %s, err: the resource is wrong type", codeRepoName)
	}

	err = c.convertProductToGroupName(ctx, codeRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to get codeRepo %s, err: unable to convert product name", codeRepoName)
	}

	return &ProjectCodeRepo{
		Project:  project,
		CodeRepo: codeRepo,
	}, nil
}

func (c *CodeRepoUsecase) ListCodeRepos(ctx context.Context, productName string) ([]*ProjectCodeRepo, error) {
	nodes, err := c.resourcesUsecase.List(ctx, productName, c)
	if err != nil {
		return nil, err
	}

	codeRepoNodes := nodestree.ListsResourceNodes(*nodes, nodestree.CodeRepo)

	items := make([]*ProjectCodeRepo, 0)
	for _, node := range codeRepoNodes {
		projectCodeRepo := c.getProjectByNode(ctx, node)
		if projectCodeRepo != nil {
			items = append(items, projectCodeRepo)
		}
	}

	return items, nil
}

func (c *CodeRepoUsecase) SaveCodeRepo(ctx context.Context, options *BizOptions, data *CodeRepoData, gitOptions *GitCodeRepoOptions) error {
	group, err := c.codeRepo.GetGroup(ctx, options.ProductName)
	if err != nil {
		return err
	}
	data.Spec.Product = fmt.Sprintf("%s%d", ProductPrefix, int(group.ID))

	project, err := c.saveRepository(ctx, group, options.ResouceName, gitOptions)
	if err != nil {
		return err
	}
	pid := int(project.ID)
	codeRepoName := fmt.Sprintf("%s%d", RepoPrefix, int(project.ID))
	data.Name = codeRepoName

	resourceOptions := &resourceOptions{
		resourceName:      options.ResouceName,
		resourceKind:      nodestree.CodeRepo,
		productName:       options.ProductName,
		insecureSkipCheck: options.InsecureSkipCheck,
		operator:          c,
	}
	err = c.resourcesUsecase.Save(ctx, resourceOptions, data)
	if err != nil {
		return err
	}

	_, err = c.saveDeployKey(ctx, pid, true)
	if err != nil {
		return err
	}

	_, err = c.saveDeployKey(ctx, pid, false)
	if err != nil {
		return err
	}

	projectAccessToken, err := c.saveAccessToken(ctx, pid)
	if err != nil {
		return err
	}

	if err := c.removeInvalidProjectAccessTokens(ctx, pid, projectAccessToken); err != nil {
		return err
	}

	nodes, err := c.resourcesUsecase.loadDefaultProjectNodes(ctx, options.ProductName)
	if err != nil {
		return err
	}

	err = c.codeRepoBindingUsecase.refreshAuthorization(ctx, nodes)
	if err != nil {
		return err
	}

	return nil
}

func (c *CodeRepoUsecase) DeleteCodeRepo(ctx context.Context, options *BizOptions) error {
	group, err := c.codeRepo.GetGroup(ctx, options.ProductName)
	if err != nil {
		return err
	}

	projectPath := fmt.Sprintf("%s/%s", group.Path, options.ResouceName)
	project, err := c.codeRepo.GetCodeRepo(ctx, projectPath)
	if err != nil {
		if commonv1.IsProjectNotFound(err) {
			return nil
		}
		return err
	}
	pid := int(project.ID)
	codeRepoName := fmt.Sprintf("%s%d", RepoPrefix, int(project.ID))

	resourceOptions := &resourceOptions{
		resourceKind:      nodestree.CodeRepo,
		productName:       options.ProductName,
		insecureSkipCheck: options.InsecureSkipCheck,
		operator:          c,
	}

	nodes, err := c.resourcesUsecase.loadDefaultProjectNodes(ctx, options.ProductName)
	if err != nil {
		return err
	}

	codeRepo, err := c.getCodeRepoByName(ctx, nodes, codeRepoName)
	if err != nil {
		return err
	}

	err = c.codeRepoBindingUsecase.refreshAuthorization(ctx, nodes, codeRepo.Name)
	if err != nil {
		return err
	}

	if err := c.resourcesUsecase.Delete(ctx, resourceOptions, func(nodes nodestree.Node) (string, error) {
		if project != nil {
			return codeRepoName, nil
		}

		return "", fmt.Errorf("failed to get resource name")
	}); err != nil {
		return err
	}

	if err = c.deleteReadOnlySecret(ctx, pid); err != nil {
		return err
	}

	if err = c.deleteReadWriteSecret(ctx, pid); err != nil {
		return err
	}

	if err = c.deleteAccessTokenSecret(ctx, pid); err != nil {
		return err
	}

	if err := c.deleteDeployKeys(ctx, pid); err != nil {
		return err
	}

	if err = c.codeRepo.DeleteCodeRepo(ctx, pid); err != nil {
		return err
	}

	return nil
}

func (c *CodeRepoUsecase) getProjectByNode(ctx context.Context, node *nodestree.Node) *ProjectCodeRepo {
	codeRepo, ok := node.Content.(*resourcev1alpha1.CodeRepo)
	if !ok {
		return nil
	}

	err := c.convertProductToGroupName(ctx, codeRepo)
	if err != nil {
		return nil
	}

	pid, err := utilstring.ExtractNumber(RepoPrefix, codeRepo.Name)
	if err != nil {
		return nil
	}

	project, err := c.codeRepo.GetCodeRepo(ctx, pid)
	if err != nil {
		return nil
	}

	return &ProjectCodeRepo{Project: project, CodeRepo: codeRepo}
}

func (c *CodeRepoUsecase) getCodeRepoByName(ctx context.Context, nodes *nodestree.Node, codeRepoName string) (*resourcev1alpha1.CodeRepo, error) {
	node := c.nodestree.GetNode(nodes, nodestree.CodeRepo, codeRepoName)
	if node == nil {
		return nil, commonv1.ErrorNodeNotFound("failed to get coderepo of repository %s", codeRepoName)
	}

	codeRepo, ok := node.Content.(*resourcev1alpha1.CodeRepo)
	if !ok {
		return nil, fmt.Errorf(" type found for %s node", node.Name)
	}

	return codeRepo, nil
}

func (c *CodeRepoUsecase) convertProductToGroupName(ctx context.Context, codeRepo *resourcev1alpha1.CodeRepo) error {
	if codeRepo.Spec.Product == "" {
		return fmt.Errorf("the product field value of coderepo %s should not be empty", codeRepo.Spec.RepoName)
	}

	groupName, err := c.resourcesUsecase.ConvertProductToGroupName(ctx, codeRepo.Spec.Product)
	if err != nil {
		return err
	}

	codeRepo.Spec.Product = groupName

	return nil
}

func nodeToCodeRepo(node *nodestree.Node) (*resourcev1alpha1.CodeRepo, error) {
	r, ok := node.Content.(*resourcev1alpha1.CodeRepo)
	if !ok {
		return nil, fmt.Errorf("failed to get instance when get %s coderepo", node.Name)
	}

	return r, nil
}

func (c *CodeRepoUsecase) saveRepository(ctx context.Context, group *Group, resourceName string, gitOptions *GitCodeRepoOptions) (*Project, error) {
	pid := fmt.Sprintf("%s/%s", group.Path, resourceName)
	project, err := c.codeRepo.GetCodeRepo(ctx, pid)
	e := errors.FromError(err)
	if err != nil && e.Code != 404 {
		return nil, err
	}

	if err != nil && e.Code == 404 {
		project, err = c.codeRepo.CreateCodeRepo(ctx, int(group.ID), gitOptions)
		if err != nil {
			return nil, err
		}
	} else {
		project, err = c.codeRepo.UpdateCodeRepo(ctx, int(project.ID), gitOptions)
		if err != nil {
			return nil, err
		}
	}

	return project, nil
}

func (c *CodeRepoUsecase) getDeployKeyFromSecretRepo(ctx context.Context, pid int, permission string) (*DeployKeySecretData, error) {
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

func (c *CodeRepoUsecase) getAccessTokenFromSecretRepo(ctx context.Context, pid int) (*AccessTokenSecretData, error) {
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

// saveDeployKey is to store the deploykey to git and the key repository separately
// It takes in an argument pid of type int and an argument canPush of type bool,
// and returns a project deploykey and error.
func (c *CodeRepoUsecase) saveDeployKey(ctx context.Context, pid int, canPush bool) (*ProjectDeployKey, error) {
	permission := string(ReadOnly)
	if canPush {
		permission = string(ReadWrite)
	}

	secretData, err := c.getDeployKeyFromSecretRepo(ctx, pid, permission)
	if err != nil {
		ok := commonv1.IsSecretNotFound(err)
		if !ok {
			return nil, err
		}

		err := c.clearDeployKeyWithSameName(ctx, pid, permission)
		if err != nil {
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

func (c *CodeRepoUsecase) clearDeployKeyWithSameName(ctx context.Context, pid int, permission string) error {
	title := fmt.Sprintf("%s%d-%s", RepoPrefix, pid, permission)

	projectDeployKeys, err := GetAllDeployKeys(ctx, c.codeRepo, pid)
	if err != nil {
		return err
	}

	for _, projectDeployKey := range projectDeployKeys {
		if projectDeployKey.Title == title {
			err := c.codeRepo.DeleteDeployKey(ctx, pid, projectDeployKey.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *CodeRepoUsecase) removeInvalidDeploykey(ctx context.Context, pid int, projectDeployKeys ...*ProjectDeployKey) error {
	validProjectDeployKeys := make(map[int]bool, len(projectDeployKeys))
	for _, deployKey := range projectDeployKeys {
		validProjectDeployKeys[deployKey.ID] = true
	}

	deployKeysToDelete, err := c.filterProjectDeployKeysByIDList(ctx, pid, validProjectDeployKeys)
	if err != nil {
		return err
	}

	err = c.deleteSameTitleDeployKey(ctx, pid, deployKeysToDelete, projectDeployKeys)
	if err != nil {
		return err
	}

	return nil
}

func (c *CodeRepoUsecase) deleteSameTitleDeployKey(ctx context.Context, pid int, deployKeysToDelete []*ProjectDeployKey, projectDeployKeys []*ProjectDeployKey) error {
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

func (c *CodeRepoUsecase) filterProjectDeployKeysByIDList(ctx context.Context, pid int, idList map[int]bool) ([]*ProjectDeployKey, error) {
	keys, err := GetAllDeployKeys(ctx, c.codeRepo, pid)
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

func (c *CodeRepoUsecase) saveAccessToken(ctx context.Context, pid int) (*ProjectAccessToken, error) {
	accessTokenSecretData, err := c.getAccessTokenFromSecretRepo(ctx, pid)
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

func (c *CodeRepoUsecase) removeInvalidProjectAccessTokens(ctx context.Context, pid int, projectAccessTokens ...*ProjectAccessToken) error {
	validAccessTokens := make(map[int]bool, len(projectAccessTokens))
	for _, projectAccessToken := range projectAccessTokens {
		validAccessTokens[projectAccessToken.ID] = true
	}

	accessTokensToDelete, err := c.filterAccessTokensByIDList(ctx, pid, validAccessTokens)
	if err != nil {
		return err
	}

	err = c.deleteSameNameAccessToken(accessTokensToDelete, projectAccessTokens, ctx, pid)
	if err != nil {
		return err
	}

	return nil
}

func (c *CodeRepoUsecase) deleteSameNameAccessToken(accessTokensToDelete []*ProjectAccessToken, projectAccessTokens []*ProjectAccessToken, ctx context.Context, pid int) error {
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

func (c *CodeRepoUsecase) filterAccessTokensByIDList(ctx context.Context, pid int, idList map[int]bool) ([]*ProjectAccessToken, error) {
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

type listFunc func(ctx context.Context, pid int, opts *ListOptions) ([]interface{}, error)

func getAllItems(ctx context.Context, pid int, opts *ListOptions, getAllFunc listFunc) ([]interface{}, error) {
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
	items, err := getAllItems(ctx, pid, opts, c.ListAccessTokens)
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

func GetAllDeployKeys(ctx context.Context, codeRepo CodeRepo, pid int) ([]*ProjectDeployKey, error) {
	opts := &ListOptions{
		Page:    1,
		PerPage: 10,
	}
	items, err := getAllItems(ctx, pid, opts, ListDeployKeys(ctx, codeRepo, pid, opts))
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

func ListDeployKeys(ctx context.Context, codeRepo CodeRepo, pid int, opts *ListOptions) listFunc {
	return func(ctx context.Context, pid int, opts *ListOptions) ([]interface{}, error) {
		keys, err := codeRepo.ListDeployKeys(ctx, pid, opts)
		if err != nil {
			return nil, err
		}

		var interfaceSlice []interface{}
		for _, v := range keys {
			interfaceSlice = append(interfaceSlice, interface{}(v))
		}

		return interfaceSlice, nil
	}
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
	err = c.secretRepo.SaveDeployKey(ctx, getCodeRepoResourceName(pid), string(privateKey), DefaultUser, permission, extendKVs)
	if err != nil {
		return nil, err
	}

	return projectDeployKey, nil
}

func (c *CodeRepoUsecase) createAccessTokenToGitAndSecretRepo(ctx context.Context, pid int) (*ProjectAccessToken, error) {
	name := AccessTokenName
	scopes := []string{string(APIPermission)}
	accessLevel := AccessLevelValue(Maintainer)
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
	err = c.secretRepo.SaveProjectAccessToken(ctx, getCodeRepoResourceName(pid), string(projectToken.Token), DefaultUser, name, extendKVs)
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

	if val.Spec.Webhook != nil && val.Spec.Webhook.Events == nil {
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

func (c *CodeRepoUsecase) UpdateNode(node *nodestree.Node, data interface{}) (*nodestree.Node, error) {
	val, ok := data.(*CodeRepoData)
	if !ok {
		return nil, fmt.Errorf("failed to get conderepo %s data when updating node", node.Name)
	}

	if len(val.Spec.Webhook.Events) == 0 {
		val.Spec.Webhook.Events = make([]string, 0)
	}

	codeRepo, ok := node.Content.(*resourcev1alpha1.CodeRepo)
	if !ok {
		return nil, fmt.Errorf("wrong type found for %s node when checking CodeRepo type", node.Name)
	}

	codeRepo.Spec = val.Spec
	codeRepoProvider, err := getCodeRepoProvider(c.client, c.config.Nautes.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get the infomation of the codeRepoProvider list when creating node")
	}

	if codeRepoProvider != nil {
		codeRepo.Spec.CodeRepoProvider = codeRepoProvider.Name
	}

	node.Content = codeRepo
	return node, nil
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

	if codeRepo.Spec.Webhook != nil {
		err = nodestree.CheckGitHooks(string(gitType), codeRepo.Spec.Webhook.Events)
		if err != nil {
			return true, err
		}
	}

	if codeRepo.Spec.Project != "" {
		ok = nodestree.IsResourceExist(options, codeRepo.Spec.Project, nodestree.Project)
		if !ok {
			return true, fmt.Errorf(_ResourceDoesNotExistOrUnavailable, nodestree.Project,
				codeRepo.Spec.Project, nodestree.CodeRepo, codeRepo.Spec.RepoName, CodeReposSubDir+"/"+codeRepo.ObjectMeta.Name)
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

func (c *CodeRepoUsecase) deleteDeployKeys(ctx context.Context, pid int) error {
	projectDeployKeys, err := GetAllDeployKeys(ctx, c.codeRepo, pid)
	if err != nil {
		return err
	}
	for _, key := range projectDeployKeys {
		if err := c.codeRepo.DeleteDeployKey(ctx, pid, key.ID); err != nil {
			return err
		}
	}

	return nil
}

func (c *CodeRepoUsecase) deleteReadOnlySecret(ctx context.Context, pid int) error {
	err := c.secretRepo.DeleteSecret(ctx, pid, DefaultUser, string(ReadOnly))
	if err != nil {
		return err
	}
	return nil
}

func (c *CodeRepoUsecase) deleteReadWriteSecret(ctx context.Context, pid int) error {
	err := c.secretRepo.DeleteSecret(ctx, pid, DefaultUser, string(ReadWrite))
	if err != nil {
		return err
	}
	return nil
}

func (c *CodeRepoUsecase) deleteAccessTokenSecret(ctx context.Context, pid int) error {
	err := c.secretRepo.DeleteSecret(ctx, pid, DefaultUser, string(AccessTokenName))
	if err != nil {
		return err
	}
	return nil
}

type ListMatchOptions func(*resourcev1alpha1.CodeRepo) bool

func nodesToCodeRepoists(nodes nodestree.Node, options ...ListMatchOptions) ([]*resourcev1alpha1.CodeRepo, error) {
	var codeReposDir *nodestree.Node
	var resources []*resourcev1alpha1.CodeRepo
	var filteredRepos []*resourcev1alpha1.CodeRepo

	for _, childNode := range nodes.Children {
		if childNode.Name == CodeReposSubDir {
			codeReposDir = childNode
			break
		}
	}

	if codeReposDir == nil {
		return nil, fmt.Errorf("the %s directory is not exist", CodeReposSubDir)
	}

	for _, subNode := range codeReposDir.Children {
		for _, node := range subNode.Children {
			if node.Kind != nodestree.CodeRepo {
				continue
			}

			resource, err := nodeToCodeRepo(node)
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
