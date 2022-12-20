// Code generated by mockery v2.13.0-beta.1. DO NOT EDIT.

package mocks

import (
	models "cozy-doc-api/models"

	mock "github.com/stretchr/testify/mock"
)

// DocsService is an autogenerated mock type for the DocsService type
type DocsService struct {
	mock.Mock
}

// InsertDocs provides a mock function with given fields: req
func (_m *DocsService) InsertDocs(req *models.DocumentRequest) error {
	ret := _m.Called(req)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.DocumentRequest) error); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewDocsServiceT interface {
	mock.TestingT
	Cleanup(func())
}

// NewDocsService creates a new instance of DocsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDocsService(t NewDocsServiceT) *DocsService {
	mock := &DocsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
