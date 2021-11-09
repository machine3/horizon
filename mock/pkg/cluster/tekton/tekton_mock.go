// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/cluster/tekton/tekton.go

// Package mock_tekton is a generated GoMock package.
package mock_tekton

import (
	context "context"
	tekton "g.hz.netease.com/horizon/pkg/cluster/tekton"
	log "g.hz.netease.com/horizon/pkg/cluster/tekton/log"
	gomock "github.com/golang/mock/gomock"
	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	reflect "reflect"
)

// MockInterface is a mock of Interface interface
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// GetPipelineRunByID mocks base method
func (m *MockInterface) GetPipelineRunByID(ctx context.Context, cluster string, clusterID, pipelinerunID uint) (*v1beta1.PipelineRun, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPipelineRunByID", ctx, cluster, clusterID, pipelinerunID)
	ret0, _ := ret[0].(*v1beta1.PipelineRun)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPipelineRunByID indicates an expected call of GetPipelineRunByID
func (mr *MockInterfaceMockRecorder) GetPipelineRunByID(ctx, cluster, clusterID, pipelinerunID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPipelineRunByID", reflect.TypeOf((*MockInterface)(nil).GetPipelineRunByID), ctx, cluster, clusterID, pipelinerunID)
}

// CreatePipelineRun mocks base method
func (m *MockInterface) CreatePipelineRun(ctx context.Context, pr *tekton.PipelineRun) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePipelineRun", ctx, pr)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePipelineRun indicates an expected call of CreatePipelineRun
func (mr *MockInterfaceMockRecorder) CreatePipelineRun(ctx, pr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePipelineRun", reflect.TypeOf((*MockInterface)(nil).CreatePipelineRun), ctx, pr)
}

// StopPipelineRun mocks base method
func (m *MockInterface) StopPipelineRun(ctx context.Context, cluster string, clusterID, pipelinerunID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopPipelineRun", ctx, cluster, clusterID, pipelinerunID)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopPipelineRun indicates an expected call of StopPipelineRun
func (mr *MockInterfaceMockRecorder) StopPipelineRun(ctx, cluster, clusterID, pipelinerunID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopPipelineRun", reflect.TypeOf((*MockInterface)(nil).StopPipelineRun), ctx, cluster, clusterID, pipelinerunID)
}

// GetPipelineRunLogByID mocks base method
func (m *MockInterface) GetPipelineRunLogByID(ctx context.Context, cluster string, clusterID, pipelinerunID uint) (<-chan log.Log, <-chan error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPipelineRunLogByID", ctx, cluster, clusterID, pipelinerunID)
	ret0, _ := ret[0].(<-chan log.Log)
	ret1, _ := ret[1].(<-chan error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPipelineRunLogByID indicates an expected call of GetPipelineRunLogByID
func (mr *MockInterfaceMockRecorder) GetPipelineRunLogByID(ctx, cluster, clusterID, pipelinerunID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPipelineRunLogByID", reflect.TypeOf((*MockInterface)(nil).GetPipelineRunLogByID), ctx, cluster, clusterID, pipelinerunID)
}

// GetPipelineRunLog mocks base method
func (m *MockInterface) GetPipelineRunLog(ctx context.Context, pr *v1beta1.PipelineRun) (<-chan log.Log, <-chan error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPipelineRunLog", ctx, pr)
	ret0, _ := ret[0].(<-chan log.Log)
	ret1, _ := ret[1].(<-chan error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPipelineRunLog indicates an expected call of GetPipelineRunLog
func (mr *MockInterfaceMockRecorder) GetPipelineRunLog(ctx, pr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPipelineRunLog", reflect.TypeOf((*MockInterface)(nil).GetPipelineRunLog), ctx, pr)
}

// DeletePipelineRun mocks base method
func (m *MockInterface) DeletePipelineRun(ctx context.Context, pr *v1beta1.PipelineRun) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePipelineRun", ctx, pr)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePipelineRun indicates an expected call of DeletePipelineRun
func (mr *MockInterfaceMockRecorder) DeletePipelineRun(ctx, pr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePipelineRun", reflect.TypeOf((*MockInterface)(nil).DeletePipelineRun), ctx, pr)
}