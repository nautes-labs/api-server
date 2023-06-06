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

// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/validate/interface.go

// Package validate_test is a generated GoMock package.
package validate

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
)

// MockValidateClient is a mock of ValidateClient interface.
type MockValidateClient struct {
	ctrl     *gomock.Controller
	recorder *MockValidateClientMockRecorder
}

// MockValidateClientMockRecorder is the mock recorder for MockValidateClient.
type MockValidateClientMockRecorder struct {
	mock *MockValidateClient
}

// NewMockValidateClient creates a new mock instance.
func NewMockValidateClient(ctrl *gomock.Controller) *MockValidateClient {
	mock := &MockValidateClient{ctrl: ctrl}
	mock.recorder = &MockValidateClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockValidateClient) EXPECT() *MockValidateClientMockRecorder {
	return m.recorder
}

// GetCodeRepo mocks base method.
func (m *MockValidateClient) GetCodeRepo(ctx context.Context, repoName string) (*v1alpha1.CodeRepo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCodeRepo", ctx, repoName)
	ret0, _ := ret[0].(*v1alpha1.CodeRepo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCodeRepo indicates an expected call of GetCodeRepo.
func (mr *MockValidateClientMockRecorder) GetCodeRepo(ctx, repoName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCodeRepo", reflect.TypeOf((*MockValidateClient)(nil).GetCodeRepo), ctx, repoName)
}

// ListCodeRepoBinding mocks base method.
func (m *MockValidateClient) ListCodeRepoBinding(ctx context.Context, productName, repoName string) (*v1alpha1.CodeRepoBindingList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCodeRepoBinding", ctx, productName, repoName)
	ret0, _ := ret[0].(*v1alpha1.CodeRepoBindingList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCodeRepoBinding indicates an expected call of ListCodeRepoBinding.
func (mr *MockValidateClientMockRecorder) ListCodeRepoBinding(ctx, productName, repoName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCodeRepoBinding", reflect.TypeOf((*MockValidateClient)(nil).ListCodeRepoBinding), ctx, productName, repoName)
}

// ListDeploymentRuntime mocks base method.
func (m *MockValidateClient) ListDeploymentRuntime(ctx context.Context, productName string) (*v1alpha1.DeploymentRuntimeList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListDeploymentRuntime", ctx, productName)
	ret0, _ := ret[0].(*v1alpha1.DeploymentRuntimeList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListDeploymentRuntime indicates an expected call of ListDeploymentRuntime.
func (mr *MockValidateClientMockRecorder) ListDeploymentRuntime(ctx, productName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDeploymentRuntime", reflect.TypeOf((*MockValidateClient)(nil).ListDeploymentRuntime), ctx, productName)
}

// ListProjectPipelineRuntime mocks base method.
func (m *MockValidateClient) ListProjectPipelineRuntime(ctx context.Context, productName string) (*v1alpha1.ProjectPipelineRuntimeList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProjectPipelineRuntime", ctx, productName)
	ret0, _ := ret[0].(*v1alpha1.ProjectPipelineRuntimeList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProjectPipelineRuntime indicates an expected call of ListProjectPipelineRuntime.
func (mr *MockValidateClientMockRecorder) ListProjectPipelineRuntime(ctx, productName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProjectPipelineRuntime", reflect.TypeOf((*MockValidateClient)(nil).ListProjectPipelineRuntime), ctx, productName)
}
