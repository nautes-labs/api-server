package biz

import (
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/nautes-labs/api-server/pkg/kubernetes"
	"github.com/nautes-labs/api-server/pkg/nodestree"
	utilstrings "github.com/nautes-labs/api-server/util/string"
	resourcev1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Get CodeRepoBinding", func() {
	var (
		resourceName          = "codeRepoBinding1"
		fakeResource          = createFakeCodeRepoBindingResource(resourceName, "project1")
		fakeNode              = createFakeCodeRepoBindingNode(fakeResource)
		fakeNodes             = createFakeContainingCodeRepoBindingNodes(fakeNode)
		fakeErrorKindResource = createFakeCodeRepoBindingErrorKindResource(resourceName, "project1")
		fakeErrorKindNode     = createFakeCodeRepoNode(fakeErrorKindResource)
		fakeErrorKindNodes    = createFakeContainingCodeRepoBindingNodes(fakeErrorKindNode)
		gid, _                = utilstrings.ExtractNumber("product-", fakeResource.Spec.Product)
		project               = &Project{Id: 1222, HttpUrlToRepo: fmt.Sprintf("ssh://git@gitlab.io/nautes-labs/%s.git", "codeRepo1")}
		projectDeployKey      = &ProjectDeployKey{
			ID:  2013,
			Key: "FingerprintData",
		}
		cloneRepositoryParam = &CloneRepositoryParam{
			URL:   project.HttpUrlToRepo,
			User:  _GitUser,
			Email: _GitEmail,
		}
		deployKeySecretData = &DeployKeySecretData{
			ID:          2013,
			Fingerprint: "Fingerprint",
		}
	)

	It("get CodeRepoBinding successfully", func() {
		codeRepo := NewMockCodeRepo(ctl)
		codeRepo.EXPECT().GetGroup(gomock.Any(), fakeResource.Spec.Product).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), gid).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), defaultGroupName).Return(defaultProductGroup, nil).AnyTimes()

		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), gomock.Any()).Return(project, nil).AnyTimes()
		codeRepo.EXPECT().GetCurrentUser(gomock.Any()).Return(_GitUser, _GitEmail, nil).AnyTimes()
		codeRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(projectDeployKey, nil).AnyTimes()

		gitRepo := NewMockGitRepo(ctl)
		gitRepo.EXPECT().Clone(gomock.Any(), cloneRepositoryParam).Return(localRepositaryPath, nil).AnyTimes()

		nodestree := nodestree.NewMockNodesTree(ctl)
		nodestree.EXPECT().AppendOperators(gomock.Any()).AnyTimes()
		nodestree.EXPECT().Load(gomock.Eq(localRepositaryPath)).Return(fakeNodes, nil).AnyTimes()
		nodestree.EXPECT().Compare(gomock.Any()).Return(nil).AnyTimes()
		nodestree.EXPECT().GetNode(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeNode).AnyTimes()

		secretRepo := NewMockSecretrepo(ctl)
		secretRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any()).Return(deployKeySecretData, nil).AnyTimes()
		resourcesUsecase := NewResourcesUsecase(logger, codeRepo, nil, gitRepo, nodestree, nautesConfigs)
		client := kubernetes.NewMockClient(ctl)

		codeRepoUsecase := NewCodeRepoUsecase(logger, codeRepo, secretRepo, nodestree, nautesConfigs, resourcesUsecase, nil)
		biz := NewCodeRepoCodeRepoBindingUsecase(logger, codeRepo, secretRepo, nodestree, codeRepoUsecase, resourcesUsecase, nautesConfigs, client)
		options := &BizOptions{
			ProductName: defaultGroupName,
			ResouceName: resourceName,
		}
		item, err := biz.GetCodeRepoBinding(ctx, options)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(item.Spec).Should(Equal(fakeResource.Spec))
	})

	It("failed to get CodeRepoBinding", func() {
		codeRepo := NewMockCodeRepo(ctl)
		codeRepo.EXPECT().GetGroup(gomock.Any(), fakeResource.Spec.Product).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), gid).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), defaultGroupName).Return(defaultProductGroup, nil).AnyTimes()

		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), gomock.Any()).Return(project, nil).AnyTimes()
		codeRepo.EXPECT().GetCurrentUser(gomock.Any()).Return(_GitUser, _GitEmail, nil).AnyTimes()
		codeRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(projectDeployKey, nil).AnyTimes()

		gitRepo := NewMockGitRepo(ctl)
		gitRepo.EXPECT().Clone(gomock.Any(), cloneRepositoryParam).Return(localRepositaryPath, nil).AnyTimes()

		nodestree := nodestree.NewMockNodesTree(ctl)
		nodestree.EXPECT().AppendOperators(gomock.Any()).AnyTimes()
		nodestree.EXPECT().Load(gomock.Eq(localRepositaryPath)).Return(fakeErrorKindNodes, nil).AnyTimes()
		nodestree.EXPECT().Compare(gomock.Any()).Return(nil).AnyTimes()
		nodestree.EXPECT().GetNode(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeErrorKindNode).AnyTimes()

		secretRepo := NewMockSecretrepo(ctl)
		secretRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any()).Return(deployKeySecretData, nil).AnyTimes()
		resourcesUsecase := NewResourcesUsecase(logger, codeRepo, nil, gitRepo, nodestree, nautesConfigs)
		client := kubernetes.NewMockClient(ctl)

		codeRepoUsecase := NewCodeRepoUsecase(logger, codeRepo, secretRepo, nodestree, nautesConfigs, resourcesUsecase, nil)
		biz := NewCodeRepoCodeRepoBindingUsecase(logger, codeRepo, secretRepo, nodestree, codeRepoUsecase, resourcesUsecase, nautesConfigs, client)
		options := &BizOptions{
			ProductName: defaultGroupName,
			ResouceName: resourceName,
		}
		_, err := biz.GetCodeRepoBinding(ctx, options)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("List CodeRepoBinding", func() {
	var (
		resourceName          = "codeRepoBinding1"
		fakeResource          = createFakeCodeRepoBindingResource(resourceName, "project1")
		fakeNode              = createFakeCodeRepoBindingNode(fakeResource)
		fakeNodes             = createFakeContainingCodeRepoBindingNodes(fakeNode)
		fakeErrorKindResource = createFakeCodeRepoBindingErrorKindResource(resourceName, "project1")
		fakeErrorKindNode     = createFakeCodeRepoNode(fakeErrorKindResource)
		fakeErrorKindNodes    = createFakeContainingCodeRepoBindingNodes(fakeErrorKindNode)
		gid, _                = utilstrings.ExtractNumber("product-", fakeResource.Spec.Product)
		project               = &Project{Id: 1222, HttpUrlToRepo: fmt.Sprintf("ssh://git@gitlab.io/nautes-labs/%s.git", "codeRepo1")}
		projectDeployKey      = &ProjectDeployKey{
			ID:  2013,
			Key: "FingerprintData",
		}
		cloneRepositoryParam = &CloneRepositoryParam{
			URL:   project.HttpUrlToRepo,
			User:  _GitUser,
			Email: _GitEmail,
		}
		deployKeySecretData = &DeployKeySecretData{
			ID:          2013,
			Fingerprint: "Fingerprint",
		}
	)

	It("list CodeRepoBinding successfully", func() {
		codeRepo := NewMockCodeRepo(ctl)
		codeRepo.EXPECT().GetGroup(gomock.Any(), fakeResource.Spec.Product).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), gid).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), defaultGroupName).Return(defaultProductGroup, nil).AnyTimes()

		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), gomock.Any()).Return(project, nil).AnyTimes()
		codeRepo.EXPECT().GetCurrentUser(gomock.Any()).Return(_GitUser, _GitEmail, nil).AnyTimes()
		codeRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(projectDeployKey, nil).AnyTimes()

		gitRepo := NewMockGitRepo(ctl)
		gitRepo.EXPECT().Clone(gomock.Any(), cloneRepositoryParam).Return(localRepositaryPath, nil).AnyTimes()

		nodestree := nodestree.NewMockNodesTree(ctl)
		nodestree.EXPECT().AppendOperators(gomock.Any()).AnyTimes()
		nodestree.EXPECT().Load(gomock.Eq(localRepositaryPath)).Return(fakeNodes, nil).AnyTimes()
		nodestree.EXPECT().Compare(gomock.Any()).Return(nil).AnyTimes()
		nodestree.EXPECT().GetNode(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeNode).AnyTimes()

		secretRepo := NewMockSecretrepo(ctl)
		secretRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any()).Return(deployKeySecretData, nil).AnyTimes()
		resourcesUsecase := NewResourcesUsecase(logger, codeRepo, nil, gitRepo, nodestree, nautesConfigs)
		client := kubernetes.NewMockClient(ctl)

		codeRepoUsecase := NewCodeRepoUsecase(logger, codeRepo, secretRepo, nodestree, nautesConfigs, resourcesUsecase, nil)
		biz := NewCodeRepoCodeRepoBindingUsecase(logger, codeRepo, secretRepo, nodestree, codeRepoUsecase, resourcesUsecase, nautesConfigs, client)
		options := &BizOptions{
			ProductName: defaultGroupName,
			ResouceName: resourceName,
		}
		items, err := biz.ListCodeRepoBindings(ctx, options)
		Expect(err).ShouldNot(HaveOccurred())
		for _, item := range items {
			Expect(item.Spec).Should(Equal(fakeResource.Spec))
		}
	})

	It("failed to get CodeRepoBinding", func() {
		codeRepo := NewMockCodeRepo(ctl)
		codeRepo.EXPECT().GetGroup(gomock.Any(), fakeResource.Spec.Product).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), gid).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), defaultGroupName).Return(defaultProductGroup, nil).AnyTimes()

		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), gomock.Any()).Return(project, nil).AnyTimes()
		codeRepo.EXPECT().GetCurrentUser(gomock.Any()).Return(_GitUser, _GitEmail, nil).AnyTimes()
		codeRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(projectDeployKey, nil).AnyTimes()

		gitRepo := NewMockGitRepo(ctl)
		gitRepo.EXPECT().Clone(gomock.Any(), cloneRepositoryParam).Return(localRepositaryPath, nil).AnyTimes()

		nodestree := nodestree.NewMockNodesTree(ctl)
		nodestree.EXPECT().AppendOperators(gomock.Any()).AnyTimes()
		nodestree.EXPECT().Load(gomock.Eq(localRepositaryPath)).Return(fakeErrorKindNodes, nil).AnyTimes()
		nodestree.EXPECT().Compare(gomock.Any()).Return(nil).AnyTimes()
		nodestree.EXPECT().GetNode(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeErrorKindNode).AnyTimes()

		secretRepo := NewMockSecretrepo(ctl)
		secretRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any()).Return(deployKeySecretData, nil).AnyTimes()
		resourcesUsecase := NewResourcesUsecase(logger, codeRepo, nil, gitRepo, nodestree, nautesConfigs)
		client := kubernetes.NewMockClient(ctl)

		codeRepoUsecase := NewCodeRepoUsecase(logger, codeRepo, secretRepo, nodestree, nautesConfigs, resourcesUsecase, nil)
		biz := NewCodeRepoCodeRepoBindingUsecase(logger, codeRepo, secretRepo, nodestree, codeRepoUsecase, resourcesUsecase, nautesConfigs, client)
		options := &BizOptions{
			ProductName: defaultGroupName,
			ResouceName: resourceName,
		}
		_, err := biz.GetCodeRepoBinding(ctx, options)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("Save CodeRepoBinding", func() {
	var (
		resourceName         = "codeRepoBinding1"
		fakeResource1        = createFakeCodeRepoBindingResource(resourceName, "project1")
		fakeNode1            = createFakeCodeRepoBindingNode(fakeResource1)
		fakeNodes1           = createFakeContainingCodeRepoBindingNodes(fakeNode1)
		fakeResource2        = createFakeCodeRepoBindingResource(resourceName, "")
		fakeNode2            = createFakeCodeRepoBindingNode(fakeResource2)
		fakeNodes2           = createFakeContainingCodeRepoBindingNodes(fakeNode2)
		project1             = &Project{Id: 122, HttpUrlToRepo: fmt.Sprintf("ssh://git@gitlab.io/nautes-labs/%s.git", "codeRepo1")}
		project2             = &Project{Id: 123, HttpUrlToRepo: fmt.Sprintf("ssh://git@gitlab.io/nautes-labs/%s.git", "codeRepo2")}
		cloneRepositoryParam = &CloneRepositoryParam{
			URL:   project1.HttpUrlToRepo,
			User:  _GitUser,
			Email: _GitEmail,
		}
		projectDeployKey = &ProjectDeployKey{
			ID:  2013,
			Key: "FingerprintData",
		}
		deployKeySecretData = &DeployKeySecretData{
			ID:          2013,
			Fingerprint: "Fingerprint",
		}
		gid, _ = utilstrings.ExtractNumber("product-", fakeResource1.Spec.Product)
	)
	It("will create CodeRepoBinding according to the product", func() {
		codeRepo := NewMockCodeRepo(ctl)
		codeRepo.EXPECT().GetGroup(gomock.Any(), fakeResource1.Spec.Product).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), gid).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), defaultGroupName).Return(defaultProductGroup, nil).AnyTimes()

		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), defaultProjectPath).Return(project1, nil).AnyTimes()
		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), fmt.Sprintf("%s/%s", defaultGroupName, fakeResource1.Spec.CodeRepo)).Return(project1, nil).AnyTimes()
		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), 122).Return(project1, nil).AnyTimes()
		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), 123).Return(project2, nil).AnyTimes()
		codeRepo.EXPECT().GetCurrentUser(gomock.Any()).Return(_GitUser, _GitEmail, nil).AnyTimes()
		codeRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(projectDeployKey, nil).AnyTimes()
		codeRepo.EXPECT().DeleteDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		codeRepo.EXPECT().EnableProjectDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(projectDeployKey, nil).AnyTimes()

		gitRepo := NewMockGitRepo(ctl)
		gitRepo.EXPECT().Clone(gomock.Any(), cloneRepositoryParam).Return(localRepositaryPath, nil).AnyTimes()
		gitRepo.EXPECT().SaveConfig(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		gitRepo.EXPECT().Fetch(gomock.Any(), gomock.Any(), "origin").Return("any", nil).AnyTimes()
		gitRepo.EXPECT().Diff(gomock.Any(), gomock.Any(), "main", "remotes/origin/main").Return("", nil).AnyTimes()

		nodestree := nodestree.NewMockNodesTree(ctl)
		nodestree.EXPECT().AppendOperators(gomock.Any()).AnyTimes()
		nodestree.EXPECT().Load(gomock.Eq(localRepositaryPath)).Return(fakeNodes2, nil).AnyTimes()
		nodestree.EXPECT().Compare(gomock.Any()).Return(nil).AnyTimes()
		nodestree.EXPECT().GetNode(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		nodestree.EXPECT().InsertNodes(gomock.Any(), gomock.Any()).Return(&fakeNodes2, nil).AnyTimes()

		secretRepo := NewMockSecretrepo(ctl)
		secretRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any()).Return(deployKeySecretData, nil).AnyTimes()
		resourcesUsecase := NewResourcesUsecase(logger, codeRepo, nil, gitRepo, nodestree, nautesConfigs)
		client := kubernetes.NewMockClient(ctl)

		codeRepoUsecase := NewCodeRepoUsecase(logger, codeRepo, secretRepo, nodestree, nautesConfigs, resourcesUsecase, nil)
		biz := NewCodeRepoCodeRepoBindingUsecase(logger, codeRepo, secretRepo, nodestree, codeRepoUsecase, resourcesUsecase, nautesConfigs, client)
		options := &BizOptions{
			ProductName: defaultGroupName,
			ResouceName: resourceName,
		}
		data := &CodeRepoBindingData{
			Name: resourceName,
			Spec: fakeResource2.Spec,
		}
		err := biz.SaveCodeRepoBinding(ctx, options, data)
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("will update CodeRepoBinding according to the product", func() {
		codeRepo := NewMockCodeRepo(ctl)
		codeRepo.EXPECT().GetGroup(gomock.Any(), fakeResource1.Spec.Product).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), gid).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), defaultGroupName).Return(defaultProductGroup, nil).AnyTimes()

		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), defaultProjectPath).Return(project1, nil).AnyTimes()
		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), fmt.Sprintf("%s/%s", defaultGroupName, fakeResource1.Spec.CodeRepo)).Return(project1, nil).AnyTimes()
		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), 122).Return(project1, nil).AnyTimes()
		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), 123).Return(project2, nil).AnyTimes()
		codeRepo.EXPECT().GetCurrentUser(gomock.Any()).Return(_GitUser, _GitEmail, nil).AnyTimes()
		codeRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(projectDeployKey, nil).AnyTimes()
		codeRepo.EXPECT().DeleteDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		codeRepo.EXPECT().EnableProjectDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(projectDeployKey, nil).AnyTimes()

		gitRepo := NewMockGitRepo(ctl)
		gitRepo.EXPECT().Clone(gomock.Any(), cloneRepositoryParam).Return(localRepositaryPath, nil).AnyTimes()
		gitRepo.EXPECT().SaveConfig(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		gitRepo.EXPECT().Fetch(gomock.Any(), gomock.Any(), "origin").Return("any", nil).AnyTimes()
		gitRepo.EXPECT().Diff(gomock.Any(), gomock.Any(), "main", "remotes/origin/main").Return("", nil).AnyTimes()

		nodestree := nodestree.NewMockNodesTree(ctl)
		nodestree.EXPECT().AppendOperators(gomock.Any()).AnyTimes()
		nodestree.EXPECT().Load(gomock.Eq(localRepositaryPath)).Return(fakeNodes2, nil).AnyTimes()
		nodestree.EXPECT().Compare(gomock.Any()).Return(nil).AnyTimes()
		nodestree.EXPECT().GetNode(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeNode1).AnyTimes()
		nodestree.EXPECT().InsertNodes(gomock.Any(), gomock.Any()).Return(&fakeNodes2, nil).AnyTimes()

		secretRepo := NewMockSecretrepo(ctl)
		secretRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any()).Return(deployKeySecretData, nil).AnyTimes()
		resourcesUsecase := NewResourcesUsecase(logger, codeRepo, nil, gitRepo, nodestree, nautesConfigs)
		client := kubernetes.NewMockClient(ctl)

		codeRepoUsecase := NewCodeRepoUsecase(logger, codeRepo, secretRepo, nodestree, nautesConfigs, resourcesUsecase, nil)
		biz := NewCodeRepoCodeRepoBindingUsecase(logger, codeRepo, secretRepo, nodestree, codeRepoUsecase, resourcesUsecase, nautesConfigs, client)
		options := &BizOptions{
			ProductName: defaultGroupName,
			ResouceName: resourceName,
		}
		data := &CodeRepoBindingData{
			Name: resourceName,
			Spec: fakeResource2.Spec,
		}
		err := biz.SaveCodeRepoBinding(ctx, options, data)
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("will save CodeRepoBinding when specifying project", func() {
		codeRepo := NewMockCodeRepo(ctl)
		codeRepo.EXPECT().GetGroup(gomock.Any(), fakeResource1.Spec.Product).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), gid).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), defaultGroupName).Return(defaultProductGroup, nil).AnyTimes()

		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), defaultProjectPath).Return(project1, nil).AnyTimes()
		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), fmt.Sprintf("%s/%s", defaultGroupName, fakeResource1.Spec.CodeRepo)).Return(project1, nil).AnyTimes()
		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), 122).Return(project1, nil).AnyTimes()
		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), 123).Return(project2, nil).AnyTimes()
		codeRepo.EXPECT().GetCurrentUser(gomock.Any()).Return(_GitUser, _GitEmail, nil).AnyTimes()
		codeRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(projectDeployKey, nil).AnyTimes()
		codeRepo.EXPECT().DeleteDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		codeRepo.EXPECT().EnableProjectDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(projectDeployKey, nil).AnyTimes()

		gitRepo := NewMockGitRepo(ctl)
		gitRepo.EXPECT().Clone(gomock.Any(), cloneRepositoryParam).Return(localRepositaryPath, nil).AnyTimes()
		gitRepo.EXPECT().SaveConfig(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		gitRepo.EXPECT().Fetch(gomock.Any(), gomock.Any(), "origin").Return("any", nil).AnyTimes()
		gitRepo.EXPECT().Diff(gomock.Any(), gomock.Any(), "main", "remotes/origin/main").Return("", nil).AnyTimes()

		nodestree := nodestree.NewMockNodesTree(ctl)
		nodestree.EXPECT().AppendOperators(gomock.Any()).AnyTimes()
		nodestree.EXPECT().Load(gomock.Eq(localRepositaryPath)).Return(fakeNodes1, nil).AnyTimes()
		nodestree.EXPECT().Compare(gomock.Any()).Return(nil).AnyTimes()
		nodestree.EXPECT().GetNode(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeNode1).AnyTimes()
		nodestree.EXPECT().InsertNodes(gomock.Any(), gomock.Any()).Return(&fakeNodes1, nil).AnyTimes()

		secretRepo := NewMockSecretrepo(ctl)
		secretRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any()).Return(deployKeySecretData, nil).AnyTimes()
		resourcesUsecase := NewResourcesUsecase(logger, codeRepo, nil, gitRepo, nodestree, nautesConfigs)
		client := kubernetes.NewMockClient(ctl)

		codeRepoUsecase := NewCodeRepoUsecase(logger, codeRepo, secretRepo, nodestree, nautesConfigs, resourcesUsecase, nil)
		biz := NewCodeRepoCodeRepoBindingUsecase(logger, codeRepo, secretRepo, nodestree, codeRepoUsecase, resourcesUsecase, nautesConfigs, client)
		options := &BizOptions{
			ProductName: defaultGroupName,
			ResouceName: resourceName,
		}
		data := &CodeRepoBindingData{
			Name: resourceName,
			Spec: fakeResource1.Spec,
		}
		err := biz.SaveCodeRepoBinding(ctx, options, data)
		Expect(err).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("Delete CodeRepoBinding", func() {
	var (
		resourceName     = "codeRepoBinding1"
		fakeResource     = createFakeCodeRepoBindingResource(resourceName, "project1")
		fakeNode         = createFakeCodeRepoBindingNode(fakeResource)
		fakeNodes        = createFakeContainingCodeRepoBindingNodes(fakeNode)
		gid, _           = utilstrings.ExtractNumber("product-", fakeResource.Spec.Product)
		project          = &Project{Id: 1222, HttpUrlToRepo: fmt.Sprintf("ssh://git@gitlab.io/nautes-labs/%s.git", "codeRepo1")}
		projectDeployKey = &ProjectDeployKey{
			ID:  2013,
			Key: "FingerprintData",
		}
		cloneRepositoryParam = &CloneRepositoryParam{
			URL:   project.HttpUrlToRepo,
			User:  _GitUser,
			Email: _GitEmail,
		}
		deployKeySecretData = &DeployKeySecretData{
			ID:          2013,
			Fingerprint: "Fingerprint",
		}
	)

	It("delete CodeRepoBinding successfully", func() {
		codeRepo := NewMockCodeRepo(ctl)
		codeRepo.EXPECT().GetGroup(gomock.Any(), fakeResource.Spec.Product).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), gid).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), defaultGroupName).Return(defaultProductGroup, nil).AnyTimes()

		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), gomock.Any()).Return(project, nil).AnyTimes()
		codeRepo.EXPECT().GetCurrentUser(gomock.Any()).Return(_GitUser, _GitEmail, nil).AnyTimes()
		codeRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(projectDeployKey, nil).AnyTimes()
		codeRepo.EXPECT().DeleteDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		gitRepo := NewMockGitRepo(ctl)
		gitRepo.EXPECT().Clone(gomock.Any(), cloneRepositoryParam).Return(localRepositaryPath, nil).AnyTimes()
		gitRepo.EXPECT().SaveConfig(gomock.Any(), gomock.Any()).Return(nil)
		gitRepo.EXPECT().Fetch(gomock.Any(), gomock.Any(), "origin").Return("any", nil)
		gitRepo.EXPECT().Diff(gomock.Any(), gomock.Any(), "main", "remotes/origin/main").Return("", nil)

		nodestree := nodestree.NewMockNodesTree(ctl)
		nodestree.EXPECT().AppendOperators(gomock.Any()).AnyTimes()
		nodestree.EXPECT().Load(gomock.Eq(localRepositaryPath)).Return(fakeNodes, nil).AnyTimes()
		nodestree.EXPECT().Compare(gomock.Any()).Return(nil).AnyTimes()
		nodestree.EXPECT().GetNode(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeNode).AnyTimes()
		nodestree.EXPECT().RemoveNode(&fakeNodes, fakeNode).Return(&emptyNodes, nil)

		secretRepo := NewMockSecretrepo(ctl)
		secretRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any()).Return(deployKeySecretData, nil).AnyTimes()
		resourcesUsecase := NewResourcesUsecase(logger, codeRepo, nil, gitRepo, nodestree, nautesConfigs)
		client := kubernetes.NewMockClient(ctl)

		codeRepoUsecase := NewCodeRepoUsecase(logger, codeRepo, secretRepo, nodestree, nautesConfigs, resourcesUsecase, nil)
		biz := NewCodeRepoCodeRepoBindingUsecase(logger, codeRepo, secretRepo, nodestree, codeRepoUsecase, resourcesUsecase, nautesConfigs, client)
		options := &BizOptions{
			ProductName: defaultGroupName,
			ResouceName: resourceName,
		}
		err := biz.DeleteCodeRepoBinding(ctx, options)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("failed to delete deploykey", func() {
		codeRepo := NewMockCodeRepo(ctl)
		codeRepo.EXPECT().GetGroup(gomock.Any(), fakeResource.Spec.Product).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), gid).Return(defaultProductGroup, nil).AnyTimes()
		codeRepo.EXPECT().GetGroup(gomock.Any(), defaultGroupName).Return(defaultProductGroup, nil).AnyTimes()

		codeRepo.EXPECT().GetCodeRepo(gomock.Any(), gomock.Any()).Return(project, nil).AnyTimes()
		codeRepo.EXPECT().GetCurrentUser(gomock.Any()).Return(_GitUser, _GitEmail, nil).AnyTimes()
		codeRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(projectDeployKey, nil).AnyTimes()
		codeRepo.EXPECT().DeleteDeployKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("failed to delete deploykey."))

		gitRepo := NewMockGitRepo(ctl)
		gitRepo.EXPECT().Clone(gomock.Any(), cloneRepositoryParam).Return(localRepositaryPath, nil).AnyTimes()

		nodestree := nodestree.NewMockNodesTree(ctl)
		nodestree.EXPECT().AppendOperators(gomock.Any()).AnyTimes()
		nodestree.EXPECT().Load(gomock.Eq(localRepositaryPath)).Return(fakeNodes, nil).AnyTimes()
		nodestree.EXPECT().Compare(gomock.Any()).Return(nil).AnyTimes()
		nodestree.EXPECT().GetNode(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeNode).AnyTimes()

		secretRepo := NewMockSecretrepo(ctl)
		secretRepo.EXPECT().GetDeployKey(gomock.Any(), gomock.Any()).Return(deployKeySecretData, nil).AnyTimes()
		resourcesUsecase := NewResourcesUsecase(logger, codeRepo, nil, gitRepo, nodestree, nautesConfigs)
		client := kubernetes.NewMockClient(ctl)

		codeRepoUsecase := NewCodeRepoUsecase(logger, codeRepo, secretRepo, nodestree, nautesConfigs, resourcesUsecase, nil)
		biz := NewCodeRepoCodeRepoBindingUsecase(logger, codeRepo, secretRepo, nodestree, codeRepoUsecase, resourcesUsecase, nautesConfigs, client)
		options := &BizOptions{
			ProductName: defaultGroupName,
			ResouceName: resourceName,
		}
		err := biz.DeleteCodeRepoBinding(ctx, options)
		Expect(err).Should(HaveOccurred())
	})
})

func createFakeCodeRepoBindingErrorKindResource(name, project string) *resourcev1alpha1.CodeRepo {
	crd := &resourcev1alpha1.CodeRepo{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		TypeMeta: v1.TypeMeta{
			Kind: nodestree.CodeRepo,
		},
	}

	return crd
}

func createFakeCodeRepoBindingResource(name, project string) *resourcev1alpha1.CodeRepoBinding {
	crd := &resourcev1alpha1.CodeRepoBinding{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		TypeMeta: v1.TypeMeta{
			Kind: nodestree.CodeRepoBinding,
		},
		Spec: resourcev1alpha1.CodeRepoBindingSpec{
			Product:     defaultProductId,
			CodeRepo:    "repo-122",
			Projects:    []string{},
			Permissions: "readonly",
		},
	}
	if project != "" {
		crd.Spec.Projects = append(crd.Spec.Projects, project)
	}
	return crd
}

func createFakeCodeRepoBindingNode(resource *resourcev1alpha1.CodeRepoBinding) *nodestree.Node {
	return &nodestree.Node{
		Name:    resource.Name,
		Path:    fmt.Sprintf("%s/%s/%s/%s.yaml", localRepositaryPath, _CodeReposSubDir, resource.Name, resource.Name),
		Level:   4,
		Content: resource,
		Kind:    nodestree.CodeRepoBinding,
	}
}

func createFakeContainingCodeRepoBindingNodes(node *nodestree.Node) nodestree.Node {
	projectID1 := fmt.Sprintf("%s%d", RepoPrefix, 122)
	codeRepoNode1 := createFakeCodeRepoNode(createFakeCodeRepoResource(projectID1))
	projectID2 := fmt.Sprintf("%s%d", RepoPrefix, 123)
	codeRepoNod2 := createFakeCodeRepoNode(createFakeCodeRepoResource(projectID2))
	return nodestree.Node{
		Name:  defaultProjectName,
		Path:  defaultProjectName,
		IsDir: true,
		Level: 1,
		Children: []*nodestree.Node{
			{
				Name:  _CodeReposSubDir,
				Path:  fmt.Sprintf("%v/%v", defaultProjectName, _CodeReposSubDir),
				IsDir: true,
				Level: 2,
				Children: []*nodestree.Node{
					{
						Name:  node.Name,
						Path:  fmt.Sprintf("%s/%s/%s", localRepositaryPath, _CodeReposSubDir, codeRepoNode1.Name),
						IsDir: true,
						Level: 3,
						Children: []*nodestree.Node{
							node,
							codeRepoNode1,
							codeRepoNod2,
						},
					},
				},
			},
		},
	}
}
