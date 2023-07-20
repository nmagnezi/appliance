// Code generated by MockGen. DO NOT EDIT.
// Source: release.go

// Package release is a generated GoMock package.
package release

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	config "github.com/openshift/appliance/pkg/asset/config"
)

// MockRelease is a mock of Release interface.
type MockRelease struct {
	ctrl     *gomock.Controller
	recorder *MockReleaseMockRecorder
}

// MockReleaseMockRecorder is the mock recorder for MockRelease.
type MockReleaseMockRecorder struct {
	mock *MockRelease
}

// NewMockRelease creates a new mock instance.
func NewMockRelease(ctrl *gomock.Controller) *MockRelease {
	mock := &MockRelease{ctrl: ctrl}
	mock.recorder = &MockReleaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRelease) EXPECT() *MockReleaseMockRecorder {
	return m.recorder
}

// ExtractFile mocks base method.
func (m *MockRelease) ExtractFile(image, filename string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExtractFile", image, filename)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExtractFile indicates an expected call of ExtractFile.
func (mr *MockReleaseMockRecorder) ExtractFile(image, filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExtractFile", reflect.TypeOf((*MockRelease)(nil).ExtractFile), image, filename)
}

// GetImageFromRelease mocks base method.
func (m *MockRelease) GetImageFromRelease(imageName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImageFromRelease", imageName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImageFromRelease indicates an expected call of GetImageFromRelease.
func (mr *MockReleaseMockRecorder) GetImageFromRelease(imageName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImageFromRelease", reflect.TypeOf((*MockRelease)(nil).GetImageFromRelease), imageName)
}

// MirrorBootstrapImages mocks base method.
func (m *MockRelease) MirrorBootstrapImages(envConfig *config.EnvConfig, applianceConfig *config.ApplianceConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MirrorBootstrapImages", envConfig, applianceConfig)
	ret0, _ := ret[0].(error)
	return ret0
}

// MirrorBootstrapImages indicates an expected call of MirrorBootstrapImages.
func (mr *MockReleaseMockRecorder) MirrorBootstrapImages(envConfig, applianceConfig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MirrorBootstrapImages", reflect.TypeOf((*MockRelease)(nil).MirrorBootstrapImages), envConfig, applianceConfig)
}

// MirrorReleaseImages mocks base method.
func (m *MockRelease) MirrorReleaseImages(envConfig *config.EnvConfig, applianceConfig *config.ApplianceConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MirrorReleaseImages", envConfig, applianceConfig)
	ret0, _ := ret[0].(error)
	return ret0
}

// MirrorReleaseImages indicates an expected call of MirrorReleaseImages.
func (mr *MockReleaseMockRecorder) MirrorReleaseImages(envConfig, applianceConfig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MirrorReleaseImages", reflect.TypeOf((*MockRelease)(nil).MirrorReleaseImages), envConfig, applianceConfig)
}
