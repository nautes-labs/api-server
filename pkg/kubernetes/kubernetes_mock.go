// Code generated by MockGen. DO NOT EDIT.
// Source: /root/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.14.5/pkg/client/interfaces.go

// Package kubernetes is a generated GoMock package.
package kubernetes

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	meta "k8s.io/apimachinery/pkg/api/meta"
	runtime "k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	resourcev1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

)

// MockPatch is a mock of Patch interface.
type MockPatch struct {
	ctrl     *gomock.Controller
	recorder *MockPatchMockRecorder
}

// MockPatchMockRecorder is the mock recorder for MockPatch.
type MockPatchMockRecorder struct {
	mock *MockPatch
}

// NewMockPatch creates a new mock instance.
func NewMockPatch(ctrl *gomock.Controller) *MockPatch {
	mock := &MockPatch{ctrl: ctrl}
	mock.recorder = &MockPatchMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPatch) EXPECT() *MockPatchMockRecorder {
	return m.recorder
}

// Data mocks base method.
func (m *MockPatch) Data(obj client.Object) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Data", obj)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Data indicates an expected call of Data.
func (mr *MockPatchMockRecorder) Data(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Data", reflect.TypeOf((*MockPatch)(nil).Data), obj)
}

// Type mocks base method.
func (m *MockPatch) Type() types.PatchType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(types.PatchType)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockPatchMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockPatch)(nil).Type))
}

// MockReader is a mock of Reader interface.
type MockReader struct {
	ctrl     *gomock.Controller
	recorder *MockReaderMockRecorder
}

// MockReaderMockRecorder is the mock recorder for MockReader.
type MockReaderMockRecorder struct {
	mock *MockReader
}

// NewMockReader creates a new mock instance.
func NewMockReader(ctrl *gomock.Controller) *MockReader {
	mock := &MockReader{ctrl: ctrl}
	mock.recorder = &MockReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReader) EXPECT() *MockReaderMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockReader) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, key, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockReaderMockRecorder) Get(ctx, key, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, key, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockReader)(nil).Get), varargs...)
}

// List mocks base method.
func (m *MockReader) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, list}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// List indicates an expected call of List.
func (mr *MockReaderMockRecorder) List(ctx, list interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, list}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockReader)(nil).List), varargs...)
}

// MockWriter is a mock of Writer interface.
type MockWriter struct {
	ctrl     *gomock.Controller
	recorder *MockWriterMockRecorder
}

// MockWriterMockRecorder is the mock recorder for MockWriter.
type MockWriterMockRecorder struct {
	mock *MockWriter
}

// NewMockWriter creates a new mock instance.
func NewMockWriter(ctrl *gomock.Controller) *MockWriter {
	mock := &MockWriter{ctrl: ctrl}
	mock.recorder = &MockWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWriter) EXPECT() *MockWriterMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockWriter) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockWriterMockRecorder) Create(ctx, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockWriter)(nil).Create), varargs...)
}

// Delete mocks base method.
func (m *MockWriter) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockWriterMockRecorder) Delete(ctx, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockWriter)(nil).Delete), varargs...)
}

// DeleteAllOf mocks base method.
func (m *MockWriter) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteAllOf", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAllOf indicates an expected call of DeleteAllOf.
func (mr *MockWriterMockRecorder) DeleteAllOf(ctx, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllOf", reflect.TypeOf((*MockWriter)(nil).DeleteAllOf), varargs...)
}

// Patch mocks base method.
func (m *MockWriter) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj, patch}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Patch", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Patch indicates an expected call of Patch.
func (mr *MockWriterMockRecorder) Patch(ctx, obj, patch interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj, patch}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockWriter)(nil).Patch), varargs...)
}

// Update mocks base method.
func (m *MockWriter) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Update", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockWriterMockRecorder) Update(ctx, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockWriter)(nil).Update), varargs...)
}

// MockStatusClient is a mock of StatusClient interface.
type MockStatusClient struct {
	ctrl     *gomock.Controller
	recorder *MockStatusClientMockRecorder
}

// MockStatusClientMockRecorder is the mock recorder for MockStatusClient.
type MockStatusClientMockRecorder struct {
	mock *MockStatusClient
}

// NewMockStatusClient creates a new mock instance.
func NewMockStatusClient(ctrl *gomock.Controller) *MockStatusClient {
	mock := &MockStatusClient{ctrl: ctrl}
	mock.recorder = &MockStatusClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStatusClient) EXPECT() *MockStatusClientMockRecorder {
	return m.recorder
}

// Status mocks base method.
func (m *MockStatusClient) Status() client.SubResourceWriter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(client.SubResourceWriter)
	return ret0
}

// Status indicates an expected call of Status.
func (mr *MockStatusClientMockRecorder) Status() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockStatusClient)(nil).Status))
}

// MockSubResourceClientConstructor is a mock of SubResourceClientConstructor interface.
type MockSubResourceClientConstructor struct {
	ctrl     *gomock.Controller
	recorder *MockSubResourceClientConstructorMockRecorder
}

// MockSubResourceClientConstructorMockRecorder is the mock recorder for MockSubResourceClientConstructor.
type MockSubResourceClientConstructorMockRecorder struct {
	mock *MockSubResourceClientConstructor
}

// NewMockSubResourceClientConstructor creates a new mock instance.
func NewMockSubResourceClientConstructor(ctrl *gomock.Controller) *MockSubResourceClientConstructor {
	mock := &MockSubResourceClientConstructor{ctrl: ctrl}
	mock.recorder = &MockSubResourceClientConstructorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubResourceClientConstructor) EXPECT() *MockSubResourceClientConstructorMockRecorder {
	return m.recorder
}

// SubResource mocks base method.
func (m *MockSubResourceClientConstructor) SubResource(subResource string) client.SubResourceClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubResource", subResource)
	ret0, _ := ret[0].(client.SubResourceClient)
	return ret0
}

// SubResource indicates an expected call of SubResource.
func (mr *MockSubResourceClientConstructorMockRecorder) SubResource(subResource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubResource", reflect.TypeOf((*MockSubResourceClientConstructor)(nil).SubResource), subResource)
}

// MockSubResourceReader is a mock of SubResourceReader interface.
type MockSubResourceReader struct {
	ctrl     *gomock.Controller
	recorder *MockSubResourceReaderMockRecorder
}

// MockSubResourceReaderMockRecorder is the mock recorder for MockSubResourceReader.
type MockSubResourceReaderMockRecorder struct {
	mock *MockSubResourceReader
}

// NewMockSubResourceReader creates a new mock instance.
func NewMockSubResourceReader(ctrl *gomock.Controller) *MockSubResourceReader {
	mock := &MockSubResourceReader{ctrl: ctrl}
	mock.recorder = &MockSubResourceReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubResourceReader) EXPECT() *MockSubResourceReaderMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockSubResourceReader) Get(ctx context.Context, obj, subResource client.Object, opts ...client.SubResourceGetOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj, subResource}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockSubResourceReaderMockRecorder) Get(ctx, obj, subResource interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj, subResource}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSubResourceReader)(nil).Get), varargs...)
}

// MockSubResourceWriter is a mock of SubResourceWriter interface.
type MockSubResourceWriter struct {
	ctrl     *gomock.Controller
	recorder *MockSubResourceWriterMockRecorder
}

// MockSubResourceWriterMockRecorder is the mock recorder for MockSubResourceWriter.
type MockSubResourceWriterMockRecorder struct {
	mock *MockSubResourceWriter
}

// NewMockSubResourceWriter creates a new mock instance.
func NewMockSubResourceWriter(ctrl *gomock.Controller) *MockSubResourceWriter {
	mock := &MockSubResourceWriter{ctrl: ctrl}
	mock.recorder = &MockSubResourceWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubResourceWriter) EXPECT() *MockSubResourceWriterMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSubResourceWriter) Create(ctx context.Context, obj, subResource client.Object, opts ...client.SubResourceCreateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj, subResource}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockSubResourceWriterMockRecorder) Create(ctx, obj, subResource interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj, subResource}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSubResourceWriter)(nil).Create), varargs...)
}

// Patch mocks base method.
func (m *MockSubResourceWriter) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj, patch}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Patch", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Patch indicates an expected call of Patch.
func (mr *MockSubResourceWriterMockRecorder) Patch(ctx, obj, patch interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj, patch}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockSubResourceWriter)(nil).Patch), varargs...)
}

// Update mocks base method.
func (m *MockSubResourceWriter) Update(ctx context.Context, obj client.Object, opts ...client.SubResourceUpdateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Update", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockSubResourceWriterMockRecorder) Update(ctx, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSubResourceWriter)(nil).Update), varargs...)
}

// MockSubResourceClient is a mock of SubResourceClient interface.
type MockSubResourceClient struct {
	ctrl     *gomock.Controller
	recorder *MockSubResourceClientMockRecorder
}

// MockSubResourceClientMockRecorder is the mock recorder for MockSubResourceClient.
type MockSubResourceClientMockRecorder struct {
	mock *MockSubResourceClient
}

// NewMockSubResourceClient creates a new mock instance.
func NewMockSubResourceClient(ctrl *gomock.Controller) *MockSubResourceClient {
	mock := &MockSubResourceClient{ctrl: ctrl}
	mock.recorder = &MockSubResourceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubResourceClient) EXPECT() *MockSubResourceClientMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSubResourceClient) Create(ctx context.Context, obj, subResource client.Object, opts ...client.SubResourceCreateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj, subResource}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockSubResourceClientMockRecorder) Create(ctx, obj, subResource interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj, subResource}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSubResourceClient)(nil).Create), varargs...)
}

// Get mocks base method.
func (m *MockSubResourceClient) Get(ctx context.Context, obj, subResource client.Object, opts ...client.SubResourceGetOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj, subResource}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockSubResourceClientMockRecorder) Get(ctx, obj, subResource interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj, subResource}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSubResourceClient)(nil).Get), varargs...)
}

// Patch mocks base method.
func (m *MockSubResourceClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj, patch}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Patch", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Patch indicates an expected call of Patch.
func (mr *MockSubResourceClientMockRecorder) Patch(ctx, obj, patch interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj, patch}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockSubResourceClient)(nil).Patch), varargs...)
}

// Update mocks base method.
func (m *MockSubResourceClient) Update(ctx context.Context, obj client.Object, opts ...client.SubResourceUpdateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Update", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockSubResourceClientMockRecorder) Update(ctx, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSubResourceClient)(nil).Update), varargs...)
}

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockClientMockRecorder) Create(ctx, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockClient)(nil).Create), varargs...)
}

// Delete mocks base method.
func (m *MockClient) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockClientMockRecorder) Delete(ctx, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockClient)(nil).Delete), varargs...)
}

// DeleteAllOf mocks base method.
func (m *MockClient) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteAllOf", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAllOf indicates an expected call of DeleteAllOf.
func (mr *MockClientMockRecorder) DeleteAllOf(ctx, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllOf", reflect.TypeOf((*MockClient)(nil).DeleteAllOf), varargs...)
}

// Get mocks base method.
func (m *MockClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	m.ctrl.T.Helper()

	// test cases mock data for deployemnt runtime
	if key.Name == "deployment-test" {
		cluster, _ := obj.(*resourcev1alpha1.Cluster)
		cluster.Spec.Usage = resourcev1alpha1.CLUSTER_USAGE_WORKER
		cluster.Spec.WorkerType = resourcev1alpha1.ClusterWorkTypeDeployment
		obj = cluster
	}

	if key.Name == "pipeline-test" {
		cluster, _ := obj.(*resourcev1alpha1.Cluster)
		cluster.Spec.Usage = resourcev1alpha1.CLUSTER_USAGE_WORKER
		cluster.Spec.WorkerType = resourcev1alpha1.ClusterWorkTypePipeline
		obj = cluster
	}


	varargs := []interface{}{ctx, key, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockClientMockRecorder) Get(ctx, key, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, key, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockClient)(nil).Get), varargs...)
}

// List mocks base method.
func (m *MockClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	m.ctrl.T.Helper()
	codeRepoList, ok := list.(*resourcev1alpha1.CodeRepoList)
	if ok {
		codeRepoList.Items = append(codeRepoList.Items, resourcev1alpha1.CodeRepo{
			ObjectMeta: v1.ObjectMeta{
				Name: "repo-22",
			},
		})
	}


	varargs := []interface{}{ctx, list}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// List indicates an expected call of List.
func (mr *MockClientMockRecorder) List(ctx, list interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, list}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockClient)(nil).List), varargs...)
}

// Patch mocks base method.
func (m *MockClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj, patch}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Patch", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Patch indicates an expected call of Patch.
func (mr *MockClientMockRecorder) Patch(ctx, obj, patch interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj, patch}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockClient)(nil).Patch), varargs...)
}

// RESTMapper mocks base method.
func (m *MockClient) RESTMapper() meta.RESTMapper {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RESTMapper")
	ret0, _ := ret[0].(meta.RESTMapper)
	return ret0
}

// RESTMapper indicates an expected call of RESTMapper.
func (mr *MockClientMockRecorder) RESTMapper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RESTMapper", reflect.TypeOf((*MockClient)(nil).RESTMapper))
}

// Scheme mocks base method.
func (m *MockClient) Scheme() *runtime.Scheme {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scheme")
	ret0, _ := ret[0].(*runtime.Scheme)
	return ret0
}

// Scheme indicates an expected call of Scheme.
func (mr *MockClientMockRecorder) Scheme() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scheme", reflect.TypeOf((*MockClient)(nil).Scheme))
}

// Status mocks base method.
func (m *MockClient) Status() client.SubResourceWriter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(client.SubResourceWriter)
	return ret0
}

// Status indicates an expected call of Status.
func (mr *MockClientMockRecorder) Status() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockClient)(nil).Status))
}

// SubResource mocks base method.
func (m *MockClient) SubResource(subResource string) client.SubResourceClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubResource", subResource)
	ret0, _ := ret[0].(client.SubResourceClient)
	return ret0
}

// SubResource indicates an expected call of SubResource.
func (mr *MockClientMockRecorder) SubResource(subResource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubResource", reflect.TypeOf((*MockClient)(nil).SubResource), subResource)
}

// Update mocks base method.
func (m *MockClient) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Update", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockClientMockRecorder) Update(ctx, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockClient)(nil).Update), varargs...)
}

// MockWithWatch is a mock of WithWatch interface.
type MockWithWatch struct {
	ctrl     *gomock.Controller
	recorder *MockWithWatchMockRecorder
}

// MockWithWatchMockRecorder is the mock recorder for MockWithWatch.
type MockWithWatchMockRecorder struct {
	mock *MockWithWatch
}

// NewMockWithWatch creates a new mock instance.
func NewMockWithWatch(ctrl *gomock.Controller) *MockWithWatch {
	mock := &MockWithWatch{ctrl: ctrl}
	mock.recorder = &MockWithWatchMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWithWatch) EXPECT() *MockWithWatchMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockWithWatch) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockWithWatchMockRecorder) Create(ctx, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockWithWatch)(nil).Create), varargs...)
}

// Delete mocks base method.
func (m *MockWithWatch) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockWithWatchMockRecorder) Delete(ctx, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockWithWatch)(nil).Delete), varargs...)
}

// DeleteAllOf mocks base method.
func (m *MockWithWatch) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteAllOf", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAllOf indicates an expected call of DeleteAllOf.
func (mr *MockWithWatchMockRecorder) DeleteAllOf(ctx, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllOf", reflect.TypeOf((*MockWithWatch)(nil).DeleteAllOf), varargs...)
}

// Get mocks base method.
func (m *MockWithWatch) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, key, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockWithWatchMockRecorder) Get(ctx, key, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, key, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockWithWatch)(nil).Get), varargs...)
}

// List mocks base method.
func (m *MockWithWatch) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, list}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// List indicates an expected call of List.
func (mr *MockWithWatchMockRecorder) List(ctx, list interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, list}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockWithWatch)(nil).List), varargs...)
}

// Patch mocks base method.
func (m *MockWithWatch) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj, patch}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Patch", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Patch indicates an expected call of Patch.
func (mr *MockWithWatchMockRecorder) Patch(ctx, obj, patch interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj, patch}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockWithWatch)(nil).Patch), varargs...)
}

// RESTMapper mocks base method.
func (m *MockWithWatch) RESTMapper() meta.RESTMapper {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RESTMapper")
	ret0, _ := ret[0].(meta.RESTMapper)
	return ret0
}

// RESTMapper indicates an expected call of RESTMapper.
func (mr *MockWithWatchMockRecorder) RESTMapper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RESTMapper", reflect.TypeOf((*MockWithWatch)(nil).RESTMapper))
}

// Scheme mocks base method.
func (m *MockWithWatch) Scheme() *runtime.Scheme {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scheme")
	ret0, _ := ret[0].(*runtime.Scheme)
	return ret0
}

// Scheme indicates an expected call of Scheme.
func (mr *MockWithWatchMockRecorder) Scheme() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scheme", reflect.TypeOf((*MockWithWatch)(nil).Scheme))
}

// Status mocks base method.
func (m *MockWithWatch) Status() client.SubResourceWriter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(client.SubResourceWriter)
	return ret0
}

// Status indicates an expected call of Status.
func (mr *MockWithWatchMockRecorder) Status() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockWithWatch)(nil).Status))
}

// SubResource mocks base method.
func (m *MockWithWatch) SubResource(subResource string) client.SubResourceClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubResource", subResource)
	ret0, _ := ret[0].(client.SubResourceClient)
	return ret0
}

// SubResource indicates an expected call of SubResource.
func (mr *MockWithWatchMockRecorder) SubResource(subResource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubResource", reflect.TypeOf((*MockWithWatch)(nil).SubResource), subResource)
}

// Update mocks base method.
func (m *MockWithWatch) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Update", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockWithWatchMockRecorder) Update(ctx, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockWithWatch)(nil).Update), varargs...)
}

// Watch mocks base method.
func (m *MockWithWatch) Watch(ctx context.Context, obj client.ObjectList, opts ...client.ListOption) (watch.Interface, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Watch", varargs...)
	ret0, _ := ret[0].(watch.Interface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch.
func (mr *MockWithWatchMockRecorder) Watch(ctx, obj interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockWithWatch)(nil).Watch), varargs...)
}

// MockFieldIndexer is a mock of FieldIndexer interface.
type MockFieldIndexer struct {
	ctrl     *gomock.Controller
	recorder *MockFieldIndexerMockRecorder
}

// MockFieldIndexerMockRecorder is the mock recorder for MockFieldIndexer.
type MockFieldIndexerMockRecorder struct {
	mock *MockFieldIndexer
}

// NewMockFieldIndexer creates a new mock instance.
func NewMockFieldIndexer(ctrl *gomock.Controller) *MockFieldIndexer {
	mock := &MockFieldIndexer{ctrl: ctrl}
	mock.recorder = &MockFieldIndexerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFieldIndexer) EXPECT() *MockFieldIndexerMockRecorder {
	return m.recorder
}

// IndexField mocks base method.
func (m *MockFieldIndexer) IndexField(ctx context.Context, obj client.Object, field string, extractValue client.IndexerFunc) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IndexField", ctx, obj, field, extractValue)
	ret0, _ := ret[0].(error)
	return ret0
}

// IndexField indicates an expected call of IndexField.
func (mr *MockFieldIndexerMockRecorder) IndexField(ctx, obj, field, extractValue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndexField", reflect.TypeOf((*MockFieldIndexer)(nil).IndexField), ctx, obj, field, extractValue)
}
