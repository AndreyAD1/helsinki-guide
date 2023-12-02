// Code generated by mockery v2.38.0. DO NOT EDIT.

package repositories

import (
	context "context"

	specifications "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/specifications"
	types "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/types"
	mock "github.com/stretchr/testify/mock"
)

// NeighbourhoodRepository_mock is an autogenerated mock type for the NeighbourhoodRepository type
type NeighbourhoodRepository_mock struct {
	mock.Mock
}

type NeighbourhoodRepository_mock_Expecter struct {
	mock *mock.Mock
}

func (_m *NeighbourhoodRepository_mock) EXPECT() *NeighbourhoodRepository_mock_Expecter {
	return &NeighbourhoodRepository_mock_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: _a0, _a1
func (_m *NeighbourhoodRepository_mock) Add(_a0 context.Context, _a1 types.Neighbourhood) (*types.Neighbourhood, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 *types.Neighbourhood
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Neighbourhood) (*types.Neighbourhood, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.Neighbourhood) *types.Neighbourhood); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Neighbourhood)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.Neighbourhood) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NeighbourhoodRepository_mock_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type NeighbourhoodRepository_mock_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 types.Neighbourhood
func (_e *NeighbourhoodRepository_mock_Expecter) Add(_a0 interface{}, _a1 interface{}) *NeighbourhoodRepository_mock_Add_Call {
	return &NeighbourhoodRepository_mock_Add_Call{Call: _e.mock.On("Add", _a0, _a1)}
}

func (_c *NeighbourhoodRepository_mock_Add_Call) Run(run func(_a0 context.Context, _a1 types.Neighbourhood)) *NeighbourhoodRepository_mock_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Neighbourhood))
	})
	return _c
}

func (_c *NeighbourhoodRepository_mock_Add_Call) Return(_a0 *types.Neighbourhood, _a1 error) *NeighbourhoodRepository_mock_Add_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NeighbourhoodRepository_mock_Add_Call) RunAndReturn(run func(context.Context, types.Neighbourhood) (*types.Neighbourhood, error)) *NeighbourhoodRepository_mock_Add_Call {
	_c.Call.Return(run)
	return _c
}

// Query provides a mock function with given fields: _a0, _a1
func (_m *NeighbourhoodRepository_mock) Query(_a0 context.Context, _a1 specifications.Specification) ([]types.Neighbourhood, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Query")
	}

	var r0 []types.Neighbourhood
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, specifications.Specification) ([]types.Neighbourhood, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, specifications.Specification) []types.Neighbourhood); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Neighbourhood)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, specifications.Specification) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NeighbourhoodRepository_mock_Query_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Query'
type NeighbourhoodRepository_mock_Query_Call struct {
	*mock.Call
}

// Query is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 specifications.Specification
func (_e *NeighbourhoodRepository_mock_Expecter) Query(_a0 interface{}, _a1 interface{}) *NeighbourhoodRepository_mock_Query_Call {
	return &NeighbourhoodRepository_mock_Query_Call{Call: _e.mock.On("Query", _a0, _a1)}
}

func (_c *NeighbourhoodRepository_mock_Query_Call) Run(run func(_a0 context.Context, _a1 specifications.Specification)) *NeighbourhoodRepository_mock_Query_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(specifications.Specification))
	})
	return _c
}

func (_c *NeighbourhoodRepository_mock_Query_Call) Return(_a0 []types.Neighbourhood, _a1 error) *NeighbourhoodRepository_mock_Query_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NeighbourhoodRepository_mock_Query_Call) RunAndReturn(run func(context.Context, specifications.Specification) ([]types.Neighbourhood, error)) *NeighbourhoodRepository_mock_Query_Call {
	_c.Call.Return(run)
	return _c
}

// Remove provides a mock function with given fields: _a0, _a1
func (_m *NeighbourhoodRepository_mock) Remove(_a0 context.Context, _a1 types.Neighbourhood) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Remove")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Neighbourhood) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NeighbourhoodRepository_mock_Remove_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Remove'
type NeighbourhoodRepository_mock_Remove_Call struct {
	*mock.Call
}

// Remove is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 types.Neighbourhood
func (_e *NeighbourhoodRepository_mock_Expecter) Remove(_a0 interface{}, _a1 interface{}) *NeighbourhoodRepository_mock_Remove_Call {
	return &NeighbourhoodRepository_mock_Remove_Call{Call: _e.mock.On("Remove", _a0, _a1)}
}

func (_c *NeighbourhoodRepository_mock_Remove_Call) Run(run func(_a0 context.Context, _a1 types.Neighbourhood)) *NeighbourhoodRepository_mock_Remove_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Neighbourhood))
	})
	return _c
}

func (_c *NeighbourhoodRepository_mock_Remove_Call) Return(_a0 error) *NeighbourhoodRepository_mock_Remove_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *NeighbourhoodRepository_mock_Remove_Call) RunAndReturn(run func(context.Context, types.Neighbourhood) error) *NeighbourhoodRepository_mock_Remove_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *NeighbourhoodRepository_mock) Update(_a0 context.Context, _a1 types.Neighbourhood) (*types.Neighbourhood, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *types.Neighbourhood
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Neighbourhood) (*types.Neighbourhood, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.Neighbourhood) *types.Neighbourhood); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Neighbourhood)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.Neighbourhood) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NeighbourhoodRepository_mock_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type NeighbourhoodRepository_mock_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 types.Neighbourhood
func (_e *NeighbourhoodRepository_mock_Expecter) Update(_a0 interface{}, _a1 interface{}) *NeighbourhoodRepository_mock_Update_Call {
	return &NeighbourhoodRepository_mock_Update_Call{Call: _e.mock.On("Update", _a0, _a1)}
}

func (_c *NeighbourhoodRepository_mock_Update_Call) Run(run func(_a0 context.Context, _a1 types.Neighbourhood)) *NeighbourhoodRepository_mock_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Neighbourhood))
	})
	return _c
}

func (_c *NeighbourhoodRepository_mock_Update_Call) Return(_a0 *types.Neighbourhood, _a1 error) *NeighbourhoodRepository_mock_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *NeighbourhoodRepository_mock_Update_Call) RunAndReturn(run func(context.Context, types.Neighbourhood) (*types.Neighbourhood, error)) *NeighbourhoodRepository_mock_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewNeighbourhoodRepository_mock creates a new instance of NeighbourhoodRepository_mock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNeighbourhoodRepository_mock(t interface {
	mock.TestingT
	Cleanup(func())
}) *NeighbourhoodRepository_mock {
	mock := &NeighbourhoodRepository_mock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}