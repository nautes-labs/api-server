package biz

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	commonv1 "github.com/nautes-labs/api-server/api/common/v1"
	"github.com/nautes-labs/api-server/pkg/nodestree"
	utilstrings "github.com/nautes-labs/api-server/util/string"
	resourcev1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
	nautesconfigs "github.com/nautes-labs/pkg/pkg/nautesconfigs"
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
	repositoryName   string
}

type CodeRepoBindingData struct {
	Name string
	Spec resourcev1alpha1.CodeRepoBindingSpec
}

type applyDeploykeyFunc func(ctx context.Context, pid interface{}, deployKey int) error

func NewCodeRepoCodeRepoBindingUsecase(logger log.Logger, codeRepo CodeRepo, secretRepo Secretrepo, nodestree nodestree.NodesTree, resourcesUsecase *ResourcesUsecase, config *nautesconfigs.Config, client client.Client) *CodeRepoBindingUsecase {
	codeRepoBindingUsecase := &CodeRepoBindingUsecase{
		log:              log.NewHelper(log.With(logger)),
		codeRepo:         codeRepo,
		secretRepo:       secretRepo,
		nodestree:        nodestree,
		resourcesUsecase: resourcesUsecase,
		config:           config,
		client:           client,
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

	repoName, err := c.resourcesUsecase.convertCodeRepoToRepoName(ctx, resource.Spec.CodeRepo)
	if err != nil {
		return nil, err
	}
	resource.Spec.CodeRepo = repoName

	groupName, err := c.resourcesUsecase.convertProductToGroupName(ctx, resource.Spec.Product)
	if err != nil {
		return nil, err
	}
	resource.Spec.Product = groupName

	return resource, nil
}

func (c *CodeRepoBindingUsecase) ListCodeRepoBindings(ctx context.Context, options *BizOptions) ([]*resourcev1alpha1.CodeRepoBinding, error) {
	nodes, err := c.resourcesUsecase.List(ctx, options.ProductName, c)
	if err != nil {
		return nil, err
	}

	codeRepoBindings, err := c.nodesToLists(*nodes)
	if err != nil {
		return nil, err
	}

	for _, codecodeRepoBinding := range codeRepoBindings {
		repoName, err := c.resourcesUsecase.convertCodeRepoToRepoName(ctx, codecodeRepoBinding.Spec.CodeRepo)
		if err != nil {
			return nil, err
		}
		codecodeRepoBinding.Spec.CodeRepo = repoName

		groupName, err := c.resourcesUsecase.convertProductToGroupName(ctx, codecodeRepoBinding.Spec.Product)
		if err != nil {
			return nil, err
		}
		codecodeRepoBinding.Spec.Product = groupName
	}

	return codeRepoBindings, nil
}

func (c *CodeRepoBindingUsecase) SaveCodeRepoBinding(ctx context.Context, options *BizOptions, data *CodeRepoBindingData) error {
	productName, err := c.resourcesUsecase.convertGroupToProduct(ctx, data.Spec.Product)
	if err != nil {
		return err
	}
	c.groupName = data.Spec.Product
	data.Spec.Product = productName

	codeRepoName, err := c.resourcesUsecase.convertRepoNameToCodeRepo(ctx, options.ProductName, data.Spec.CodeRepo)
	if err != nil {
		return err
	}
	c.repositoryName = data.Spec.CodeRepo
	data.Spec.CodeRepo = codeRepoName

	nodes, err := c.resourcesUsecase.loadDefaultProjectNodes(ctx, options.ProductName)
	if err != nil {
		return err
	}

	node := c.nodestree.GetNode(nodes, nodestree.CodeRepoBinding, options.ResouceName)
	if node != nil {
		lastCodeRepoBinding, ok := node.Content.(*resourcev1alpha1.CodeRepoBinding)
		if !ok {
			return fmt.Errorf("resource type is inconsistent, please check if this resource %s is legal", options.ResouceName)
		}

		if lastCodeRepoBinding.Spec.CodeRepo != data.Spec.CodeRepo {
			return fmt.Errorf("It is not allowed to modify the authorized repository. If you want to change the authorized repository, please delete the authorization")
		}
	}

	resourceOptions := &resourceOptions{
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

func (c *CodeRepoBindingUsecase) authorizeDeployKey(ctx context.Context, codeRepos []*resourcev1alpha1.CodeRepo, pid interface{}, permisions string) error {
	var repositories []*Project
	for _, repo := range codeRepos {
		pid, err := utilstrings.ExtractNumber(RepoPrefix, repo.Name)
		if err != nil {
			return err
		}
		project, err := c.codeRepo.GetCodeRepo(ctx, pid)
		if err != nil {
			return err
		}
		repositories = append(repositories, project)
	}

	if err := c.applyDeploykey(ctx, pid, permisions, repositories, func(ctx context.Context, pid interface{}, deployKey int) error {
		_, err := c.codeRepo.GetDeployKey(ctx, pid, deployKey)
		if err != nil {
			if commonv1.IsDeploykeyNotFound(err) {
				_, err := c.codeRepo.EnableProjectDeployKey(ctx, pid, deployKey)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (c *CodeRepoBindingUsecase) applyDeploykey(ctx context.Context, pid interface{}, permissions string, repositories []*Project, fn applyDeploykeyFunc) error {
	for _, repository := range repositories {
		repoid := fmt.Sprintf("%s%d", RepoPrefix, repository.Id)
		secretData, err := c.GetDeployKeyFromSecretRepo(ctx, repoid, DefaultUser, permissions)
		if err != nil {
			return commonv1.ErrorDeploykeyNotFound("failed to get the %s deploykey from secret repo, please check if the key under /%s/%s exists or is invalid", permissions, c.config.Git.GitType, repoid)
		}

		//Revoke in-product authorization, not revoke the deploykey of the authorization repository.
		if pid == int(repository.Id) {
			continue
		}

		deploykey, err := c.codeRepo.GetDeployKey(ctx, int(repository.Id), secretData.ID)
		if err != nil {
			return commonv1.ErrorDeploykeyNotFound("failed to get the %s deploykey from git, please check if the key exists or is invalid for the repository %s under organization %s", permissions, repository.Name, c.groupName)
		}

		err = fn(ctx, pid, deploykey.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CodeRepoBindingUsecase) refreshAuthorization(ctx context.Context, nodes nodestree.Node, codeRepoName string) error {
	if err := c.clearInvalidDeployKey(ctx, codeRepoName); err != nil {
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
	codeRepo, err := c.getAuthorizedRepoCodeRepo(ctx, nodes, authRepoName)
	if err != nil {
		return err
	}

	scopes := c.calculateAuthorizationScopes(ctx, codeRepoBindings, permissions, codeRepo.Spec.Project)
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

func (c *CodeRepoBindingUsecase) getCodeRepoBindingsInAuthorizedRepo(ctx context.Context, nodes nodestree.Node, codeRepoName, permisions string) []*resourcev1alpha1.CodeRepoBinding {
	codeRepoBindingNodes := nodestree.ListsResourceNodes(nodes, nodestree.CodeRepoBinding, func(node *nodestree.Node) bool {
		val, ok := node.Content.(*resourcev1alpha1.CodeRepoBinding)
		if !ok {
			return false
		}

		if val.Spec.CodeRepo == codeRepoName && val.Spec.Permissions == permisions {
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
func (c *CodeRepoBindingUsecase) calculateAuthorizationScopes(ctx context.Context, codeRepoBindings []*resourcev1alpha1.CodeRepoBinding, permissions, authorizedCodeRepoProject string) []*ProductAuthorization {
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
			if !isProductPermission && authorizedCodeRepoProject != "" && !projectScopes[authorizedCodeRepoProject] {
				projectScopes[authorizedCodeRepoProject] = true
			}

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
func (c *CodeRepoBindingUsecase) updateAuthorization(ctx context.Context, scopes map[string]bool, nodes nodestree.Node, pid interface{}, permissions, product string) error {
	if len(scopes) == 0 {
		return nil
	}

	codeRepos, err := nodesToCodeRepoists(nodes, func(codeRepo *resourcev1alpha1.CodeRepo) bool {
		for key, _ := range scopes {
			if codeRepo.Spec.Project == key {
				return true
			}
		}

		return false
	})
	if err != nil {
		return err
	}

	err = c.authorizeDeployKey(ctx, codeRepos, pid, permissions)
	if err != nil {
		return err
	}

	return nil
}

// recycleAuthorization Recycle according to the authorization scope of the project.
func (c *CodeRepoBindingUsecase) recycleAuthorization(ctx context.Context, projectScopes map[string]bool, nodes nodestree.Node, pid interface{}, permissions, product string) error {
	var codeRepos []*resourcev1alpha1.CodeRepo

	repository, err := c.codeRepo.GetCodeRepo(ctx, pid)
	if err != nil {
		return err
	}

	currentProductName := fmt.Sprintf("%s%d", _ProductPrefix, repository.Namespace.ID)
	if product != "" && product != currentProductName {
		// TODO: Increase cross-product processing
	} else {
		tmpProjects := make([]string, 0)

		projectNodes := nodestree.ListsResourceNodes(nodes, nodestree.Project)
		for _, node := range projectNodes {
			val, ok := node.Content.(*resourcev1alpha1.Project)
			if !ok {
				return fmt.Errorf("resource type is inconsistent, please check if this resource %s is legal", node.Name)
			}
			if _, ok := projectScopes[val.Name]; !ok {
				tmpProjects = append(tmpProjects, val.Name)
			}
		}

		codeRepos, err = nodesToCodeRepoists(nodes, func(codeRepo *resourcev1alpha1.CodeRepo) bool {
			for _, project := range tmpProjects {
				if codeRepo.Spec.Project == project {
					return true
				}
			}

			return false
		})
		if err != nil {
			return err
		}
	}

	if err := c.RevokeDeployKey(ctx, codeRepos, pid, permissions); err != nil {
		return err
	}
	return nil
}

func (c *CodeRepoBindingUsecase) clearInvalidDeployKey(ctx context.Context, authorizedRepository string) error {
	pid, err := utilstrings.ExtractNumber(RepoPrefix, authorizedRepository)
	if err != nil {
		return err
	}

	projectDeployKeys, err := GetAllDeployKeys(ctx, c.codeRepo, pid)
	if err != nil {
		return err
	}
	for _, projectDeployKey := range projectDeployKeys {
		re := regexp.MustCompile(`repo-(\d+)-`)
		match := re.FindStringSubmatch(projectDeployKey.Title)
		if len(match) == 0 {
			continue
		}

		matchpid, err := strconv.Atoi(match[1])
		if err != nil {
			return err
		}

		if matchpid == pid {
			continue
		}

		repository, err := c.codeRepo.GetCodeRepo(ctx, matchpid)
		if err != nil {
			if !commonv1.IsProjectNotFound(err) {
				return err
			}
			if err := c.codeRepo.DeleteDeployKey(ctx, pid, projectDeployKey.ID); err != nil {
				return err
			}
		}

		if repository != nil {
			_, err = c.codeRepo.GetDeployKey(ctx, int(repository.Id), projectDeployKey.ID)
			if err != nil {
				if commonv1.IsDeploykeyNotFound(err) {
					err := c.codeRepo.DeleteDeployKey(ctx, pid, projectDeployKey.ID)
					if err != nil {
						return err
					}
				} else {
					return err
				}
			}
		}
	}

	return nil
}

func (c *CodeRepoBindingUsecase) ClearOthersDeployKey(ctx context.Context, pid int) error {
	projectDeployKeys, err := GetAllDeployKeys(ctx, c.codeRepo, pid)
	if err != nil {
		return err
	}
	for _, projectDeployKey := range projectDeployKeys {
		re := regexp.MustCompile(`repo-(\d+)-`)
		match := re.FindStringSubmatch(projectDeployKey.Title)
		if len(match) == 0 {
			continue
		}

		matchpid, err := strconv.Atoi(match[1])
		if err != nil {
			return err
		}

		if matchpid == pid {
			continue
		}

		repository, err := c.codeRepo.GetCodeRepo(ctx, pid)
		if err != nil && commonv1.IsProjectNotFound(err) {
			return nil
		}
		if err != nil {
			return err
		}

		err = c.codeRepo.DeleteDeployKey(ctx, int(repository.Id), projectDeployKey.ID)
		if err != nil && !commonv1.IsDeploykeyNotFound(err) {
			return err
		}
	}

	return nil
}

func (c *CodeRepoBindingUsecase) RevokeDeployKey(ctx context.Context, concodeRepos []*resourcev1alpha1.CodeRepo, authpid interface{}, permissions string) error {
	for _, repo := range concodeRepos {
		pid, err := utilstrings.ExtractNumber(RepoPrefix, repo.Name)
		if err != nil {
			return err
		}

		repository, err := c.codeRepo.GetCodeRepo(ctx, pid)
		if err != nil {
			return err
		}

		if authpid == int(repository.Id) {
			continue
		}

		secretData, err := c.GetDeployKeyFromSecretRepo(ctx, repo.Name, DefaultUser, permissions)
		if err != nil {
			if commonv1.IsDeploykeyNotFound(err) {
				return commonv1.ErrorDeploykeyNotFound("failed to get the %s deploykey from secret repo, please check if the key under /%s/%s exists or is invalid", permissions, c.config.Git.GitType, repo.Name)
			}
			return err
		}

		deploykey, err := c.codeRepo.GetDeployKey(ctx, pid, secretData.ID)
		if err != nil {
			if commonv1.IsDeploykeyNotFound(err) {
				return commonv1.ErrorDeploykeyNotFound("failed to get the %s deploykey from git, please check if the key exists or is invalid for the repository %s under organization %s", permissions, repository.Name, c.groupName)
			}
			return err
		}

		err = c.codeRepo.DeleteDeployKey(ctx, authpid, deploykey.ID)
		if err != nil && !commonv1.IsDeploykeyNotFound(err) {
			return err
		}

	}

	return nil
}

func (c *CodeRepoBindingUsecase) GetAuthorizedRepositories(ctx context.Context, spec resourcev1alpha1.CodeRepoBindingSpec) ([]*Project, error) {
	gid, err := utilstrings.ExtractNumber(_ProductPrefix, spec.Product)
	if err != nil {
		return nil, err
	}
	nodes, err := c.resourcesUsecase.List(ctx, gid, c)
	if err != nil {
		return nil, err
	}

	// Get all CodeRepoBindings for the same authorized repository.
	codeRepoBindingNodes := nodestree.ListsResourceNodes(*nodes, nodestree.CodeRepoBinding, func(node *nodestree.Node) bool {
		val, ok := node.Content.(*resourcev1alpha1.CodeRepoBinding)
		if !ok {
			return false
		}

		if val.Spec.CodeRepo == spec.CodeRepo {
			return true
		}

		return false
	})

	// if product-level authorization is found, all repository are directly authorized, otherwise it will try to merge project authorization scope.
	tmpProjects := make([]string, 0)
	for _, node := range codeRepoBindingNodes {
		val, ok := node.Content.(*resourcev1alpha1.CodeRepoBinding)
		if !ok {
			return nil, fmt.Errorf("resource type is inconsistent, please check if this resource %s is legal", node.Name)
		}
		if len(val.Spec.Projects) == 0 {
			tmpProjects = make([]string, 0)
			break
		}
		tmpProjects = append(tmpProjects, val.Spec.Projects...)
	}

	// Get CodeRepos within the specified authorization scope.
	codeRepos, err := nodesToCodeRepoists(*nodes, func(codeRepo *resourcev1alpha1.CodeRepo) bool {
		for _, project := range tmpProjects {
			if codeRepo.Spec.Project == project {
				return true
			}
		}

		return false
	})
	if err != nil {
		return nil, err
	}

	var repositories []*Project
	for _, repo := range codeRepos {
		pid, err := utilstrings.ExtractNumber(RepoPrefix, repo.Name)
		if err != nil {
			return nil, err
		}
		project, err := c.codeRepo.GetCodeRepo(ctx, pid)
		if err != nil {
			return nil, err
		}
		repositories = append(repositories, project)
	}

	return repositories, nil
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

	codeRepo, ok := resourceNode.Content.(*resourcev1alpha1.CodeRepoBinding)
	if !ok {
		return nil, fmt.Errorf("failed to get coderepo %s when updating node", resourceNode.Name)
	}

	codeRepo.Spec = val.Spec
	resourceNode.Content = codeRepo

	return resourceNode, nil
}

func (c *CodeRepoBindingUsecase) CreateResource(kind string) interface{} {
	if kind != nodestree.CodeRepoBinding {
		return nil
	}

	return &resourcev1alpha1.CodeRepoBinding{}
}

func (c *CodeRepoBindingUsecase) nodesToLists(nodes nodestree.Node) ([]*resourcev1alpha1.CodeRepoBinding, error) {
	var codeReposDir *nodestree.Node
	var resources []*resourcev1alpha1.CodeRepoBinding

	for _, child := range nodes.Children {
		if child.Name == _CodeReposSubDir {
			codeReposDir = child
			break
		}
	}

	if codeReposDir == nil {
		return nil, fmt.Errorf("the %s directory is not exist", _CodeReposSubDir)
	}

	for _, subNode := range codeReposDir.Children {
		for _, node := range subNode.Children {
			if node.Kind != nodestree.CodeRepoBinding {
				continue
			}

			resource, err := c.nodeToResource(node)
			if err != nil {
				return nil, err
			}
			resources = append(resources, resource)
		}
	}

	return resources, nil
}

func (c *CodeRepoBindingUsecase) nodeToResource(node *nodestree.Node) (*resourcev1alpha1.CodeRepoBinding, error) {
	r, ok := node.Content.(*resourcev1alpha1.CodeRepoBinding)
	if !ok {
		return nil, fmt.Errorf("failed to get instance when get %s coderepoBinding", node.Name)
	}

	return r, nil
}
