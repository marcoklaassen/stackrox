// Code generated by MockGen. DO NOT EDIT.
// Source: set.go
//
// Generated by this command:
//
//	mockgen -package mocks -destination mocks/set.go -source set.go
//
// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	storage "github.com/stackrox/rox/generated/storage"
	types "github.com/stackrox/rox/pkg/scanners/types"
	gomock "go.uber.org/mock/gomock"
)

// MockSet is a mock of Set interface.
type MockSet struct {
	ctrl     *gomock.Controller
	recorder *MockSetMockRecorder
}

// MockSetMockRecorder is the mock recorder for MockSet.
type MockSetMockRecorder struct {
	mock *MockSet
}

// NewMockSet creates a new mock instance.
func NewMockSet(ctrl *gomock.Controller) *MockSet {
	mock := &MockSet{ctrl: ctrl}
	mock.recorder = &MockSetMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSet) EXPECT() *MockSetMockRecorder {
	return m.recorder
}

// Clear mocks base method.
func (m *MockSet) Clear() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Clear")
}

// Clear indicates an expected call of Clear.
func (mr *MockSetMockRecorder) Clear() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockSet)(nil).Clear))
}

// GetAll mocks base method.
func (m *MockSet) GetAll() []types.ImageScannerWithDataSource {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]types.ImageScannerWithDataSource)
	return ret0
}

// GetAll indicates an expected call of GetAll.
func (mr *MockSetMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockSet)(nil).GetAll))
}

// IsEmpty mocks base method.
func (m *MockSet) IsEmpty() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsEmpty")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsEmpty indicates an expected call of IsEmpty.
func (mr *MockSetMockRecorder) IsEmpty() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEmpty", reflect.TypeOf((*MockSet)(nil).IsEmpty))
}

// RemoveImageIntegration mocks base method.
func (m *MockSet) RemoveImageIntegration(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveImageIntegration", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveImageIntegration indicates an expected call of RemoveImageIntegration.
func (mr *MockSetMockRecorder) RemoveImageIntegration(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveImageIntegration", reflect.TypeOf((*MockSet)(nil).RemoveImageIntegration), id)
}

// UpdateImageIntegration mocks base method.
func (m *MockSet) UpdateImageIntegration(integration *storage.ImageIntegration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateImageIntegration", integration)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateImageIntegration indicates an expected call of UpdateImageIntegration.
func (mr *MockSetMockRecorder) UpdateImageIntegration(integration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateImageIntegration", reflect.TypeOf((*MockSet)(nil).UpdateImageIntegration), integration)
}
