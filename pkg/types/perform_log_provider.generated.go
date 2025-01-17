// Code generated by mockery v2.12.1. DO NOT EDIT.

package types

import (
	context "context"
	testing "testing"

	mock "github.com/stretchr/testify/mock"
)

// MockPerformLogProvider is an autogenerated mock type for the PerformLogProvider type
type MockPerformLogProvider struct {
	mock.Mock
}

// PerformLogs provides a mock function with given fields: _a0
func (_m *MockPerformLogProvider) PerformLogs(_a0 context.Context) ([]PerformLog, error) {
	ret := _m.Called(_a0)

	var r0 []PerformLog
	if rf, ok := ret.Get(0).(func(context.Context) []PerformLog); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]PerformLog)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StaleReportLogs provides a mock function with given fields: _a0
func (_m *MockPerformLogProvider) StaleReportLogs(_a0 context.Context) ([]StaleReportLog, error) {
	ret := _m.Called(_a0)

	var r0 []StaleReportLog
	if rf, ok := ret.Get(0).(func(context.Context) []StaleReportLog); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]StaleReportLog)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockPerformLogProvider creates a new instance of MockPerformLogProvider. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockPerformLogProvider(t testing.TB) *MockPerformLogProvider {
	mock := &MockPerformLogProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
