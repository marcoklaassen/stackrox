// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/stackrox/rox/central/logimbue/handler (interfaces: Compressor)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCompressor is a mock of Compressor interface
type MockCompressor struct {
	ctrl     *gomock.Controller
	recorder *MockCompressorMockRecorder
}

// MockCompressorMockRecorder is the mock recorder for MockCompressor
type MockCompressorMockRecorder struct {
	mock *MockCompressor
}

// NewMockCompressor creates a new mock instance
func NewMockCompressor(ctrl *gomock.Controller) *MockCompressor {
	mock := &MockCompressor{ctrl: ctrl}
	mock.recorder = &MockCompressorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCompressor) EXPECT() *MockCompressorMockRecorder {
	return m.recorder
}

// Bytes mocks base method
func (m *MockCompressor) Bytes() []byte {
	ret := m.ctrl.Call(m, "Bytes")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Bytes indicates an expected call of Bytes
func (mr *MockCompressorMockRecorder) Bytes() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bytes", reflect.TypeOf((*MockCompressor)(nil).Bytes))
}

// Close mocks base method
func (m *MockCompressor) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockCompressorMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCompressor)(nil).Close))
}

// Write mocks base method
func (m *MockCompressor) Write(arg0 []byte) (int, error) {
	ret := m.ctrl.Call(m, "Write", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write
func (mr *MockCompressorMockRecorder) Write(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockCompressor)(nil).Write), arg0)
}
