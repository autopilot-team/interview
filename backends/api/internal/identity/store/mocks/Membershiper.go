// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	core "autopilot/backends/internal/core"
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "autopilot/backends/api/internal/identity/model"

	store "autopilot/backends/api/internal/identity/store"
)

// MockMembershiper is an autogenerated mock type for the Membershiper type
type MockMembershiper struct {
	mock.Mock
}

type MockMembershiper_Expecter struct {
	mock *mock.Mock
}

func (_m *MockMembershiper) EXPECT() *MockMembershiper_Expecter {
	return &MockMembershiper_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, membership
func (_m *MockMembershiper) Create(ctx context.Context, membership *model.Membership) (*model.Membership, error) {
	ret := _m.Called(ctx, membership)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *model.Membership
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Membership) (*model.Membership, error)); ok {
		return rf(ctx, membership)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.Membership) *model.Membership); ok {
		r0 = rf(ctx, membership)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Membership)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.Membership) error); ok {
		r1 = rf(ctx, membership)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMembershiper_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockMembershiper_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - membership *model.Membership
func (_e *MockMembershiper_Expecter) Create(ctx interface{}, membership interface{}) *MockMembershiper_Create_Call {
	return &MockMembershiper_Create_Call{Call: _e.mock.On("Create", ctx, membership)}
}

func (_c *MockMembershiper_Create_Call) Run(run func(ctx context.Context, membership *model.Membership)) *MockMembershiper_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Membership))
	})
	return _c
}

func (_c *MockMembershiper_Create_Call) Return(_a0 *model.Membership, _a1 error) *MockMembershiper_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMembershiper_Create_Call) RunAndReturn(run func(context.Context, *model.Membership) (*model.Membership, error)) *MockMembershiper_Create_Call {
	_c.Call.Return(run)
	return _c
}

// GetByEntityID provides a mock function with given fields: ctx, entityID
func (_m *MockMembershiper) GetByEntityID(ctx context.Context, entityID string) ([]*model.Membership, error) {
	ret := _m.Called(ctx, entityID)

	if len(ret) == 0 {
		panic("no return value specified for GetByEntityID")
	}

	var r0 []*model.Membership
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*model.Membership, error)); ok {
		return rf(ctx, entityID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Membership); ok {
		r0 = rf(ctx, entityID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Membership)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, entityID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMembershiper_GetByEntityID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByEntityID'
type MockMembershiper_GetByEntityID_Call struct {
	*mock.Call
}

// GetByEntityID is a helper method to define mock.On call
//   - ctx context.Context
//   - entityID string
func (_e *MockMembershiper_Expecter) GetByEntityID(ctx interface{}, entityID interface{}) *MockMembershiper_GetByEntityID_Call {
	return &MockMembershiper_GetByEntityID_Call{Call: _e.mock.On("GetByEntityID", ctx, entityID)}
}

func (_c *MockMembershiper_GetByEntityID_Call) Run(run func(ctx context.Context, entityID string)) *MockMembershiper_GetByEntityID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockMembershiper_GetByEntityID_Call) Return(_a0 []*model.Membership, _a1 error) *MockMembershiper_GetByEntityID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMembershiper_GetByEntityID_Call) RunAndReturn(run func(context.Context, string) ([]*model.Membership, error)) *MockMembershiper_GetByEntityID_Call {
	_c.Call.Return(run)
	return _c
}

// GetByEntityIDWithInheritance provides a mock function with given fields: ctx, userID, entityID
func (_m *MockMembershiper) GetByEntityIDWithInheritance(ctx context.Context, userID string, entityID string) ([]*model.Membership, error) {
	ret := _m.Called(ctx, userID, entityID)

	if len(ret) == 0 {
		panic("no return value specified for GetByEntityIDWithInheritance")
	}

	var r0 []*model.Membership
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]*model.Membership, error)); ok {
		return rf(ctx, userID, entityID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []*model.Membership); ok {
		r0 = rf(ctx, userID, entityID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Membership)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, userID, entityID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMembershiper_GetByEntityIDWithInheritance_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByEntityIDWithInheritance'
type MockMembershiper_GetByEntityIDWithInheritance_Call struct {
	*mock.Call
}

// GetByEntityIDWithInheritance is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
//   - entityID string
func (_e *MockMembershiper_Expecter) GetByEntityIDWithInheritance(ctx interface{}, userID interface{}, entityID interface{}) *MockMembershiper_GetByEntityIDWithInheritance_Call {
	return &MockMembershiper_GetByEntityIDWithInheritance_Call{Call: _e.mock.On("GetByEntityIDWithInheritance", ctx, userID, entityID)}
}

func (_c *MockMembershiper_GetByEntityIDWithInheritance_Call) Run(run func(ctx context.Context, userID string, entityID string)) *MockMembershiper_GetByEntityIDWithInheritance_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockMembershiper_GetByEntityIDWithInheritance_Call) Return(_a0 []*model.Membership, _a1 error) *MockMembershiper_GetByEntityIDWithInheritance_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMembershiper_GetByEntityIDWithInheritance_Call) RunAndReturn(run func(context.Context, string, string) ([]*model.Membership, error)) *MockMembershiper_GetByEntityIDWithInheritance_Call {
	_c.Call.Return(run)
	return _c
}

// GetByUserID provides a mock function with given fields: ctx, userID
func (_m *MockMembershiper) GetByUserID(ctx context.Context, userID string) ([]*model.Membership, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetByUserID")
	}

	var r0 []*model.Membership
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*model.Membership, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Membership); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Membership)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMembershiper_GetByUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByUserID'
type MockMembershiper_GetByUserID_Call struct {
	*mock.Call
}

// GetByUserID is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
func (_e *MockMembershiper_Expecter) GetByUserID(ctx interface{}, userID interface{}) *MockMembershiper_GetByUserID_Call {
	return &MockMembershiper_GetByUserID_Call{Call: _e.mock.On("GetByUserID", ctx, userID)}
}

func (_c *MockMembershiper_GetByUserID_Call) Run(run func(ctx context.Context, userID string)) *MockMembershiper_GetByUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockMembershiper_GetByUserID_Call) Return(_a0 []*model.Membership, _a1 error) *MockMembershiper_GetByUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMembershiper_GetByUserID_Call) RunAndReturn(run func(context.Context, string) ([]*model.Membership, error)) *MockMembershiper_GetByUserID_Call {
	_c.Call.Return(run)
	return _c
}

// GetByUserIDWithInheritance provides a mock function with given fields: ctx, userID
func (_m *MockMembershiper) GetByUserIDWithInheritance(ctx context.Context, userID string) ([]*model.Membership, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetByUserIDWithInheritance")
	}

	var r0 []*model.Membership
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*model.Membership, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Membership); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Membership)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMembershiper_GetByUserIDWithInheritance_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByUserIDWithInheritance'
type MockMembershiper_GetByUserIDWithInheritance_Call struct {
	*mock.Call
}

// GetByUserIDWithInheritance is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
func (_e *MockMembershiper_Expecter) GetByUserIDWithInheritance(ctx interface{}, userID interface{}) *MockMembershiper_GetByUserIDWithInheritance_Call {
	return &MockMembershiper_GetByUserIDWithInheritance_Call{Call: _e.mock.On("GetByUserIDWithInheritance", ctx, userID)}
}

func (_c *MockMembershiper_GetByUserIDWithInheritance_Call) Run(run func(ctx context.Context, userID string)) *MockMembershiper_GetByUserIDWithInheritance_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockMembershiper_GetByUserIDWithInheritance_Call) Return(_a0 []*model.Membership, _a1 error) *MockMembershiper_GetByUserIDWithInheritance_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMembershiper_GetByUserIDWithInheritance_Call) RunAndReturn(run func(context.Context, string) ([]*model.Membership, error)) *MockMembershiper_GetByUserIDWithInheritance_Call {
	_c.Call.Return(run)
	return _c
}

// WithQuerier provides a mock function with given fields: q
func (_m *MockMembershiper) WithQuerier(q core.Querier) store.Membershiper {
	ret := _m.Called(q)

	if len(ret) == 0 {
		panic("no return value specified for WithQuerier")
	}

	var r0 store.Membershiper
	if rf, ok := ret.Get(0).(func(core.Querier) store.Membershiper); ok {
		r0 = rf(q)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.Membershiper)
		}
	}

	return r0
}

// MockMembershiper_WithQuerier_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithQuerier'
type MockMembershiper_WithQuerier_Call struct {
	*mock.Call
}

// WithQuerier is a helper method to define mock.On call
//   - q core.Querier
func (_e *MockMembershiper_Expecter) WithQuerier(q interface{}) *MockMembershiper_WithQuerier_Call {
	return &MockMembershiper_WithQuerier_Call{Call: _e.mock.On("WithQuerier", q)}
}

func (_c *MockMembershiper_WithQuerier_Call) Run(run func(q core.Querier)) *MockMembershiper_WithQuerier_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(core.Querier))
	})
	return _c
}

func (_c *MockMembershiper_WithQuerier_Call) Return(_a0 store.Membershiper) *MockMembershiper_WithQuerier_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMembershiper_WithQuerier_Call) RunAndReturn(run func(core.Querier) store.Membershiper) *MockMembershiper_WithQuerier_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockMembershiper creates a new instance of MockMembershiper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMembershiper(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMembershiper {
	mock := &MockMembershiper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
