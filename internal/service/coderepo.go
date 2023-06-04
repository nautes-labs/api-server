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
	"encoding/json"
	"errors"
	"fmt"

	coderepov1 "github.com/nautes-labs/api-server/api/coderepo/v1"
	commonv1 "github.com/nautes-labs/api-server/api/common/v1"
	"github.com/nautes-labs/api-server/internal/biz"
	"github.com/nautes-labs/api-server/pkg/nodestree"
	resourcev1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
	nautesconfigs "github.com/nautes-labs/pkg/pkg/nautesconfigs"
)

type CodeRepoService struct {
	coderepov1.UnimplementedCodeRepoServer
	codeRepo *biz.CodeRepoUsecase
	configs  *nautesconfigs.Config
}

func NewCodeRepoService(codeRepo *biz.CodeRepoUsecase, configs *nautesconfigs.Config) *CodeRepoService {
	return &CodeRepoService{
		codeRepo: codeRepo,
		configs:  configs,
	}
}

func (s *CodeRepoService) CovertCodeRepoValueToReply(codeRepo *resourcev1alpha1.CodeRepo, project *biz.Project) *coderepov1.GetReply {
	var git *coderepov1.GitProject
	if s.configs.Git.GitType == nautesconfigs.GIT_TYPE_GITLAB {
		git = &coderepov1.GitProject{
			Gitlab: &coderepov1.GitlabProject{
				Name:          project.Name,
				Description:   project.Description,
				Path:          project.Path,
				Visibility:    project.Visibility,
				HttpUrlToRepo: project.HttpUrlToRepo,
				SshUrlToRepo:  project.SshUrlToRepo,
			},
		}
	} else {
		git = &coderepov1.GitProject{
			Github: &coderepov1.GithubProject{
				Name:          project.Name,
				Description:   project.Description,
				Path:          project.Path,
				Visibility:    project.Visibility,
				HttpUrlToRepo: project.HttpUrlToRepo,
				SshUrlToRepo:  project.SshUrlToRepo,
			},
		}
	}

	return &coderepov1.GetReply{
		Product: codeRepo.Spec.Product,
		Name:    codeRepo.Spec.RepoName,
		Project: codeRepo.Spec.Project,
		Webhook: &coderepov1.Webhook{
			Events: codeRepo.Spec.Webhook.Events,
		},
		PipelineRuntime:   codeRepo.Spec.PipelineRuntime,
		DeploymentRuntime: codeRepo.Spec.DeploymentRuntime,
		Git:               git,
	}
}

func (s *CodeRepoService) GetCodeRepo(ctx context.Context, req *coderepov1.GetRequest) (*coderepov1.GetReply, error) {
	codeRepo, project, err := s.codeRepo.GetCodeRepo(ctx, req.CoderepoName, req.ProductName)
	if err != nil {
		return nil, err
	}

	return s.CovertCodeRepoValueToReply(codeRepo, project), nil
}

func (s *CodeRepoService) ListCodeRepos(ctx context.Context, req *coderepov1.ListsRequest) (*coderepov1.ListsReply, error) {
	codeRepoAndProjects, err := s.codeRepo.ListCodeRepos(ctx, req.ProductName)
	if err != nil {
		return nil, err
	}

	var items []*coderepov1.GetReply
	for _, cp := range codeRepoAndProjects {
		items = append(items, s.CovertCodeRepoValueToReply(cp.CodeRepo, cp.Project))
	}

	return &coderepov1.ListsReply{
		Items: items,
	}, nil
}

func (s *CodeRepoService) SaveCodeRepo(ctx context.Context, req *coderepov1.SaveRequest) (*coderepov1.SaveReply, error) {
	ctx = biz.SetResourceContext(ctx, "", biz.SaveMethod, "", "", nodestree.CodeRepo, req.CoderepoName)

	gitOptions := &biz.GitCodeRepoOptions{
		Gitlab: &biz.GitlabCodeRepoOptions{},
	}

	// TODO
	// Coming soon to support github
	if s.configs.Git.GitType == nautesconfigs.GIT_TYPE_GITLAB {
		bytes, err := json.Marshal(req.Body.Git.Gitlab)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(bytes, gitOptions.Gitlab)
		if err != nil {
			return nil, err
		}

		if gitOptions.Gitlab.Name == "" {
			gitOptions.Gitlab.Name = req.CoderepoName
		}

		if gitOptions.Gitlab.Path == "" {
			gitOptions.Gitlab.Path = gitOptions.Gitlab.Name
		}

		if gitOptions.Gitlab.Path != req.CoderepoName {
			return nil, fmt.Errorf("the name of the codeRepo resource must be the same as the repository path")
		}
	} else {
		if gitOptions.Github != "" {
			return nil, errors.New("coming soon to support github")
		}
	}

	data := &biz.CodeRepoData{
		Spec: resourcev1alpha1.CodeRepoSpec{
			Project:           req.Body.Project,
			RepoName:          req.CoderepoName,
			DeploymentRuntime: req.Body.DeploymentRuntime,
			PipelineRuntime:   req.Body.PipelineRuntime,
			Webhook: &resourcev1alpha1.Webhook{
				Events: req.Body.Webhook.Events,
			},
		},
	}
	options := &biz.BizOptions{
		ResouceName:       req.CoderepoName,
		ProductName:       req.ProductName,
		InsecureSkipCheck: req.InsecureSkipCheck,
	}

	err := s.codeRepo.SaveCodeRepo(ctx, options, data, gitOptions)
	if err != nil {
		if commonv1.IsRefreshPermissionsAccessDenied(err) {
			return &coderepov1.SaveReply{
				Msg: fmt.Sprintf("[warning] the CodeRepo %s was deleted, but refresh permission failed, please check the current user permission, err: %s", req.CoderepoName, err.Error()),
			}, nil
		}

		return nil, err
	}

	return &coderepov1.SaveReply{
		Msg: fmt.Sprintf("Successfully saved %v configuration", req.CoderepoName),
	}, nil
}

func (s *CodeRepoService) DeleteCodeRepo(ctx context.Context, req *coderepov1.DeleteRequest) (*coderepov1.DeleteReply, error) {
	ctx = biz.SetResourceContext(ctx, "", biz.DeleteMethod, "", "", nodestree.CodeRepo, req.CoderepoName)

	options := &biz.BizOptions{
		ResouceName:       req.CoderepoName,
		ProductName:       req.ProductName,
		InsecureSkipCheck: req.InsecureSkipCheck,
	}
	err := s.codeRepo.DeleteCodeRepo(ctx, options)
	if err != nil {
		return nil, err
	}

	return &coderepov1.DeleteReply{
		Msg: fmt.Sprintf("Successfully deleted %v configuration", req.CoderepoName),
	}, nil
}
