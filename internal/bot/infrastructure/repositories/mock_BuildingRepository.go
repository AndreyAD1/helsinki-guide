// Code generated by mockery v2.38.0. DO NOT EDIT.

package repositories

import (
	context "context"

	specifications "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/specifications"
	types "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/types"
	mock "github.com/stretchr/testify/mock"
)

// BuildingRepository_mock is an autogenerated mock type for the BuildingRepository type
type BuildingRepository_mock struct {
	mock.Mock
}

type BuildingRepository_mock_Expecter struct {
	mock *mock.Mock
}

func (_m *BuildingRepository_mock) EXPECT() *BuildingRepository_mock_Expecter {
	return &BuildingRepository_mock_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: _a0, _a1
func (_m *BuildingRepository_mock) Add(_a0 context.Context, _a1 types.Building) (*types.Building, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 *types.Building
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Building) (*types.Building, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.Building) *types.Building); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Building)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.Building) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BuildingRepository_mock_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type BuildingRepository_mock_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 types.Building
func (_e *BuildingRepository_mock_Expecter) Add(_a0 interface{}, _a1 interface{}) *BuildingRepository_mock_Add_Call {
	return &BuildingRepository_mock_Add_Call{Call: _e.mock.On("Add", _a0, _a1)}
}

func (_c *BuildingRepository_mock_Add_Call) Run(run func(_a0 context.Context, _a1 types.Building)) *BuildingRepository_mock_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Building))
	})
	return _c
}

func (_c *BuildingRepository_mock_Add_Call) Return(_a0 *types.Building, _a1 error) *BuildingRepository_mock_Add_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BuildingRepository_mock_Add_Call) RunAndReturn(run func(context.Context, types.Building) (*types.Building, error)) *BuildingRepository_mock_Add_Call {
	_c.Call.Return(run)
	return _c
}

// Query provides a mock function with given fields: _a0, _a1
func (_m *BuildingRepository_mock) Query(_a0 context.Context, _a1 specifications.Specification) ([]types.Building, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Query")
	}

	var r0 []types.Building
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, specifications.Specification) ([]types.Building, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, specifications.Specification) []types.Building); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Building)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, specifications.Specification) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BuildingRepository_mock_Query_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Query'
type BuildingRepository_mock_Query_Call struct {
	*mock.Call
}

// Query is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 specifications.Specification
func (_e *BuildingRepository_mock_Expecter) Query(_a0 interface{}, _a1 interface{}) *BuildingRepository_mock_Query_Call {
	return &BuildingRepository_mock_Query_Call{Call: _e.mock.On("Query", _a0, _a1)}
}

func (_c *BuildingRepository_mock_Query_Call) Run(run func(_a0 context.Context, _a1 specifications.Specification)) *BuildingRepository_mock_Query_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(specifications.Specification))
	})
	return _c
}

func (_c *BuildingRepository_mock_Query_Call) Return(_a0 []types.Building, _a1 error) *BuildingRepository_mock_Query_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BuildingRepository_mock_Query_Call) RunAndReturn(run func(context.Context, specifications.Specification) ([]types.Building, error)) *BuildingRepository_mock_Query_Call {
	_c.Call.Return(run)
	return _c
}

// Remove provides a mock function with given fields: _a0, _a1
func (_m *BuildingRepository_mock) Remove(_a0 context.Context, _a1 types.Building) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Remove")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Building) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BuildingRepository_mock_Remove_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Remove'
type BuildingRepository_mock_Remove_Call struct {
	*mock.Call
}

// Remove is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 types.Building
func (_e *BuildingRepository_mock_Expecter) Remove(_a0 interface{}, _a1 interface{}) *BuildingRepository_mock_Remove_Call {
	return &BuildingRepository_mock_Remove_Call{Call: _e.mock.On("Remove", _a0, _a1)}
}

func (_c *BuildingRepository_mock_Remove_Call) Run(run func(_a0 context.Context, _a1 types.Building)) *BuildingRepository_mock_Remove_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Building))
	})
	return _c
}

func (_c *BuildingRepository_mock_Remove_Call) Return(_a0 error) *BuildingRepository_mock_Remove_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BuildingRepository_mock_Remove_Call) RunAndReturn(run func(context.Context, types.Building) error) *BuildingRepository_mock_Remove_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *BuildingRepository_mock) Update(_a0 context.Context, _a1 types.Building) (*types.Building, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *types.Building
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Building) (*types.Building, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.Building) *types.Building); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Building)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.Building) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BuildingRepository_mock_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type BuildingRepository_mock_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 types.Building
func (_e *BuildingRepository_mock_Expecter) Update(_a0 interface{}, _a1 interface{}) *BuildingRepository_mock_Update_Call {
	return &BuildingRepository_mock_Update_Call{Call: _e.mock.On("Update", _a0, _a1)}
}

func (_c *BuildingRepository_mock_Update_Call) Run(run func(_a0 context.Context, _a1 types.Building)) *BuildingRepository_mock_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Building))
	})
	return _c
}

func (_c *BuildingRepository_mock_Update_Call) Return(_a0 *types.Building, _a1 error) *BuildingRepository_mock_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BuildingRepository_mock_Update_Call) RunAndReturn(run func(context.Context, types.Building) (*types.Building, error)) *BuildingRepository_mock_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewBuildingRepository_mock creates a new instance of BuildingRepository_mock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBuildingRepository_mock(t interface {
	mock.TestingT
	Cleanup(func())
}) *BuildingRepository_mock {
	mock := &BuildingRepository_mock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
