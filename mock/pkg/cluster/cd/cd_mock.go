// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/cluster/cd/cd.go

// Package mock_cd is a generated GoMock package.
package mock_cd

import (
	context "context"
	cd "g.hz.netease.com/horizon/pkg/cluster/cd"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCD is a mock of CD interface
type MockCD struct {
	ctrl     *gomock.Controller
	recorder *MockCDMockRecorder
}

// MockCDMockRecorder is the mock recorder for MockCD
type MockCDMockRecorder struct {
	mock *MockCD
}

// NewMockCD creates a new mock instance
func NewMockCD(ctrl *gomock.Controller) *MockCD {
	mock := &MockCD{ctrl: ctrl}
	mock.recorder = &MockCDMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCD) EXPECT() *MockCDMockRecorder {
	return m.recorder
}

// CreateCluster mocks base method
func (m *MockCD) CreateCluster(ctx context.Context, params *cd.CreateClusterParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCluster", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCluster indicates an expected call of CreateCluster
func (mr *MockCDMockRecorder) CreateCluster(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCluster", reflect.TypeOf((*MockCD)(nil).CreateCluster), ctx, params)
}

// GetClusterState mocks base method
func (m *MockCD) GetClusterState(ctx context.Context, params *cd.GetClusterStateParams) (*cd.ClusterState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClusterState", ctx, params)
	ret0, _ := ret[0].(*cd.ClusterState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClusterState indicates an expected call of GetClusterState
func (mr *MockCDMockRecorder) GetClusterState(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusterState", reflect.TypeOf((*MockCD)(nil).GetClusterState), ctx, params)
}