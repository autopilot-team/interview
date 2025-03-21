// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	core "autopilot/backends/internal/core"
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "autopilot/backends/api/internal/identity/model"

	store "autopilot/backends/api/internal/identity/store"
)

// MockAuditLoger is an autogenerated mock type for the AuditLoger type
type MockAuditLoger struct {
	mock.Mock
}

type MockAuditLoger_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAuditLoger) EXPECT() *MockAuditLoger_Expecter {
	return &MockAuditLoger_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, log
func (_m *MockAuditLoger) Create(ctx context.Context, log *model.AuditLog) (*model.AuditLog, error) {
	ret := _m.Called(ctx, log)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *model.AuditLog
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AuditLog) (*model.AuditLog, error)); ok {
		return rf(ctx, log)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.AuditLog) *model.AuditLog); ok {
		r0 = rf(ctx, log)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AuditLog)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.AuditLog) error); ok {
		r1 = rf(ctx, log)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAuditLoger_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockAuditLoger_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - log *model.AuditLog
func (_e *MockAuditLoger_Expecter) Create(ctx interface{}, log interface{}) *MockAuditLoger_Create_Call {
	return &MockAuditLoger_Create_Call{Call: _e.mock.On("Create", ctx, log)}
}

func (_c *MockAuditLoger_Create_Call) Run(run func(ctx context.Context, log *model.AuditLog)) *MockAuditLoger_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.AuditLog))
	})
	return _c
}

func (_c *MockAuditLoger_Create_Call) Return(_a0 *model.AuditLog, _a1 error) *MockAuditLoger_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAuditLoger_Create_Call) RunAndReturn(run func(context.Context, *model.AuditLog) (*model.AuditLog, error)) *MockAuditLoger_Create_Call {
	_c.Call.Return(run)
	return _c
}

// WithQuerier provides a mock function with given fields: q
func (_m *MockAuditLoger) WithQuerier(q core.Querier) store.AuditLoger {
	ret := _m.Called(q)

	if len(ret) == 0 {
		panic("no return value specified for WithQuerier")
	}

	var r0 store.AuditLoger
	if rf, ok := ret.Get(0).(func(core.Querier) store.AuditLoger); ok {
		r0 = rf(q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.AuditLoger)
		}
	}

	return r0
}

// MockAuditLoger_WithQuerier_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithQuerier'
type MockAuditLoger_WithQuerier_Call struct {
	*mock.Call
}

// WithQuerier is a helper method to define mock.On call
//   - q core.Querier
func (_e *MockAuditLoger_Expecter) WithQuerier(q interface{}) *MockAuditLoger_WithQuerier_Call {
	return &MockAuditLoger_WithQuerier_Call{Call: _e.mock.On("WithQuerier", q)}
}

func (_c *MockAuditLoger_WithQuerier_Call) Run(run func(q core.Querier)) *MockAuditLoger_WithQuerier_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(core.Querier))
	})
	return _c
}

func (_c *MockAuditLoger_WithQuerier_Call) Return(_a0 store.AuditLoger) *MockAuditLoger_WithQuerier_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuditLoger_WithQuerier_Call) RunAndReturn(run func(core.Querier) store.AuditLoger) *MockAuditLoger_WithQuerier_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAuditLoger creates a new instance of MockAuditLoger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAuditLoger(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAuditLoger {
	mock := &MockAuditLoger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
