package service

import (
	"context"
	"fmt"

	coderepobindingv1 "github.com/nautes-labs/api-server/api/coderepobinding/v1"
	"github.com/nautes-labs/api-server/internal/biz"
	resourcev1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
)

type CodeRepoBindingService struct {
	coderepobindingv1.UnimplementedCodeRepoBindingServer
	codeRepoBindingUsecase *biz.CodeRepoBindingUsecase
}

func NewCodeRepoBindingService(codeRepoBindingUsecase *biz.CodeRepoBindingUsecase) *CodeRepoBindingService {
	return &CodeRepoBindingService{codeRepoBindingUsecase: codeRepoBindingUsecase}
}

func (s *CodeRepoBindingService) CovertCodeRepoBindingValueToReply(codeRepoBinding *resourcev1alpha1.CodeRepoBinding) *coderepobindingv1.GetReply {
	return &coderepobindingv1.GetReply{
		Name:        codeRepoBinding.Name,
		Coderepo:    codeRepoBinding.Spec.CodeRepo,
		Product:     codeRepoBinding.Spec.Product,
		Projects:    codeRepoBinding.Spec.Projects,
		Permissions: codeRepoBinding.Spec.Permissions,
	}
}

func (s *CodeRepoBindingService) GetCodeRepoBinding(ctx context.Context, req *coderepobindingv1.GetRequest) (*coderepobindingv1.GetReply, error) {
	options := &biz.BizOptions{
		ProductName: req.ProductName,
		ResouceName: req.CoderepoBindingName,
	}
	codeRepoBinding, err := s.codeRepoBindingUsecase.GetCodeRepoBinding(ctx, options)
	if err != nil {
		return nil, err
	}

	return s.CovertCodeRepoBindingValueToReply(codeRepoBinding), nil
}

func (s *CodeRepoBindingService) ListCodeRepoBindings(ctx context.Context, req *coderepobindingv1.ListsRequest) (*coderepobindingv1.ListsReply, error) {
	options := &biz.BizOptions{
		ProductName: req.ProductName,
	}

	codeRepoBindings, err := s.codeRepoBindingUsecase.ListCodeRepoBindings(ctx, options)
	if err != nil {
		return nil, err
	}

	var items []*coderepobindingv1.GetReply
	for _, codeRepoBinding := range codeRepoBindings {
		items = append(items, s.CovertCodeRepoBindingValueToReply(codeRepoBinding))
	}

	return &coderepobindingv1.ListsReply{Items: items}, nil
}

func (s *CodeRepoBindingService) SaveCodeRepoBinding(ctx context.Context, req *coderepobindingv1.SaveRequest) (*coderepobindingv1.SaveReply, error) {
	options := &biz.BizOptions{
		ProductName:       req.ProductName,
		ResouceName:       req.CoderepoBindingName,
		InsecureSkipCheck: req.InsecureSkipCheck,
	}

	data := &biz.CodeRepoBindingData{
		Name: req.CoderepoBindingName,
		Spec: resourcev1alpha1.CodeRepoBindingSpec{
			CodeRepo:    req.Body.Coderepo,
			Product:     req.Body.Product,
			Projects:    req.Body.Projects,
			Permissions: req.Body.Permissions,
		},
	}
	if err := s.codeRepoBindingUsecase.SaveCodeRepoBinding(ctx, options, data); err != nil {
		return nil, err
	}

	return &coderepobindingv1.SaveReply{
		Msg: fmt.Sprintf("Successfully saved %v configuration", req.CoderepoBindingName),
	}, nil
}

func (s *CodeRepoBindingService) DeleteCodeRepoBinding(ctx context.Context, req *coderepobindingv1.DeleteRequest) (*coderepobindingv1.DeleteReply, error) {
	options := &biz.BizOptions{
		ProductName:       req.ProductName,
		ResouceName:       req.CoderepoBindingName,
		InsecureSkipCheck: req.InsecureSkipCheck,
	}

	if err := s.codeRepoBindingUsecase.DeleteCodeRepoBinding(ctx, options); err != nil {
		return nil, err
	}

	return &coderepobindingv1.DeleteReply{
		Msg: fmt.Sprintf("Successfully deleted %v configuration", req.CoderepoBindingName),
	}, nil
}
