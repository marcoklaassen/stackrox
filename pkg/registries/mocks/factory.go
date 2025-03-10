// Code generated by MockGen. DO NOT EDIT.
// Source: factory.go
//
// Generated by this command:
//
//	mockgen -package mocks -destination mocks/factory.go -source factory.go
//
// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	storage "github.com/stackrox/rox/generated/storage"
	types "github.com/stackrox/rox/pkg/registries/types"
	gomock "go.uber.org/mock/gomock"
)

// MockFactory is a mock of Factory interface.
type MockFactory struct {
	ctrl     *gomock.Controller
	recorder *MockFactoryMockRecorder
}

// MockFactoryMockRecorder is the mock recorder for MockFactory.
type MockFactoryMockRecorder struct {
	mock *MockFactory
}

// NewMockFactory creates a new mock instance.
func NewMockFactory(ctrl *gomock.Controller) *MockFactory {
	mock := &MockFactory{ctrl: ctrl}
	mock.recorder = &MockFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFactory) EXPECT() *MockFactoryMockRecorder {
	return m.recorder
}

// CreateRegistry mocks base method.
func (m *MockFactory) CreateRegistry(source *storage.ImageIntegration) (types.ImageRegistry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRegistry", source)
	ret0, _ := ret[0].(types.ImageRegistry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRegistry indicates an expected call of CreateRegistry.
func (mr *MockFactoryMockRecorder) CreateRegistry(source any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRegistry", reflect.TypeOf((*MockFactory)(nil).CreateRegistry), source)
}
