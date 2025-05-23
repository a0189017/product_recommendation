// Code generated by mockery v2.53.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// DBInterface is an autogenerated mock type for the DBInterface type
type DBInterface struct {
	mock.Mock
}

type DBInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *DBInterface) EXPECT() *DBInterface_Expecter {
	return &DBInterface_Expecter{mock: &_m.Mock}
}

// GetDB provides a mock function with no fields
func (_m *DBInterface) GetDB() *gorm.DB {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetDB")
	}

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func() *gorm.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// DBInterface_GetDB_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDB'
type DBInterface_GetDB_Call struct {
	*mock.Call
}

// GetDB is a helper method to define mock.On call
func (_e *DBInterface_Expecter) GetDB() *DBInterface_GetDB_Call {
	return &DBInterface_GetDB_Call{Call: _e.mock.On("GetDB")}
}

func (_c *DBInterface_GetDB_Call) Run(run func()) *DBInterface_GetDB_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DBInterface_GetDB_Call) Return(_a0 *gorm.DB) *DBInterface_GetDB_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DBInterface_GetDB_Call) RunAndReturn(run func() *gorm.DB) *DBInterface_GetDB_Call {
	_c.Call.Return(run)
	return _c
}

// LoadSchemaFields provides a mock function with no fields
func (_m *DBInterface) LoadSchemaFields() {
	_m.Called()
}

// DBInterface_LoadSchemaFields_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LoadSchemaFields'
type DBInterface_LoadSchemaFields_Call struct {
	*mock.Call
}

// LoadSchemaFields is a helper method to define mock.On call
func (_e *DBInterface_Expecter) LoadSchemaFields() *DBInterface_LoadSchemaFields_Call {
	return &DBInterface_LoadSchemaFields_Call{Call: _e.mock.On("LoadSchemaFields")}
}

func (_c *DBInterface_LoadSchemaFields_Call) Run(run func()) *DBInterface_LoadSchemaFields_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DBInterface_LoadSchemaFields_Call) Return() *DBInterface_LoadSchemaFields_Call {
	_c.Call.Return()
	return _c
}

func (_c *DBInterface_LoadSchemaFields_Call) RunAndReturn(run func()) *DBInterface_LoadSchemaFields_Call {
	_c.Run(run)
	return _c
}

// NewDBInterface creates a new instance of DBInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDBInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *DBInterface {
	mock := &DBInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
