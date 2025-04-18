// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	core "autopilot/backends/internal/core"
	context "context"

	mock "github.com/stretchr/testify/mock"

	sqlx "github.com/jmoiron/sqlx"

	time "time"
)

// MockDBer is an autogenerated mock type for the DBer type
type MockDBer struct {
	mock.Mock
}

type MockDBer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDBer) EXPECT() *MockDBer_Expecter {
	return &MockDBer_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with no fields
func (_m *MockDBer) Close() {
	_m.Called()
}

// MockDBer_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type MockDBer_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *MockDBer_Expecter) Close() *MockDBer_Close_Call {
	return &MockDBer_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *MockDBer_Close_Call) Run(run func()) *MockDBer_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDBer_Close_Call) Return() *MockDBer_Close_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockDBer_Close_Call) RunAndReturn(run func()) *MockDBer_Close_Call {
	_c.Run(run)
	return _c
}

// GenMigration provides a mock function with given fields: name
func (_m *MockDBer) GenMigration(name string) error {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for GenMigration")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDBer_GenMigration_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenMigration'
type MockDBer_GenMigration_Call struct {
	*mock.Call
}

// GenMigration is a helper method to define mock.On call
//   - name string
func (_e *MockDBer_Expecter) GenMigration(name interface{}) *MockDBer_GenMigration_Call {
	return &MockDBer_GenMigration_Call{Call: _e.mock.On("GenMigration", name)}
}

func (_c *MockDBer_GenMigration_Call) Run(run func(name string)) *MockDBer_GenMigration_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockDBer_GenMigration_Call) Return(_a0 error) *MockDBer_GenMigration_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDBer_GenMigration_Call) RunAndReturn(run func(string) error) *MockDBer_GenMigration_Call {
	_c.Call.Return(run)
	return _c
}

// HealthCheck provides a mock function with given fields: ctx
func (_m *MockDBer) HealthCheck(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for HealthCheck")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDBer_HealthCheck_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HealthCheck'
type MockDBer_HealthCheck_Call struct {
	*mock.Call
}

// HealthCheck is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockDBer_Expecter) HealthCheck(ctx interface{}) *MockDBer_HealthCheck_Call {
	return &MockDBer_HealthCheck_Call{Call: _e.mock.On("HealthCheck", ctx)}
}

func (_c *MockDBer_HealthCheck_Call) Run(run func(ctx context.Context)) *MockDBer_HealthCheck_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockDBer_HealthCheck_Call) Return(_a0 error) *MockDBer_HealthCheck_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDBer_HealthCheck_Call) RunAndReturn(run func(context.Context) error) *MockDBer_HealthCheck_Call {
	_c.Call.Return(run)
	return _c
}

// Identifier provides a mock function with no fields
func (_m *MockDBer) Identifier() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Identifier")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockDBer_Identifier_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Identifier'
type MockDBer_Identifier_Call struct {
	*mock.Call
}

// Identifier is a helper method to define mock.On call
func (_e *MockDBer_Expecter) Identifier() *MockDBer_Identifier_Call {
	return &MockDBer_Identifier_Call{Call: _e.mock.On("Identifier")}
}

func (_c *MockDBer_Identifier_Call) Run(run func()) *MockDBer_Identifier_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDBer_Identifier_Call) Return(_a0 string) *MockDBer_Identifier_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDBer_Identifier_Call) RunAndReturn(run func() string) *MockDBer_Identifier_Call {
	_c.Call.Return(run)
	return _c
}

// Migrate provides a mock function with given fields: ctx
func (_m *MockDBer) Migrate(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Migrate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDBer_Migrate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Migrate'
type MockDBer_Migrate_Call struct {
	*mock.Call
}

// Migrate is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockDBer_Expecter) Migrate(ctx interface{}) *MockDBer_Migrate_Call {
	return &MockDBer_Migrate_Call{Call: _e.mock.On("Migrate", ctx)}
}

func (_c *MockDBer_Migrate_Call) Run(run func(ctx context.Context)) *MockDBer_Migrate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockDBer_Migrate_Call) Return(_a0 error) *MockDBer_Migrate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDBer_Migrate_Call) RunAndReturn(run func(context.Context) error) *MockDBer_Migrate_Call {
	_c.Call.Return(run)
	return _c
}

// MigrateStatus provides a mock function with given fields: ctx
func (_m *MockDBer) MigrateStatus(ctx context.Context) (int64, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for MigrateStatus")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (int64, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) int64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDBer_MigrateStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MigrateStatus'
type MockDBer_MigrateStatus_Call struct {
	*mock.Call
}

// MigrateStatus is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockDBer_Expecter) MigrateStatus(ctx interface{}) *MockDBer_MigrateStatus_Call {
	return &MockDBer_MigrateStatus_Call{Call: _e.mock.On("MigrateStatus", ctx)}
}

func (_c *MockDBer_MigrateStatus_Call) Run(run func(ctx context.Context)) *MockDBer_MigrateStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockDBer_MigrateStatus_Call) Return(_a0 int64, _a1 error) *MockDBer_MigrateStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDBer_MigrateStatus_Call) RunAndReturn(run func(context.Context) (int64, error)) *MockDBer_MigrateStatus_Call {
	_c.Call.Return(run)
	return _c
}

// Name provides a mock function with no fields
func (_m *MockDBer) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockDBer_Name_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Name'
type MockDBer_Name_Call struct {
	*mock.Call
}

// Name is a helper method to define mock.On call
func (_e *MockDBer_Expecter) Name() *MockDBer_Name_Call {
	return &MockDBer_Name_Call{Call: _e.mock.On("Name")}
}

func (_c *MockDBer_Name_Call) Run(run func()) *MockDBer_Name_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDBer_Name_Call) Return(_a0 string) *MockDBer_Name_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDBer_Name_Call) RunAndReturn(run func() string) *MockDBer_Name_Call {
	_c.Call.Return(run)
	return _c
}

// Options provides a mock function with no fields
func (_m *MockDBer) Options() core.DBOptions {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Options")
	}

	var r0 core.DBOptions
	if rf, ok := ret.Get(0).(func() core.DBOptions); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(core.DBOptions)
	}

	return r0
}

// MockDBer_Options_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Options'
type MockDBer_Options_Call struct {
	*mock.Call
}

// Options is a helper method to define mock.On call
func (_e *MockDBer_Expecter) Options() *MockDBer_Options_Call {
	return &MockDBer_Options_Call{Call: _e.mock.On("Options")}
}

func (_c *MockDBer_Options_Call) Run(run func()) *MockDBer_Options_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDBer_Options_Call) Return(_a0 core.DBOptions) *MockDBer_Options_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDBer_Options_Call) RunAndReturn(run func() core.DBOptions) *MockDBer_Options_Call {
	_c.Call.Return(run)
	return _c
}

// Reader provides a mock function with no fields
func (_m *MockDBer) Reader() *sqlx.DB {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Reader")
	}

	var r0 *sqlx.DB
	if rf, ok := ret.Get(0).(func() *sqlx.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.DB)
		}
	}

	return r0
}

// MockDBer_Reader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Reader'
type MockDBer_Reader_Call struct {
	*mock.Call
}

// Reader is a helper method to define mock.On call
func (_e *MockDBer_Expecter) Reader() *MockDBer_Reader_Call {
	return &MockDBer_Reader_Call{Call: _e.mock.On("Reader")}
}

func (_c *MockDBer_Reader_Call) Run(run func()) *MockDBer_Reader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDBer_Reader_Call) Return(_a0 *sqlx.DB) *MockDBer_Reader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDBer_Reader_Call) RunAndReturn(run func() *sqlx.DB) *MockDBer_Reader_Call {
	_c.Call.Return(run)
	return _c
}

// Seed provides a mock function with given fields: ctx
func (_m *MockDBer) Seed(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Seed")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDBer_Seed_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Seed'
type MockDBer_Seed_Call struct {
	*mock.Call
}

// Seed is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockDBer_Expecter) Seed(ctx interface{}) *MockDBer_Seed_Call {
	return &MockDBer_Seed_Call{Call: _e.mock.On("Seed", ctx)}
}

func (_c *MockDBer_Seed_Call) Run(run func(ctx context.Context)) *MockDBer_Seed_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockDBer_Seed_Call) Return(_a0 error) *MockDBer_Seed_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDBer_Seed_Call) RunAndReturn(run func(context.Context) error) *MockDBer_Seed_Call {
	_c.Call.Return(run)
	return _c
}

// WithTx provides a mock function with given fields: ctx, fn
func (_m *MockDBer) WithTx(ctx context.Context, fn func(context.Context, *sqlx.Tx) error) error {
	ret := _m.Called(ctx, fn)

	if len(ret) == 0 {
		panic("no return value specified for WithTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context, *sqlx.Tx) error) error); ok {
		r0 = rf(ctx, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDBer_WithTx_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTx'
type MockDBer_WithTx_Call struct {
	*mock.Call
}

// WithTx is a helper method to define mock.On call
//   - ctx context.Context
//   - fn func(context.Context , *sqlx.Tx) error
func (_e *MockDBer_Expecter) WithTx(ctx interface{}, fn interface{}) *MockDBer_WithTx_Call {
	return &MockDBer_WithTx_Call{Call: _e.mock.On("WithTx", ctx, fn)}
}

func (_c *MockDBer_WithTx_Call) Run(run func(ctx context.Context, fn func(context.Context, *sqlx.Tx) error)) *MockDBer_WithTx_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(context.Context, *sqlx.Tx) error))
	})
	return _c
}

func (_c *MockDBer_WithTx_Call) Return(_a0 error) *MockDBer_WithTx_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDBer_WithTx_Call) RunAndReturn(run func(context.Context, func(context.Context, *sqlx.Tx) error) error) *MockDBer_WithTx_Call {
	_c.Call.Return(run)
	return _c
}

// WithTxTimeout provides a mock function with given fields: ctx, timeout, fn
func (_m *MockDBer) WithTxTimeout(ctx context.Context, timeout time.Duration, fn func(context.Context, *sqlx.Tx) error) error {
	ret := _m.Called(ctx, timeout, fn)

	if len(ret) == 0 {
		panic("no return value specified for WithTxTimeout")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Duration, func(context.Context, *sqlx.Tx) error) error); ok {
		r0 = rf(ctx, timeout, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDBer_WithTxTimeout_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTxTimeout'
type MockDBer_WithTxTimeout_Call struct {
	*mock.Call
}

// WithTxTimeout is a helper method to define mock.On call
//   - ctx context.Context
//   - timeout time.Duration
//   - fn func(context.Context , *sqlx.Tx) error
func (_e *MockDBer_Expecter) WithTxTimeout(ctx interface{}, timeout interface{}, fn interface{}) *MockDBer_WithTxTimeout_Call {
	return &MockDBer_WithTxTimeout_Call{Call: _e.mock.On("WithTxTimeout", ctx, timeout, fn)}
}

func (_c *MockDBer_WithTxTimeout_Call) Run(run func(ctx context.Context, timeout time.Duration, fn func(context.Context, *sqlx.Tx) error)) *MockDBer_WithTxTimeout_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(time.Duration), args[2].(func(context.Context, *sqlx.Tx) error))
	})
	return _c
}

func (_c *MockDBer_WithTxTimeout_Call) Return(_a0 error) *MockDBer_WithTxTimeout_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDBer_WithTxTimeout_Call) RunAndReturn(run func(context.Context, time.Duration, func(context.Context, *sqlx.Tx) error) error) *MockDBer_WithTxTimeout_Call {
	_c.Call.Return(run)
	return _c
}

// Writer provides a mock function with no fields
func (_m *MockDBer) Writer() *sqlx.DB {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Writer")
	}

	var r0 *sqlx.DB
	if rf, ok := ret.Get(0).(func() *sqlx.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.DB)
		}
	}

	return r0
}

// MockDBer_Writer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Writer'
type MockDBer_Writer_Call struct {
	*mock.Call
}

// Writer is a helper method to define mock.On call
func (_e *MockDBer_Expecter) Writer() *MockDBer_Writer_Call {
	return &MockDBer_Writer_Call{Call: _e.mock.On("Writer")}
}

func (_c *MockDBer_Writer_Call) Run(run func()) *MockDBer_Writer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDBer_Writer_Call) Return(_a0 *sqlx.DB) *MockDBer_Writer_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDBer_Writer_Call) RunAndReturn(run func() *sqlx.DB) *MockDBer_Writer_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDBer creates a new instance of MockDBer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDBer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDBer {
	mock := &MockDBer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
