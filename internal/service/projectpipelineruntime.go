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

	projectpipelineruntimev1 "github.com/nautes-labs/api-server/api/projectpipelineruntime/v1"
	"github.com/nautes-labs/api-server/internal/biz"
	resourcev1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
)

type ProjectPipelineRuntimeService struct {
	projectpipelineruntimev1.UnimplementedProjectPipelineRuntimeServer
	projectPipelineRuntime *biz.ProjectPipelineRuntimeUsecase
}

func NewProjectPipelineRuntimeService(projectPipelineRuntime *biz.ProjectPipelineRuntimeUsecase) *ProjectPipelineRuntimeService {
	return &ProjectPipelineRuntimeService{projectPipelineRuntime: projectPipelineRuntime}
}

func (s *ProjectPipelineRuntimeService) covertCodeRepoValueToReply(projectPipelineRuntime *resourcev1alpha1.ProjectPipelineRuntime, productName string) *projectpipelineruntimev1.GetReply {
	var pipelines []*projectpipelineruntimev1.Pipeline
	for _, pipeline := range projectPipelineRuntime.Spec.Pipelines {
		pipelines = append(pipelines, &projectpipelineruntimev1.Pipeline{
			Name:  pipeline.Name,
			Label: pipeline.Label,
			Path:  pipeline.Path,
		})
	}

	var eventSources []*projectpipelineruntimev1.EventSource
	for _, source := range projectPipelineRuntime.Spec.EventSources {
		eventSources = append(eventSources, &projectpipelineruntimev1.EventSource{
			Name: source.Name,
			Gitlab: &projectpipelineruntimev1.Gitlab{
				RepoName: source.Gitlab.RepoName,
				Revision: source.Gitlab.Revision,
				Events:   source.Gitlab.Events,
			},
		})
	}

	var pipelineTriggers []*projectpipelineruntimev1.PipelineTriggers
	for _, trigger := range projectPipelineRuntime.Spec.PipelineTriggers {
		pipelineTriggers = append(pipelineTriggers, &projectpipelineruntimev1.PipelineTriggers{
			EventSource: trigger.EventSource,
			Pipeline:    trigger.Pipeline,
			Revision:    trigger.Revision,
		})
	}

	return &projectpipelineruntimev1.GetReply{
		Name:             projectPipelineRuntime.Name,
		Project:          projectPipelineRuntime.Spec.Project,
		PipelineSource:   projectPipelineRuntime.Spec.PipelineSource,
		EventSources:     eventSources,
		Pipelines:        pipelines,
		PipelineTriggers: pipelineTriggers,
		Destination:      projectPipelineRuntime.Spec.Destination,
	}
}

func (s *ProjectPipelineRuntimeService) GetProjectPipelineRuntime(ctx context.Context, req *projectpipelineruntimev1.GetRequest) (*projectpipelineruntimev1.GetReply, error) {
	runtime, err := s.projectPipelineRuntime.GetProjectPipelineRuntime(ctx, req.ProjectPipelineRuntimeName, req.ProductName)
	if err != nil {
		return nil, err
	}

	return s.covertCodeRepoValueToReply(runtime, req.ProductName), nil
}

func (s *ProjectPipelineRuntimeService) ListProjectPipelineRuntimes(ctx context.Context, req *projectpipelineruntimev1.ListsRequest) (*projectpipelineruntimev1.ListsReply, error) {
	runtimes, err := s.projectPipelineRuntime.ListProjectPipelineRuntimes(ctx, req.ProductName)
	if err != nil {
		return nil, err
	}

	var items []*projectpipelineruntimev1.GetReply
	for _, runtime := range runtimes {
		items = append(items, s.covertCodeRepoValueToReply(runtime, req.ProductName))
	}

	return &projectpipelineruntimev1.ListsReply{
		Items: items,
	}, nil
}

func (s *ProjectPipelineRuntimeService) SaveProjectPipelineRuntime(ctx context.Context, req *projectpipelineruntimev1.SaveRequest) (*projectpipelineruntimev1.SaveReply, error) {
	data := &biz.ProjectPipelineRuntimeData{
		Name: req.ProjectPipelineRuntimeName,
		Spec: resourcev1alpha1.ProjectPipelineRuntimeSpec{
			Project:          req.Body.Project,
			PipelineSource:   req.Body.PipelineSource,
			EventSources:     s.convertEventSources(req.Body.EventSources),
			Pipelines:        s.convertPipelines(req.Body.Pipelines),
			PipelineTriggers: s.convertPipelineTriggers(req.Body.PipelineTriggers),
			Destination:      req.Body.Destination,
			Isolation:        req.Body.Isolation,
		},
	}

	options := &biz.BizOptions{
		ResouceName:       req.ProjectPipelineRuntimeName,
		ProductName:       req.ProductName,
		InsecureSkipCheck: req.InsecureSkipCheck,
	}
	err := s.projectPipelineRuntime.SaveProjectPipelineRuntime(ctx, options, data)
	if err != nil {
		return nil, err
	}

	return &projectpipelineruntimev1.SaveReply{
		Msg: fmt.Sprintf("Successfully saved %s configuration", req.ProjectPipelineRuntimeName),
	}, nil
}

func (s *ProjectPipelineRuntimeService) convertPipelineTriggers(triggers []*projectpipelineruntimev1.PipelineTriggers) (pipelineTriggers []resourcev1alpha1.PipelineTrigger) {
	for _, trigger := range triggers {
		resourcePipelineTrigger := resourcev1alpha1.PipelineTrigger{
			EventSource: trigger.EventSource,
			Pipeline:    trigger.Pipeline,
			Revision:    trigger.Revision,
		}

		pipelineTriggers = append(pipelineTriggers, resourcePipelineTrigger)
	}

	return pipelineTriggers
}

func (s *ProjectPipelineRuntimeService) convertPipelines(pipelines []*projectpipelineruntimev1.Pipeline) (resourcePipelines []resourcev1alpha1.Pipeline) {

	for _, pipeline := range pipelines {
		resourcePipeline := resourcev1alpha1.Pipeline{
			Name:  pipeline.Name,
			Label: pipeline.Label,
			Path:  pipeline.Path,
		}

		resourcePipelines = append(resourcePipelines, resourcePipeline)

	}

	return resourcePipelines
}

func (s *ProjectPipelineRuntimeService) convertEventSources(events []*projectpipelineruntimev1.EventSource) (eventSources []resourcev1alpha1.EventSource) {
	for _, eventSource := range events {
		resourceEventSource := resourcev1alpha1.EventSource{
			Name: eventSource.Name,
			Gitlab: &resourcev1alpha1.Gitlab{
				RepoName: eventSource.Gitlab.RepoName,
				Revision: eventSource.Gitlab.Revision,
				Events:   eventSource.Gitlab.Events,
			},
			Calendar: &resourcev1alpha1.Calendar{
				Schedule:       eventSource.Calendar.Schedule,
				Interval:       eventSource.Calendar.Interval,
				ExclusionDates: eventSource.Calendar.ExclusionDates,
				Timezone:       eventSource.Calendar.Timezone,
			},
		}

		eventSources = append(eventSources, resourceEventSource)
	}

	return eventSources
}

func (s *ProjectPipelineRuntimeService) DeleteProjectPipelineRuntime(ctx context.Context, req *projectpipelineruntimev1.DeleteRequest) (*projectpipelineruntimev1.DeleteReply, error) {
	options := &biz.BizOptions{
		ResouceName:       req.ProjectPipelineRuntimeName,
		ProductName:       req.ProductName,
		InsecureSkipCheck: req.InsecureSkipCheck,
	}
	err := s.projectPipelineRuntime.DeleteProjectPipelineRuntime(ctx, options)
	if err != nil {
		return nil, err
	}

	return &projectpipelineruntimev1.DeleteReply{
		Msg: fmt.Sprintf("Successfully deleted %s configuration", req.ProjectPipelineRuntimeName),
	}, nil
}
