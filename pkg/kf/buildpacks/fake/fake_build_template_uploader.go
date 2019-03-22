// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/GoogleCloudPlatform/kf/pkg/kf/buildpacks/fake (interfaces: BuildTemplateUploader)

// Package fake is a generated GoMock package.
package fake

import (
	buildpacks "github.com/GoogleCloudPlatform/kf/pkg/kf/buildpacks"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// FakeBuildTemplateUploader is a mock of BuildTemplateUploader interface
type FakeBuildTemplateUploader struct {
	ctrl     *gomock.Controller
	recorder *FakeBuildTemplateUploaderMockRecorder
}

// FakeBuildTemplateUploaderMockRecorder is the mock recorder for FakeBuildTemplateUploader
type FakeBuildTemplateUploaderMockRecorder struct {
	mock *FakeBuildTemplateUploader
}

// NewFakeBuildTemplateUploader creates a new mock instance
func NewFakeBuildTemplateUploader(ctrl *gomock.Controller) *FakeBuildTemplateUploader {
	mock := &FakeBuildTemplateUploader{ctrl: ctrl}
	mock.recorder = &FakeBuildTemplateUploaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *FakeBuildTemplateUploader) EXPECT() *FakeBuildTemplateUploaderMockRecorder {
	return m.recorder
}

// UploadBuildTemplate mocks base method
func (m *FakeBuildTemplateUploader) UploadBuildTemplate(arg0 string, arg1 ...buildpacks.UploadBuildTemplateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UploadBuildTemplate", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadBuildTemplate indicates an expected call of UploadBuildTemplate
func (mr *FakeBuildTemplateUploaderMockRecorder) UploadBuildTemplate(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadBuildTemplate", reflect.TypeOf((*FakeBuildTemplateUploader)(nil).UploadBuildTemplate), varargs...)
}
