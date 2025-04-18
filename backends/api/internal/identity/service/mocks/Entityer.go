// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	model "autopilot/backends/api/internal/identity/model"
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "autopilot/backends/internal/types"
)

// MockEntityer is an autogenerated mock type for the Entityer type
type MockEntityer struct {
	mock.Mock
}

type MockEntityer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEntityer) EXPECT() *MockEntityer_Expecter {
	return &MockEntityer_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, id
func (_m *MockEntityer) Get(ctx context.Context, id string) (*model.Entity, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *model.Entity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Entity, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Entity); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Entity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEntityer_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockEntityer_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *MockEntityer_Expecter) Get(ctx interface{}, id interface{}) *MockEntityer_Get_Call {
	return &MockEntityer_Get_Call{Call: _e.mock.On("Get", ctx, id)}
}

func (_c *MockEntityer_Get_Call) Run(run func(ctx context.Context, id string)) *MockEntityer_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockEntityer_Get_Call) Return(_a0 *model.Entity, _a1 error) *MockEntityer_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEntityer_Get_Call) RunAndReturn(run func(context.Context, string) (*model.Entity, error)) *MockEntityer_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *MockEntityer) GetByID(ctx context.Context, id string) (*model.Entity, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *model.Entity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Entity, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Entity); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Entity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEntityer_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type MockEntityer_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *MockEntityer_Expecter) GetByID(ctx interface{}, id interface{}) *MockEntityer_GetByID_Call {
	return &MockEntityer_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *MockEntityer_GetByID_Call) Run(run func(ctx context.Context, id string)) *MockEntityer_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockEntityer_GetByID_Call) Return(_a0 *model.Entity, _a1 error) *MockEntityer_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEntityer_GetByID_Call) RunAndReturn(run func(context.Context, string) (*model.Entity, error)) *MockEntityer_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetBySlug provides a mock function with given fields: ctx, mode, slug
func (_m *MockEntityer) GetBySlug(ctx context.Context, mode types.OperationMode, slug string) (*model.Entity, error) {
	ret := _m.Called(ctx, mode, slug)

	if len(ret) == 0 {
		panic("no return value specified for GetBySlug")
	}

	var r0 *model.Entity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.OperationMode, string) (*model.Entity, error)); ok {
		return rf(ctx, mode, slug)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.OperationMode, string) *model.Entity); ok {
		r0 = rf(ctx, mode, slug)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Entity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.OperationMode, string) error); ok {
		r1 = rf(ctx, mode, slug)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEntityer_GetBySlug_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBySlug'
type MockEntityer_GetBySlug_Call struct {
	*mock.Call
}

// GetBySlug is a helper method to define mock.On call
//   - ctx context.Context
//   - mode types.OperationMode
//   - slug string
func (_e *MockEntityer_Expecter) GetBySlug(ctx interface{}, mode interface{}, slug interface{}) *MockEntityer_GetBySlug_Call {
	return &MockEntityer_GetBySlug_Call{Call: _e.mock.On("GetBySlug", ctx, mode, slug)}
}

func (_c *MockEntityer_GetBySlug_Call) Run(run func(ctx context.Context, mode types.OperationMode, slug string)) *MockEntityer_GetBySlug_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.OperationMode), args[2].(string))
	})
	return _c
}

func (_c *MockEntityer_GetBySlug_Call) Return(_a0 *model.Entity, _a1 error) *MockEntityer_GetBySlug_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEntityer_GetBySlug_Call) RunAndReturn(run func(context.Context, types.OperationMode, string) (*model.Entity, error)) *MockEntityer_GetBySlug_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockEntityer creates a new instance of MockEntityer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEntityer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEntityer {
	mock := &MockEntityer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
