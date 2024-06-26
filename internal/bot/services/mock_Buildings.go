// Code generated by mockery v2.43.2. DO NOT EDIT.

package services

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Buildings_mock is an autogenerated mock type for the Buildings type
type Buildings_mock struct {
	mock.Mock
}

type Buildings_mock_Expecter struct {
	mock *mock.Mock
}

func (_m *Buildings_mock) EXPECT() *Buildings_mock_Expecter {
	return &Buildings_mock_Expecter{mock: &_m.Mock}
}

// GetBuildingByID provides a mock function with given fields: c, ID
func (_m *Buildings_mock) GetBuildingByID(c context.Context, ID int64) (*BuildingDTO, error) {
	ret := _m.Called(c, ID)

	if len(ret) == 0 {
		panic("no return value specified for GetBuildingByID")
	}

	var r0 *BuildingDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*BuildingDTO, error)); ok {
		return rf(c, ID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *BuildingDTO); ok {
		r0 = rf(c, ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*BuildingDTO)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(c, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Buildings_mock_GetBuildingByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBuildingByID'
type Buildings_mock_GetBuildingByID_Call struct {
	*mock.Call
}

// GetBuildingByID is a helper method to define mock.On call
//   - c context.Context
//   - ID int64
func (_e *Buildings_mock_Expecter) GetBuildingByID(c interface{}, ID interface{}) *Buildings_mock_GetBuildingByID_Call {
	return &Buildings_mock_GetBuildingByID_Call{Call: _e.mock.On("GetBuildingByID", c, ID)}
}

func (_c *Buildings_mock_GetBuildingByID_Call) Run(run func(c context.Context, ID int64)) *Buildings_mock_GetBuildingByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *Buildings_mock_GetBuildingByID_Call) Return(_a0 *BuildingDTO, _a1 error) *Buildings_mock_GetBuildingByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Buildings_mock_GetBuildingByID_Call) RunAndReturn(run func(context.Context, int64) (*BuildingDTO, error)) *Buildings_mock_GetBuildingByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetBuildings provides a mock function with given fields: ctx, addressPrefix, limit, offset
func (_m *Buildings_mock) GetBuildings(ctx context.Context, addressPrefix string, limit int, offset int) ([]BuildingDTO, error) {
	ret := _m.Called(ctx, addressPrefix, limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetBuildings")
	}

	var r0 []BuildingDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) ([]BuildingDTO, error)); ok {
		return rf(ctx, addressPrefix, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []BuildingDTO); ok {
		r0 = rf(ctx, addressPrefix, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]BuildingDTO)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) error); ok {
		r1 = rf(ctx, addressPrefix, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Buildings_mock_GetBuildings_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBuildings'
type Buildings_mock_GetBuildings_Call struct {
	*mock.Call
}

// GetBuildings is a helper method to define mock.On call
//   - ctx context.Context
//   - addressPrefix string
//   - limit int
//   - offset int
func (_e *Buildings_mock_Expecter) GetBuildings(ctx interface{}, addressPrefix interface{}, limit interface{}, offset interface{}) *Buildings_mock_GetBuildings_Call {
	return &Buildings_mock_GetBuildings_Call{Call: _e.mock.On("GetBuildings", ctx, addressPrefix, limit, offset)}
}

func (_c *Buildings_mock_GetBuildings_Call) Run(run func(ctx context.Context, addressPrefix string, limit int, offset int)) *Buildings_mock_GetBuildings_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(int), args[3].(int))
	})
	return _c
}

func (_c *Buildings_mock_GetBuildings_Call) Return(_a0 []BuildingDTO, _a1 error) *Buildings_mock_GetBuildings_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Buildings_mock_GetBuildings_Call) RunAndReturn(run func(context.Context, string, int, int) ([]BuildingDTO, error)) *Buildings_mock_GetBuildings_Call {
	_c.Call.Return(run)
	return _c
}

// GetBuildingsByAddress provides a mock function with given fields: c, address
func (_m *Buildings_mock) GetBuildingsByAddress(c context.Context, address string) ([]BuildingDTO, error) {
	ret := _m.Called(c, address)

	if len(ret) == 0 {
		panic("no return value specified for GetBuildingsByAddress")
	}

	var r0 []BuildingDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]BuildingDTO, error)); ok {
		return rf(c, address)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []BuildingDTO); ok {
		r0 = rf(c, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]BuildingDTO)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(c, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Buildings_mock_GetBuildingsByAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBuildingsByAddress'
type Buildings_mock_GetBuildingsByAddress_Call struct {
	*mock.Call
}

// GetBuildingsByAddress is a helper method to define mock.On call
//   - c context.Context
//   - address string
func (_e *Buildings_mock_Expecter) GetBuildingsByAddress(c interface{}, address interface{}) *Buildings_mock_GetBuildingsByAddress_Call {
	return &Buildings_mock_GetBuildingsByAddress_Call{Call: _e.mock.On("GetBuildingsByAddress", c, address)}
}

func (_c *Buildings_mock_GetBuildingsByAddress_Call) Run(run func(c context.Context, address string)) *Buildings_mock_GetBuildingsByAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Buildings_mock_GetBuildingsByAddress_Call) Return(_a0 []BuildingDTO, _a1 error) *Buildings_mock_GetBuildingsByAddress_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Buildings_mock_GetBuildingsByAddress_Call) RunAndReturn(run func(context.Context, string) ([]BuildingDTO, error)) *Buildings_mock_GetBuildingsByAddress_Call {
	_c.Call.Return(run)
	return _c
}

// GetNearestBuildings provides a mock function with given fields: ctx, distance, latitude, longitude, limit, offset
func (_m *Buildings_mock) GetNearestBuildings(ctx context.Context, distance int, latitude float64, longitude float64, limit int, offset int) ([]BuildingDTO, error) {
	ret := _m.Called(ctx, distance, latitude, longitude, limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetNearestBuildings")
	}

	var r0 []BuildingDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, float64, float64, int, int) ([]BuildingDTO, error)); ok {
		return rf(ctx, distance, latitude, longitude, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, float64, float64, int, int) []BuildingDTO); ok {
		r0 = rf(ctx, distance, latitude, longitude, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]BuildingDTO)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, float64, float64, int, int) error); ok {
		r1 = rf(ctx, distance, latitude, longitude, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Buildings_mock_GetNearestBuildings_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNearestBuildings'
type Buildings_mock_GetNearestBuildings_Call struct {
	*mock.Call
}

// GetNearestBuildings is a helper method to define mock.On call
//   - ctx context.Context
//   - distance int
//   - latitude float64
//   - longitude float64
//   - limit int
//   - offset int
func (_e *Buildings_mock_Expecter) GetNearestBuildings(ctx interface{}, distance interface{}, latitude interface{}, longitude interface{}, limit interface{}, offset interface{}) *Buildings_mock_GetNearestBuildings_Call {
	return &Buildings_mock_GetNearestBuildings_Call{Call: _e.mock.On("GetNearestBuildings", ctx, distance, latitude, longitude, limit, offset)}
}

func (_c *Buildings_mock_GetNearestBuildings_Call) Run(run func(ctx context.Context, distance int, latitude float64, longitude float64, limit int, offset int)) *Buildings_mock_GetNearestBuildings_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(float64), args[3].(float64), args[4].(int), args[5].(int))
	})
	return _c
}

func (_c *Buildings_mock_GetNearestBuildings_Call) Return(_a0 []BuildingDTO, _a1 error) *Buildings_mock_GetNearestBuildings_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Buildings_mock_GetNearestBuildings_Call) RunAndReturn(run func(context.Context, int, float64, float64, int, int) ([]BuildingDTO, error)) *Buildings_mock_GetNearestBuildings_Call {
	_c.Call.Return(run)
	return _c
}

// NewBuildings_mock creates a new instance of Buildings_mock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBuildings_mock(t interface {
	mock.TestingT
	Cleanup(func())
}) *Buildings_mock {
	mock := &Buildings_mock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
