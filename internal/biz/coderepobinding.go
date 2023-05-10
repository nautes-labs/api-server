package biz

import (
	"context"
	"fmt"
	"reflect"

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
	codeRepoUsecase  *CodeRepoUsecase
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

type ApplyDeploykeyFunc func(ctx context.Context, pid interface{}, deployKey int) error

func NewCodeRepoCodeRepoBindingUsecase(logger log.Logger, codeRepo CodeRepo, secretRepo Secretrepo, nodestree nodestree.NodesTree, codeRepoUsecase *CodeRepoUsecase, resourcesUsecase *ResourcesUsecase, config *nautesconfigs.Config, client client.Client) *CodeRepoBindingUsecase {
	codeRepoBindingUsecase := &CodeRepoBindingUsecase{log: log.NewHelper(log.With(logger)), codeRepo: codeRepo, secretRepo: secretRepo, nodestree: nodestree, codeRepoUsecase: codeRepoUsecase, resourcesUsecase: resourcesUsecase, config: config, client: client}
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

	resourceOptions := &resourceOptions{
		resourceKind:      nodestree.CodeRepoBinding,
		productName:       options.ProductName,
		insecureSkipCheck: options.InsecureSkipCheck,
		operator:          c,
	}

	lastSpec, err := c.GetCodeRepoBindingSpec(ctx, options.ProductName, options.ResouceName)
	if err != nil {
		return err
	}

	if lastSpec != nil {
		equal := reflect.DeepEqual(*lastSpec, data.Spec)
		if equal {
			return nil
		}

		if err := c.ClearLastDeployKeyResult(ctx, lastSpec); err != nil {
			return err
		}
	}

	if err := c.resourcesUsecase.Save(ctx, resourceOptions, data); err != nil {
		return err
	}

	if err := c.AuthorizeDeployKey(ctx, data.Spec); err != nil {
		return err
	}

	return nil
}

func (c *CodeRepoBindingUsecase) DeleteCodeRepoBinding(ctx context.Context, options *BizOptions) error {
	resourceOptions := &resourceOptions{
		resourceName:      options.ResouceName,
		resourceKind:      nodestree.CodeRepoBinding,
		productName:       options.ProductName,
		insecureSkipCheck: options.InsecureSkipCheck,
		operator:          c,
	}

	lastSpec, err := c.GetCodeRepoBindingSpec(ctx, options.ProductName, options.ResouceName)
	if err != nil {
		return err
	}

	if err := c.ClearLastDeployKeyResult(ctx, lastSpec); err != nil {
		return err
	}

	if err := c.resourcesUsecase.Delete(ctx, resourceOptions, func(nodes nodestree.Node) (string, error) {
		return resourceOptions.resourceName, nil
	}); err != nil {
		return err
	}

	return nil
}

func (c *CodeRepoBindingUsecase) GetCodeRepoBindingSpec(ctx context.Context, productName, coderepoBindingName string) (*resourcev1alpha1.CodeRepoBindingSpec, error) {
	_, project, err := c.resourcesUsecase.GetProductAndCodeRepo(ctx, productName)
	if err != nil {
		return nil, err
	}

	localPath, err := c.resourcesUsecase.CloneCodeRepo(ctx, project.HttpUrlToRepo)
	if err != nil {
		return nil, err
	}

	nodes, err := c.nodestree.Load(localPath)
	if err != nil {
		return nil, err
	}

	node := c.nodestree.GetNode(&nodes, nodestree.CodeRepoBinding, coderepoBindingName)
	if node == nil {
		return nil, nil
	}

	resouece, ok := node.Content.(*resourcev1alpha1.CodeRepoBinding)
	if !ok {
		return nil, fmt.Errorf("resource type is inconsistent, please check if this resource %s is legal", coderepoBindingName)
	}

	return &resouece.Spec, nil
}

func (c *CodeRepoBindingUsecase) AuthorizeDeployKey(ctx context.Context, spec resourcev1alpha1.CodeRepoBindingSpec) error {
	if err := c.ApplyDeploykey(ctx, spec, func(ctx context.Context, pid interface{}, deployKey int) error {
		_, err := c.codeRepo.EnableProjectDeployKey(ctx, pid, deployKey)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (c *CodeRepoBindingUsecase) RevokeDeployKey(ctx context.Context, spec resourcev1alpha1.CodeRepoBindingSpec) error {
	if err := c.ApplyDeploykey(ctx, spec, func(ctx context.Context, pid interface{}, deployKey int) error {
		err := c.codeRepo.DeleteDeployKey(ctx, pid, deployKey)
		if err != nil {
			if commonv1.IsDeploykeyNotFound(err) {
				return nil
			}
			return fmt.Errorf("failed to delete the code repository %s deploykey under group %s, err: %v", c.repositoryName, c.groupName, err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (c *CodeRepoBindingUsecase) ApplyDeploykey(ctx context.Context, spec resourcev1alpha1.CodeRepoBindingSpec, fn ApplyDeploykeyFunc) error {
	repositories, err := c.GetAuthorizedRepositories(ctx, spec)
	if err != nil {
		return err
	}

	for _, repository := range repositories {
		repoid := fmt.Sprintf("%s%d", RepoPrefix, repository.Id)
		secretData, err := c.GetDeployKeyFromSecretRepo(ctx, repoid, DefaultUser, spec.Permissions)
		if err != nil {
			return commonv1.ErrorDeploykeyNotFound("failed to get the %s deploykey from secret repo, please check if the key under /%s/%s exists or is invalid", spec.Permissions, c.config.Git.GitType, repoid)
		}

		deploykey, err := c.codeRepo.GetDeployKey(ctx, int(repository.Id), secretData.ID)
		if err != nil {
			return commonv1.ErrorDeploykeyNotFound("failed to get the %s deploykey from git, please check if the key exists or is invalid for the repository %s under organization %s", spec.Permissions, repository.Name, c.groupName)
		}

		pid, err := utilstrings.ExtractNumber(RepoPrefix, spec.CodeRepo)
		if err != nil {
			return err
		}

		//Revoke in-product authorization, not revoke the deploykey of the authorization repository.
		if pid == int(repository.Id) {
			continue
		}

		err = fn(ctx, pid, deploykey.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CodeRepoBindingUsecase) ClearLastDeployKeyResult(ctx context.Context, spec *resourcev1alpha1.CodeRepoBindingSpec) error {
	if err := c.RevokeDeployKey(ctx, *spec); err != nil {
		return err
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

	codeRepos, err := c.codeRepoUsecase.nodesToLists(*nodes, func(codeRepo *resourcev1alpha1.CodeRepo) bool {
		for _, project := range spec.Projects {
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

	objKey := client.ObjectKey{
		Namespace: codeRepoBinding.Spec.Product,
		Name:      codeRepoBinding.Spec.CodeRepo,
	}

	err := k8sClient.Get(context.TODO(), objKey, &resourcev1alpha1.CodeRepo{})
	if err != nil {
		return true, err
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
