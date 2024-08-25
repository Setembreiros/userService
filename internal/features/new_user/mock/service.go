// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_newuser is a generated GoMock package.
package mock_newuser

import (
	reflect "reflect"
	new_user "userservice/internal/features/new_user"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddNewUser mocks base method.
func (m *MockRepository) AddNewUser(data *new_user.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewUser", data)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewUser indicates an expected call of AddNewUser.
func (mr *MockRepositoryMockRecorder) AddNewUser(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewUser", reflect.TypeOf((*MockRepository)(nil).AddNewUser), data)
}