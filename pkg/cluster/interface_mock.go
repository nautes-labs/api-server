// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/cluster/interface.go

// Package cluster is a generated GoMock package.
package cluster

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
)

// MockClusterRegistrationOperator is a mock of ClusterRegistrationOperator interface.
type MockClusterRegistrationOperator struct {
	ctrl     *gomock.Controller
	recorder *MockClusterRegistrationOperatorMockRecorder
}

// MockClusterRegistrationOperatorMockRecorder is the mock recorder for MockClusterRegistrationOperator.
type MockClusterRegistrationOperatorMockRecorder struct {
	mock *MockClusterRegistrationOperator
}

// NewMockClusterRegistrationOperator creates a new mock instance.
func NewMockClusterRegistrationOperator(ctrl *gomock.Controller) *MockClusterRegistrationOperator {
	mock := &MockClusterRegistrationOperator{ctrl: ctrl}
	mock.recorder = &MockClusterRegistrationOperatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterRegistrationOperator) EXPECT() *MockClusterRegistrationOperatorMockRecorder {
	return m.recorder
}

// GetArgocdURL mocks base method.
func (m *MockClusterRegistrationOperator) GetArgocdURL() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetArgocdURL")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetArgocdURL indicates an expected call of GetArgocdURL.
func (mr *MockClusterRegistrationOperatorMockRecorder) GetArgocdURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArgocdURL", reflect.TypeOf((*MockClusterRegistrationOperator)(nil).GetArgocdURL))
}

// GetClsuter mocks base method.
func (m *MockClusterRegistrationOperator) GetClsuter(tenantLocalPath, clusterName string) (*v1alpha1.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClsuter", tenantLocalPath, clusterName)
	ret0, _ := ret[0].(*v1alpha1.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClsuter indicates an expected call of GetClsuter.
func (mr *MockClusterRegistrationOperatorMockRecorder) GetClsuter(tenantLocalPath, clusterName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClsuter", reflect.TypeOf((*MockClusterRegistrationOperator)(nil).GetClsuter), tenantLocalPath, clusterName)
}

// GetClsuters mocks base method.
func (m *MockClusterRegistrationOperator) GetClsuters(tenantLocalPath string) ([]*v1alpha1.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClsuters", tenantLocalPath)
	ret0, _ := ret[0].([]*v1alpha1.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClsuters indicates an expected call of GetClsuters.
func (mr *MockClusterRegistrationOperatorMockRecorder) GetClsuters(tenantLocalPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClsuters", reflect.TypeOf((*MockClusterRegistrationOperator)(nil).GetClsuters), tenantLocalPath)
}

// GetTektonOAuthURL mocks base method.
func (m *MockClusterRegistrationOperator) GetTektonOAuthURL() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTektonOAuthURL")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTektonOAuthURL indicates an expected call of GetTektonOAuthURL.
func (mr *MockClusterRegistrationOperatorMockRecorder) GetTektonOAuthURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTektonOAuthURL", reflect.TypeOf((*MockClusterRegistrationOperator)(nil).GetTektonOAuthURL))
}

// InitializeClusterConfig mocks base method.
func (m *MockClusterRegistrationOperator) InitializeClusterConfig(param *ClusterRegistrationParam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitializeClusterConfig", param)
	ret0, _ := ret[0].(error)
	return ret0
}

// InitializeClusterConfig indicates an expected call of InitializeClusterConfig.
func (mr *MockClusterRegistrationOperatorMockRecorder) InitializeClusterConfig(param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitializeClusterConfig", reflect.TypeOf((*MockClusterRegistrationOperator)(nil).InitializeClusterConfig), param)
}

// Remove mocks base method.
func (m *MockClusterRegistrationOperator) Remove() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove")
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockClusterRegistrationOperatorMockRecorder) Remove() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockClusterRegistrationOperator)(nil).Remove))
}

// Save mocks base method.
func (m *MockClusterRegistrationOperator) Save() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save")
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockClusterRegistrationOperatorMockRecorder) Save() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockClusterRegistrationOperator)(nil).Save))
}
