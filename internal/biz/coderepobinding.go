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
	"regexp"
	"strconv"
	"sync"

	errors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	commonv1 "github.com/nautes-labs/api-server/api/common/v1"
	"github.com/nautes-labs/api-server/pkg/nodestree"
	utilstrings "github.com/nautes-labs/api-server/util/string"
	resourcev1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
	nautesconfigs "github.com/nautes-labs/pkg/pkg/nautesconfigs"
	"golang.org/x/sync/singleflight"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CodeRepoBindingUsecase struct {
	log              *log.Helper
	codeRepo         CodeRepo
	secretRepo       Secretrepo
	nodestree        nodestree.NodesTree
	resourcesUsecase *ResourcesUsecase
	config           *nautesconfigs.Config
	client           client.Client
	groupName        string
	lock             sync.RWMutex
	wg               sync.WaitGroup
	cacheStore       *CacheStore
}

type CacheStore struct {
	projectDeployKeyMap map[int]map[int]*ProjectDeployKey
}

type CodeRepoBindingData struct {
	Name string
	Spec resourcev1alpha1.CodeRepoBindingSpec
}

type applyDeploykeyFunc func(ctx context.Context, pid interface{}, deployKey int) error

func NewCodeRepoCodeRepoBindingUsecase(logger log.Logger, codeRepo CodeRepo, secretRepo Secretrepo, nodestree nodestree.NodesTree, resourcesUsecase *ResourcesUsecase, config *nautesconfigs.Config, client client.Client) *CodeRepoBindingUsecase {
	cacheStore := &CacheStore{
		projectDeployKeyMap: map[int]map[int]*ProjectDeployKey{},
	}
	codeRepoBindingUsecase := &CodeRepoBindingUsecase{
		log:              log.NewHelper(log.With(logger)),
		codeRepo:         codeRepo,
		secretRepo:       secretRepo,
		nodestree:        nodestree,
		resourcesUsecase: resourcesUsecase,
		config:           config,
		client:           client,
		cacheStore:       cacheStore,
	}
	nodestree.AppendOperators(codeRepoBindingUsecase)

	return codeRepoBindingUsecase
}

func (c *CodeRepoBindingUsecase) GetCodeRepoBinding(ctx context.Context, options *BizOptions) (*resourcev1alpha1.CodeRepoBinding, error) {
	node, err := c.resourcesUsecase.Get(ctx, nodestree.CodeRepoBinding, options.ProductName, c, func(nodes nodestree.Node) (string, error) {
		return options.ResouceName, nil
	})
	if err != nil {
		return nil, err
	}

	resource, err := c.nodeToResource(node)
	if err != nil {
		return nil, err
	}

	err = c.ConvertRuntime(ctx, resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (c *CodeRepoBindingUsecase) ConvertRuntime(ctx context.Context, resource *resourcev1alpha1.CodeRepoBinding) error {
	repoName, err := c.resourcesUsecase.ConvertCodeRepoToRepoName(ctx, resource.Spec.CodeRepo)
	if err != nil {
		return err
	}
	resource.Spec.CodeRepo = repoName

	groupName, err := c.resourcesUsecase.ConvertProductToGroupName(ctx, resource.Spec.Product)
	if err != nil {
		return err
	}
	resource.Spec.Product = groupName

	return nil
}

func (c *CodeRepoBindingUsecase) ListCodeRepoBindings(ctx context.Context, options *BizOptions) ([]*nodestree.Node, error) {
	nodes, err := c.resourcesUsecase.List(ctx, options.ProductName, c)
	if err != nil {
		return nil, err
	}

	codeRepoBindingNodes := nodestree.ListsResourceNodes(*nodes, nodestree.CodeRepoBinding)

	return codeRepoBindingNodes, nil
}

func (c *CodeRepoBindingUsecase) SaveCodeRepoBinding(ctx context.Context, options *BizOptions, data *CodeRepoBindingData) error {
	c.groupName = options.ProductName

	resourceOptions := &resourceOptions{
		resourceName:      options.ResouceName,
		resourceKind:      nodestree.CodeRepoBinding,
		productName:       options.ProductName,
		insecureSkipCheck: options.InsecureSkipCheck,
		operator:          c,
	}
	if err := c.resourcesUsecase.Save(ctx, resourceOptions, data); err != nil {
		return err
	}

	latestNodes, err := c.resourcesUsecase.GetNodes()
	if err != nil {
		return err
	}

	return c.refreshAuthorization(ctx, *latestNodes, data.Spec.CodeRepo)
}

func (c *CodeRepoBindingUsecase) DeleteCodeRepoBinding(ctx context.Context, options *BizOptions) error {
	resourceOptions := &resourceOptions{
		resourceName:      options.ResouceName,
		resourceKind:      nodestree.CodeRepoBinding,
		productName:       options.ProductName,
		insecureSkipCheck: options.InsecureSkipCheck,
		operator:          c,
	}

	nodes, err := c.resourcesUsecase.loadDefaultProjectNodes(ctx, options.ProductName)
	if err != nil {
		return err
	}
	node := c.nodestree.GetNode(nodes, nodestree.CodeRepoBinding, options.ResouceName)
	if node == nil {
		return fmt.Errorf("%s resource %s not found or invalid. Please check whether the resource exists under the default project", resourceOptions.resourceKind, options.ResouceName)
	}
	lastCodeRepoBinding, ok := node.Content.(*resourcev1alpha1.CodeRepoBinding)
	if !ok {
		return fmt.Errorf("resource type is inconsistent, please check if this resource %s is legal", options.ResouceName)
	}

	if err := c.resourcesUsecase.Delete(ctx, resourceOptions, func(nodes nodestree.Node) (string, error) {
		return resourceOptions.resourceName, nil
	}); err != nil {
		return err
	}

	nodes, err = c.resourcesUsecase.GetNodes()
	if err != nil {
		return err
	}

	return c.refreshAuthorization(ctx, *nodes, lastCodeRepoBinding.Spec.CodeRepo)
}

func (c *CodeRepoBindingUsecase) getCodeRepoBindings(nodes nodestree.Node, codeRepoName string) ([]*resourcev1alpha1.CodeRepoBinding, error) {
	codeRepoBindingNodes := nodestree.ListsResourceNodes(nodes, nodestree.CodeRepoBinding, func(node *nodestree.Node) bool {
		val, ok := node.Content.(*resourcev1alpha1.CodeRepoBinding)
		if !ok {
			return false
		}

		if val.Spec.CodeRepo == codeRepoName {
			return true
		}

		return false
	})

	var codeRepoBindings []*resourcev1alpha1.CodeRepoBinding
	for _, node := range codeRepoBindingNodes {
		codeRepoBinding, ok := node.Content.(*resourcev1alpha1.CodeRepoBinding)
		if ok {
			codeRepoBindings = append(codeRepoBindings, codeRepoBinding)
		}
	}

	return codeRepoBindings, nil
}

func (c *CodeRepoBindingUsecase) authorizeDeployKey(ctx context.Context, codeRepos []*resourcev1alpha1.CodeRepo, authorizationpid interface{}, permissions string) error {
	var repoIDs []int
	for _, repo := range codeRepos {
		pid, err := utilstrings.ExtractNumber(RepoPrefix, repo.Name)
		if err != nil {
			return err
		}

		repoIDs = append(repoIDs, pid)
	}

	if err := c.applyDeploykey(ctx, authorizationpid, permissions, repoIDs, func(ctx context.Context, authorizationpid interface{}, deployKeyID int) error {
		value, _ := authorizationpid.(int)

		projectDeployKey, err := c.codeRepo.EnableProjectDeployKey(ctx, value, deployKeyID)
		if err != nil {
			return err
		}

		if permissions != string(ReadWrite) {
			return nil
		}

		_, err = c.codeRepo.UpdateDeployKey(ctx, value, projectDeployKey.ID, projectDeployKey.Title, true)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (c *CodeRepoBindingUsecase) applyDeploykey(ctx context.Context, authorizationpid interface{}, permissions string, repoIDs []int, fn applyDeploykeyFunc) error {
	for _, repoID := range repoIDs {
		//Revoke in-product authorization, not revoke the deploykey of the authorization repository.
		pid, ok := authorizationpid.(int)
		if !ok {
			return fmt.Errorf("the ID of the authorized repository is not of type int during applyDeploykey: %v", authorizationpid)
		}
		if pid == repoID {
			continue
		}

		secretData, err := c.GetDeployKeyFromSecretRepo(ctx, fmt.Sprintf("%s%d", RepoPrefix, repoID), DefaultUser, permissions)
		if err != nil {
			if commonv1.IsDeploykeyNotFound(err) {
				return nil
			}
			return err
		}

		deploykey, ok := c.cacheStore.projectDeployKeyMap[repoID][secretData.ID]
		if ok {
			err = fn(ctx, authorizationpid, deploykey.ID)
			if err != nil {
				return err
			}
			continue
		}

		deploykey, err = c.codeRepo.GetDeployKey(ctx, repoID, secretData.ID)
		if err != nil {
			if commonv1.IsDeploykeyNotFound(err) {
				c.log.Debugf("failed to get deploykey during applyDeploykey, repo id: %d, err: %w", repoID, err)
				continue
			}
			return fmt.Errorf("failed to get deploykey during applyDeploykey, repo id: %d, err: %w", repoID, err)
		}

		err = fn(ctx, authorizationpid, deploykey.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CodeRepoBindingUsecase) refreshAuthorization(ctx context.Context, nodes nodestree.Node, codeRepoName string) error {
	err := c.clearInvalidDeployKey(ctx, nodes)
	if err != nil {
		return err
	}

	err = c.authorizeForSameProjectRepo(ctx, nodes)
	if err != nil {
		return err
	}

	if err := c.processAuthorization(ctx, nodes, string(ReadOnly), codeRepoName); err != nil {
		return err
	}

	if err := c.processAuthorization(ctx, nodes, string(ReadWrite), codeRepoName); err != nil {
		return err
	}

	return nil
}

// processAuthorization Calculate the authorization scopes, Perform corresponding operations based on the authorization scopes.
func (c *CodeRepoBindingUsecase) processAuthorization(ctx context.Context, nodes nodestree.Node, permissions, authRepoName string) error {
	pid, err := utilstrings.ExtractNumber(RepoPrefix, authRepoName)
	if err != nil {
		return err
	}

	codeRepoBindings := c.getCodeRepoBindingsInAuthorizedRepo(ctx, nodes, authRepoName, permissions)
	scopes := c.calculateAuthorizationScopes(ctx, codeRepoBindings, permissions)
	for _, scope := range scopes {
		if scope.isProductPermission {
			err := c.updateAllAuthorization(ctx, nodes, pid, permissions, scope.ProductName)
			if err != nil {
				return err
			}
		} else {
			err := c.recycleAuthorization(ctx, scope.ProjectScopes, nodes, pid, permissions, scope.ProductName)
			if err != nil {
				return err
			}

			err = c.updateAuthorization(ctx, scope.ProjectScopes, nodes, pid, permissions, scope.ProductName)
			if err != nil {
				return err
			}
		}
	}

	if len(codeRepoBindings) == 0 {
		err := c.recycleAuthorization(ctx, nil, nodes, pid, permissions, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CodeRepoBindingUsecase) authorizeForSameProjectRepo(ctx context.Context, nodes nodestree.Node) error {
	codeRepos := c.getCodeRepos(nodes)

	errChan := make(chan error)
	semaphore := make(chan struct{}, 10)

	wg := sync.WaitGroup{}
	for i := 0; i < len(codeRepos); i++ {
		for j := i + 1; j < len(codeRepos); j++ {
			repo1 := codeRepos[i]
			repo2 := codeRepos[j]

			if repo1.Name == repo2.Name {
				continue
			}

			if repo1.Spec.Project != repo2.Spec.Project {
				continue
			}

			wg.Add(1)
			go func(repo1, repo2 *resourcev1alpha1.CodeRepo) {
				defer wg.Done()

				semaphore <- struct{}{}
				defer func() {
					<-semaphore
				}()

				c.authorizeRepositories(ctx, repo1, repo2, errChan)
			}(repo1, repo2)
		}
	}

	go func() {
		wg.Wait()
		close(errChan)
		close(semaphore)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (*CodeRepoBindingUsecase) getCodeRepos(nodes nodestree.Node) []*resourcev1alpha1.CodeRepo {
	codeRepoNodes := nodestree.ListsResourceNodes(nodes, nodestree.CodeRepo)
	tmpCodeRepos := make([]*resourcev1alpha1.CodeRepo, 0)

	for _, codeRepoNode := range codeRepoNodes {
		codeRepo, ok := codeRepoNode.Content.(*resourcev1alpha1.CodeRepo)
		if ok {
			tmpCodeRepos = append(tmpCodeRepos, codeRepo)
		}
	}
	return tmpCodeRepos
}

func (c *CodeRepoBindingUsecase) authorizeRepositories(ctx context.Context, repo1, repo2 *resourcev1alpha1.CodeRepo, errChan chan error) {

	c.lock.Lock()

	tmpSecretDeploykeyMap := make(map[string]*DeployKeySecretData, 0)

	roDeployKey1Info, err := c.getDeployKey(ctx, repo1, tmpSecretDeploykeyMap, ReadOnly)
	if err != nil {
		c.lock.Unlock()
		if commonv1.IsDeploykeyNotFound(err) {
			return
		}
		errChan <- err
		return
	}

	rwDeployKey1Info, err := c.getDeployKey(ctx, repo1, tmpSecretDeploykeyMap, ReadWrite)
	if err != nil {
		c.lock.Unlock()
		if commonv1.IsDeploykeyNotFound(err) {
			return
		}
		errChan <- err
		return
	}

	roDeployKey2Info, err := c.getDeployKey(ctx, repo2, tmpSecretDeploykeyMap, ReadOnly)
	if err != nil {
		c.lock.Unlock()
		if commonv1.IsDeploykeyNotFound(err) {
			return
		}
		errChan <- err
		return
	}

	rwDeployKey2Info, err := c.getDeployKey(ctx, repo2, tmpSecretDeploykeyMap, ReadWrite)
	if err != nil {
		c.lock.Unlock()
		if commonv1.IsDeploykeyNotFound(err) {
			return
		}
		errChan <- err
		return
	}

	c.lock.Unlock()

	err = c.enableProjectDeployKey(ctx, repo1, roDeployKey2Info, rwDeployKey2Info)
	if err != nil {
		errChan <- err
		return
	}

	err = c.enableProjectDeployKey(ctx, repo2, roDeployKey1Info, rwDeployKey1Info)
	if err != nil {
		errChan <- err
		return
	}
}

func (c *CodeRepoBindingUsecase) getDeployKey(ctx context.Context, repo *resourcev1alpha1.CodeRepo, tmpSecretDeploykeyMap map[string]*DeployKeySecretData, permission DeployKeyType) (*ProjectDeployKey, error) {
	var err error

	pid, err := utilstrings.ExtractNumber(RepoPrefix, repo.Name)
	if err != nil {
		return nil, err
	}

	// Building cache keys for deploy keys.
	key := fmt.Sprintf("%s-%s", repo.Name, permission)

	// Check if there is a deploy key in the cache, and if it does not exist, obtain it from the keystore.
	deployKey, ok := tmpSecretDeploykeyMap[key]
	if !ok {
		deployKey, err = c.GetDeployKeyFromSecretRepo(ctx, repo.Name, DefaultUser, string(permission))
		if err != nil {
			if commonv1.IsDeploykeyNotFound(err) {
				return nil, nil
			}
			return nil, err
		}
		tmpSecretDeploykeyMap[key] = deployKey
	}

	// Check if there is a deploy key for the project in the cache, and if it does not exist, obtain it from codeRepo.
	projectDeploykey, ok := c.cacheStore.projectDeployKeyMap[pid][deployKey.ID]
	if !ok {
		projectDeploykey, err = c.codeRepo.GetDeployKey(ctx, pid, deployKey.ID)
		if err != nil {
			return nil, err
		}
		c.cacheStore.projectDeployKeyMap[pid] = map[int]*ProjectDeployKey{
			deployKey.ID: projectDeploykey,
		}
	}

	return projectDeploykey, nil
}
func (c *CodeRepoBindingUsecase) enableProjectDeployKey(ctx context.Context, repo *resourcev1alpha1.CodeRepo, roDeployKey2Info *ProjectDeployKey, rwDeployKey2Info *ProjectDeployKey) error {
	pid, err := utilstrings.ExtractNumber(RepoPrefix, repo.Name)
	if err != nil {
		return err
	}

	if _, ok := c.cacheStore.projectDeployKeyMap[pid][roDeployKey2Info.ID]; !ok {
		err = c.enableDeployKey(ctx, pid, roDeployKey2Info.ID)
		if err != nil {
			return err
		}
	}

	if _, ok := c.cacheStore.projectDeployKeyMap[pid][rwDeployKey2Info.ID]; !ok {
		err = c.enableDeployKey(ctx, pid, rwDeployKey2Info.ID)
		if err != nil {
			return err
		}

		err = c.updateDeployKey(ctx, pid, rwDeployKey2Info.ID, rwDeployKey2Info.Title, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CodeRepoBindingUsecase) enableDeployKey(ctx context.Context, pid int, deployKeyID int) error {
	_, err := c.codeRepo.EnableProjectDeployKey(ctx, pid, deployKeyID)
	if err != nil {
		return err
	}
	return nil
}

func (c *CodeRepoBindingUsecase) updateDeployKey(ctx context.Context, pid int, deployKeyID int, title string, enabled bool) error {
	_, err := c.codeRepo.UpdateDeployKey(ctx, pid, deployKeyID, title, enabled)
	if err != nil {
		return err
	}
	return nil
}

func Contains(arr []int, target int) bool {
	for _, element := range arr {
		if element == target {
			return true
		}
	}
	return false
}

func (c *CodeRepoBindingUsecase) getCodeRepoBindingsInAuthorizedRepo(ctx context.Context, nodes nodestree.Node, codeRepoName, permissions string) []*resourcev1alpha1.CodeRepoBinding {
	codeRepoBindingNodes := nodestree.ListsResourceNodes(nodes, nodestree.CodeRepoBinding, func(node *nodestree.Node) bool {
		val, ok := node.Content.(*resourcev1alpha1.CodeRepoBinding)
		if !ok {
			return false
		}

		if val.Spec.CodeRepo == codeRepoName && val.Spec.Permissions == permissions {
			return true
		}

		return false
	})

	codeRepoBindings := make([]*resourcev1alpha1.CodeRepoBinding, 0)
	for _, node := range codeRepoBindingNodes {
		val, ok := node.Content.(*resourcev1alpha1.CodeRepoBinding)
		if ok {
			codeRepoBindings = append(codeRepoBindings, val)
		}
	}

	return codeRepoBindings
}

func (c *CodeRepoBindingUsecase) getAuthorizedRepoCodeRepo(ctx context.Context, nodes nodestree.Node, authRepoName string) (*resourcev1alpha1.CodeRepo, error) {
	node := c.nodestree.GetNode(&nodes, nodestree.CodeRepo, authRepoName)
	codeRepo, ok := node.Content.(*resourcev1alpha1.CodeRepo)
	if !ok {
		return nil, fmt.Errorf("wrong type found for %s node when checking CodeRepo type", node.Name)
	}
	return codeRepo, nil
}

type ProductAuthorization struct {
	ProductName         string
	ProjectScopes       map[string]bool
	isProductPermission bool
}

// calculateAuthorizationScopes calculate the permissions for each product based on CodeRepobindings.
// Return a list of ProductAuthorization entities, recording the authorization scope for each product.
// Each entity contains permissions for product name, product level, and project scopes.
// Projectscopes have merged all project permissions under CodeRepoBinding for the product, This is a reserved data.
func (c *CodeRepoBindingUsecase) calculateAuthorizationScopes(ctx context.Context, codeRepoBindings []*resourcev1alpha1.CodeRepoBinding, permissions string) []*ProductAuthorization {
	scopes := []*ProductAuthorization{}

	for _, codeRepoBinding := range codeRepoBindings {
		isProductPermission := len(codeRepoBinding.Spec.Projects) == 0
		projectScopes := make(map[string]bool)

		for _, project := range codeRepoBinding.Spec.Projects {
			projectScopes[project] = true
		}

		existingScope := findExistingScope(scopes, codeRepoBinding.Spec.Product)
		if existingScope != nil {
			if isProductPermission {
				existingScope.isProductPermission = isProductPermission
			}
			for key, _ := range projectScopes {
				existingScope.ProjectScopes[key] = true
			}
		} else {
			// authorizedCodeRepoProject is project of the authorized repository, when isProductPermission is false that it is required.
			scopes = append(scopes, &ProductAuthorization{
				ProductName:         codeRepoBinding.Spec.Product,
				isProductPermission: isProductPermission,
				ProjectScopes:       projectScopes,
			})
		}
	}

	return scopes
}

func findExistingScope(scopes []*ProductAuthorization, productName string) *ProductAuthorization {
	for _, scope := range scopes {
		if scope.ProductName == productName {
			return scope
		}
	}
	return nil
}

func (c *CodeRepoBindingUsecase) updateAllAuthorization(ctx context.Context, nodes nodestree.Node, pid interface{}, permissions, product string) error {
	codeRepos, err := nodesToCodeRepoists(nodes)
	if err != nil {
		return err
	}

	err = c.authorizeDeployKey(ctx, codeRepos, pid, permissions)
	if err != nil {
		return err
	}

	return nil
}

// updateAuthorization Authorize according to the authorization scope of the project
func (c *CodeRepoBindingUsecase) updateAuthorization(ctx context.Context, projectScopes map[string]bool, nodes nodestree.Node, pid interface{}, permissions, product string) error {
	if len(projectScopes) == 0 {
		return nil
	}

	codeRepos, err := nodesToCodeRepoists(nodes, func(codeRepo *resourcev1alpha1.CodeRepo) bool {
		_, ok := projectScopes[codeRepo.Spec.Project]
		return ok
	})

	err = c.authorizeDeployKey(ctx, codeRepos, pid, permissions)
	if err != nil {
		return err
	}

	return nil
}

// recycleAuthorization Recycle according to the authorization scope of the project.
func (c *CodeRepoBindingUsecase) recycleAuthorization(ctx context.Context, projectScopes map[string]bool, nodes nodestree.Node, pid interface{}, permissions, productName string) error {
	if len(projectScopes) == 0 {
		return nil
	}

	repository, err := c.codeRepo.GetCodeRepo(ctx, pid)
	if err != nil {
		return err
	}

	codeRepos := []*resourcev1alpha1.CodeRepo{}
	currentProductName := fmt.Sprintf("%s%d", _ProductPrefix, repository.Namespace.ID)
	if productName == currentProductName {
		codeRepos, err = nodesToCodeRepoists(nodes, func(codeRepo *resourcev1alpha1.CodeRepo) bool {
			if codeRepo.Spec.Project == "" {
				return false
			}

			_, ok := projectScopes[codeRepo.Spec.Project]
			return !ok
		})
		if err != nil {
			return err
		}
	} else {
		// TODO: Increase cross-product processing
	}

	if err := c.RevokeDeployKey(ctx, codeRepos, pid, permissions); err != nil {
		return err
	}
	return nil
}

type deployKeyMapValue struct {
	deployKey *ProjectDeployKey
	codeRepos []*resourcev1alpha1.CodeRepo
}

func (c *CodeRepoBindingUsecase) clearInvalidDeployKey(ctx context.Context, nodes nodestree.Node) error {
	codeRepoNodes := nodestree.ListsResourceNodes(nodes, nodestree.CodeRepo)
	cacheProjectMap := make(map[int]*Project)
	cacheProjectsDeploykeyMap := make(map[int]map[int]*ProjectDeployKey, len(codeRepoNodes))
	deploykeyInAllProjectsMap := make(map[string]*deployKeyMapValue, len(codeRepoNodes))

	for _, codeRepoNode := range codeRepoNodes {
		codeRepo, ok := codeRepoNode.Content.(*resourcev1alpha1.CodeRepo)
		if !ok {
			continue
		}

		c.wg.Add(1)

		go func(codeRepo *resourcev1alpha1.CodeRepo) {
			defer c.wg.Done()
			err := c.deduplicateAndCacheDeployKeys(ctx, codeRepo, cacheProjectsDeploykeyMap, deploykeyInAllProjectsMap)
			if err != nil {
				return
			}
		}(codeRepo)

	}

	c.wg.Wait()

	for key, val := range deploykeyInAllProjectsMap {
		// Parsing the name of deploykey to obtain the ID of the code repository.
		// eg: repo-22-readwrite, repository id is 22.
		re := regexp.MustCompile(`repo-(\d+)-`)
		match := re.FindStringSubmatch(key)
		if len(match) == 0 {
			continue
		}
		pid, err := strconv.Atoi(match[1])
		if err != nil {
			return err
		}

		deployKeyID := val.deployKey.ID

		c.wg.Add(1)

		go func(pid, deployKeyID int, codeRepos []*resourcev1alpha1.CodeRepo) {
			var isDeleteDeployKey bool
			var DeleteDeployKeys = make(map[int]bool)

			defer c.wg.Done()
			defer func(pid int) {
				if !isDeleteDeployKey {
					return
				}

				// Delete the deploykey associated with the repository id.
				err := c.deleteAssociatedRepositoryDeployKey(ctx, codeRepos, deployKeyID)
				if err != nil {
					return
				}
			}(pid)

			isDeleteDeployKey, err = c.checkRepositoryExistence(ctx, cacheProjectMap, pid)
			if err != nil {
				return
			}

			isDeleteDeployKey, err = c.checkDeployKeyExistence(ctx, pid, deployKeyID, DeleteDeployKeys)
			if err != nil {
				return
			}

			err := c.checkProjectConsistency(ctx, nodes, codeRepos, pid, deployKeyID)
			if err != nil {
				return
			}

		}(pid, deployKeyID, val.codeRepos)
	}

	c.wg.Wait()

	c.cacheStore.projectDeployKeyMap = cacheProjectsDeploykeyMap

	return nil
}

func (c *CodeRepoBindingUsecase) deleteAssociatedRepositoryDeployKey(ctx context.Context, codeRepos []*resourcev1alpha1.CodeRepo, deployKeyID int) error {
	for _, codeRepo := range codeRepos {
		pid, err := utilstrings.ExtractNumber(RepoPrefix, codeRepo.Name)
		if err != nil {
			return err
		}
		if err := c.codeRepo.DeleteDeployKey(ctx, pid, deployKeyID); err != nil {
			return err
		}
	}
	return nil
}

func (c *CodeRepoBindingUsecase) deduplicateAndCacheDeployKeys(ctx context.Context, codeRepo *resourcev1alpha1.CodeRepo, cacheProjectsDeploykeyMap map[int]map[int]*ProjectDeployKey, deploykeyInAllProjectsMap map[string]*deployKeyMapValue) error {
	pid, err := utilstrings.ExtractNumber(RepoPrefix, codeRepo.Name)
	if err != nil {
		return err
	}

	projectDeployKeys, err := GetAllDeployKeys(ctx, c.codeRepo, pid)
	if err != nil {
		return err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	for _, projectDeployKey := range projectDeployKeys {
		if cacheProjectsDeploykeyMap[pid] == nil {
			cacheProjectsDeploykeyMap[pid] = make(map[int]*ProjectDeployKey)
			cacheProjectsDeploykeyMap[pid][projectDeployKey.ID] = projectDeployKey
		} else {
			if _, ok := cacheProjectsDeploykeyMap[pid][projectDeployKey.ID]; !ok {
				cacheProjectsDeploykeyMap[pid][projectDeployKey.ID] = projectDeployKey
			}
		}

		if _, ok := deploykeyInAllProjectsMap[projectDeployKey.Title]; !ok {
			deploykeyInAllProjectsMap[projectDeployKey.Title] = &deployKeyMapValue{
				deployKey: projectDeployKey,
				codeRepos: []*resourcev1alpha1.CodeRepo{codeRepo},
			}
		} else {
			deploykeyInAllProjectsMap[projectDeployKey.Title].codeRepos = append(deploykeyInAllProjectsMap[projectDeployKey.Title].codeRepos, codeRepo)
		}
	}

	return nil
}

func (c *CodeRepoBindingUsecase) checkRepositoryExistence(ctx context.Context, cacheProjectMap map[int]*Project, pid int) (bool, error) {
	var repository *Project
	var ok bool
	var err error

	c.lock.Lock()

	repository, ok = cacheProjectMap[pid]
	if !ok {
		sg := &singleflight.Group{}
		sg.Do(fmt.Sprintf("%d", pid), func() (interface{}, error) {
			repository, err = c.codeRepo.GetCodeRepo(ctx, pid)
			if err != nil {
				if !commonv1.IsProjectNotFound(err) {
					c.lock.Unlock()
					return false, err
				}

				return true, nil
			}

			cacheProjectMap[pid] = repository

			return false, nil
		})
	}
	c.lock.Unlock()

	return false, nil
}

func (c *CodeRepoBindingUsecase) checkDeployKeyExistence(ctx context.Context, pid int, deployKeyID int, DeleteDeployKeys map[int]bool) (bool, error) {
	_, err := c.codeRepo.GetDeployKey(ctx, pid, deployKeyID)
	if err != nil {
		if !commonv1.IsDeploykeyNotFound(err) {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func (c *CodeRepoBindingUsecase) checkProjectConsistency(ctx context.Context, nodes nodestree.Node, codeRepos []*resourcev1alpha1.CodeRepo, authorizationpid int, deployKeyID int) error {
	authorizationCodeRepoNode := c.nodestree.GetNode(&nodes, nodestree.CodeRepo, fmt.Sprintf("%s%d", RepoPrefix, authorizationpid))
	authorizationCodeRepo, ok := authorizationCodeRepoNode.Content.(*resourcev1alpha1.CodeRepo)
	if !ok {
		return fmt.Errorf("failed to check project consistency, illegal codeRepo: %s", authorizationCodeRepoNode.Name)
	}

	for _, codeRepo := range codeRepos {
		if codeRepo.Spec.Project != authorizationCodeRepo.Spec.Project {
			pid, err := utilstrings.ExtractNumber(RepoPrefix, codeRepo.Name)
			if err != nil {
				return err
			}

			if err := c.codeRepo.DeleteDeployKey(ctx, pid, deployKeyID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *CodeRepoBindingUsecase) RevokeDeployKey(ctx context.Context, codeRepos []*resourcev1alpha1.CodeRepo, authpid interface{}, permissions string) error {
	for _, codeRepo := range codeRepos {
		pid, err := utilstrings.ExtractNumber(RepoPrefix, codeRepo.Name)
		if err != nil {
			return err
		}

		if authpid == pid {
			continue
		}

		secretData, err := c.GetDeployKeyFromSecretRepo(ctx, codeRepo.Name, DefaultUser, permissions)
		if err != nil {
			if commonv1.IsDeploykeyNotFound(err) {
				return nil
			}
			return err
		}

		err = c.codeRepo.DeleteDeployKey(ctx, authpid, secretData.ID)
		if err != nil && !commonv1.IsDeploykeyNotFound(err) {
			return err
		}
	}

	return nil
}

func (c *CodeRepoBindingUsecase) GetDeployKeyFromSecretRepo(ctx context.Context, repoName, user, permissions string) (*DeployKeySecretData, error) {
	gitType := c.config.Git.GitType
	secretsEngine := SecretsGitEngine
	secretsKey := SecretsDeployKey
	secretPath := fmt.Sprintf("%s/%s/%s/%s", gitType, repoName, user, permissions)
	secretOptions := &SecretOptions{
		SecretPath:   secretPath,
		SecretEngine: secretsEngine,
		SecretKey:    secretsKey,
	}

	deployKeySecretData, err := c.secretRepo.GetDeployKey(ctx, secretOptions)
	if err != nil {
		return nil, err
	}

	return deployKeySecretData, nil
}

func (c *CodeRepoBindingUsecase) CheckReference(options nodestree.CompareOptions, node *nodestree.Node, k8sClient client.Client) (bool, error) {
	if node.Kind != nodestree.CodeRepoBinding {
		return false, nil
	}

	codeRepoBinding, ok := node.Content.(*resourcev1alpha1.CodeRepoBinding)
	if !ok {
		return true, fmt.Errorf("wrong type found for %s node when checking CodeRepoBinding type", node.Name)
	}

	ok = nodestree.IsResourceExist(options, codeRepoBinding.Spec.CodeRepo, nodestree.CodeRepo)
	if !ok {
		objKey := client.ObjectKey{
			Namespace: codeRepoBinding.Spec.Product,
			Name:      codeRepoBinding.Spec.CodeRepo,
		}

		err := k8sClient.Get(context.TODO(), objKey, &resourcev1alpha1.CodeRepo{})
		if err != nil {
			return true, err
		}
	}

	// TODO:
	// Query through k8s cluster when crossing products.
	for _, project := range codeRepoBinding.Spec.Projects {
		ok = nodestree.IsResourceExist(options, project, nodestree.Project)
		if !ok {
			return true, fmt.Errorf("project resource %s does not exist or is invalid", project)
		}
	}

	return true, nil
}

func (c *CodeRepoBindingUsecase) CreateNode(path string, data interface{}) (*nodestree.Node, error) {
	val, ok := data.(*CodeRepoBindingData)
	if !ok {
		return nil, fmt.Errorf("failed to creating specify node, the path: %s", path)
	}

	if len(val.Spec.Projects) == 0 {
		val.Spec.Projects = make([]string, 0)
	}

	codeRepo := &resourcev1alpha1.CodeRepoBinding{
		TypeMeta: v1.TypeMeta{
			Kind:       nodestree.CodeRepoBinding,
			APIVersion: resourcev1alpha1.GroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name: val.Name,
		},
		Spec: val.Spec,
	}

	resourceDirectory := fmt.Sprintf("%s/%s", path, "code-repos")
	resourcePath := fmt.Sprintf("%s/%s/%s.yaml", resourceDirectory, val.Spec.CodeRepo, val.Name)

	return &nodestree.Node{
		Name:    val.Name,
		Path:    resourcePath,
		Content: codeRepo,
		Kind:    nodestree.CodeRepoBinding,
		Level:   4,
	}, nil
}

func (c *CodeRepoBindingUsecase) UpdateNode(resourceNode *nodestree.Node, data interface{}) (*nodestree.Node, error) {
	val, ok := data.(*CodeRepoBindingData)
	if !ok {
		return nil, fmt.Errorf("failed to get conderepo %s data when updating node", resourceNode.Name)
	}

	if len(val.Spec.Projects) == 0 {
		val.Spec.Projects = make([]string, 0)
	}

	codeRepoBinding, ok := resourceNode.Content.(*resourcev1alpha1.CodeRepoBinding)
	if !ok {
		return nil, fmt.Errorf("failed to get coderepo %s when updating node", resourceNode.Name)
	}

	if val.Spec.CodeRepo != codeRepoBinding.Spec.CodeRepo {
		return nil, errors.New(500, "NOT_ALLOWED_MODIFY", "It is not allowed to modify the authorized repository. If you want to change the authorized repository, please delete the authorization")
	}

	codeRepoBinding.Spec = val.Spec
	resourceNode.Content = codeRepoBinding

	return resourceNode, nil
}

func (c *CodeRepoBindingUsecase) CreateResource(kind string) interface{} {
	if kind != nodestree.CodeRepoBinding {
		return nil
	}

	return &resourcev1alpha1.CodeRepoBinding{}
}

func (c *CodeRepoBindingUsecase) nodeToResource(node *nodestree.Node) (*resourcev1alpha1.CodeRepoBinding, error) {
	r, ok := node.Content.(*resourcev1alpha1.CodeRepoBinding)
	if !ok {
		return nil, fmt.Errorf("failed to get instance when get %s coderepoBinding", node.Name)
	}

	return r, nil
}
