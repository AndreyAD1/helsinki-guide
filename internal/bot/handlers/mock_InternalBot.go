// Code generated by mockery v2.38.0. DO NOT EDIT.

package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	mock "github.com/stretchr/testify/mock"
)

// InternalBot_mock is an autogenerated mock type for the InternalBot type
type InternalBot_mock struct {
	mock.Mock
}

type InternalBot_mock_Expecter struct {
	mock *mock.Mock
}

func (_m *InternalBot_mock) EXPECT() *InternalBot_mock_Expecter {
	return &InternalBot_mock_Expecter{mock: &_m.Mock}
}

// GetUpdatesChan provides a mock function with given fields: _a0
func (_m *InternalBot_mock) GetUpdatesChan(_a0 tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetUpdatesChan")
	}

	var r0 tgbotapi.UpdatesChannel
	if rf, ok := ret.Get(0).(func(tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(tgbotapi.UpdatesChannel)
		}
	}

	return r0
}

// InternalBot_mock_GetUpdatesChan_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUpdatesChan'
type InternalBot_mock_GetUpdatesChan_Call struct {
	*mock.Call
}

// GetUpdatesChan is a helper method to define mock.On call
//   - _a0 tgbotapi.UpdateConfig
func (_e *InternalBot_mock_Expecter) GetUpdatesChan(_a0 interface{}) *InternalBot_mock_GetUpdatesChan_Call {
	return &InternalBot_mock_GetUpdatesChan_Call{Call: _e.mock.On("GetUpdatesChan", _a0)}
}

func (_c *InternalBot_mock_GetUpdatesChan_Call) Run(run func(_a0 tgbotapi.UpdateConfig)) *InternalBot_mock_GetUpdatesChan_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(tgbotapi.UpdateConfig))
	})
	return _c
}

func (_c *InternalBot_mock_GetUpdatesChan_Call) Return(_a0 tgbotapi.UpdatesChannel) *InternalBot_mock_GetUpdatesChan_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *InternalBot_mock_GetUpdatesChan_Call) RunAndReturn(run func(tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel) *InternalBot_mock_GetUpdatesChan_Call {
	_c.Call.Return(run)
	return _c
}

// Request provides a mock function with given fields: _a0
func (_m *InternalBot_mock) Request(_a0 tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Request")
	}

	var r0 *tgbotapi.APIResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(tgbotapi.Chattable) (*tgbotapi.APIResponse, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(tgbotapi.Chattable) *tgbotapi.APIResponse); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tgbotapi.APIResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(tgbotapi.Chattable) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InternalBot_mock_Request_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Request'
type InternalBot_mock_Request_Call struct {
	*mock.Call
}

// Request is a helper method to define mock.On call
//   - _a0 tgbotapi.Chattable
func (_e *InternalBot_mock_Expecter) Request(_a0 interface{}) *InternalBot_mock_Request_Call {
	return &InternalBot_mock_Request_Call{Call: _e.mock.On("Request", _a0)}
}

func (_c *InternalBot_mock_Request_Call) Run(run func(_a0 tgbotapi.Chattable)) *InternalBot_mock_Request_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(tgbotapi.Chattable))
	})
	return _c
}

func (_c *InternalBot_mock_Request_Call) Return(_a0 *tgbotapi.APIResponse, _a1 error) *InternalBot_mock_Request_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *InternalBot_mock_Request_Call) RunAndReturn(run func(tgbotapi.Chattable) (*tgbotapi.APIResponse, error)) *InternalBot_mock_Request_Call {
	_c.Call.Return(run)
	return _c
}

// Send provides a mock function with given fields: _a0
func (_m *InternalBot_mock) Send(_a0 tgbotapi.Chattable) (tgbotapi.Message, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 tgbotapi.Message
	var r1 error
	if rf, ok := ret.Get(0).(func(tgbotapi.Chattable) (tgbotapi.Message, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(tgbotapi.Chattable) tgbotapi.Message); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(tgbotapi.Message)
	}

	if rf, ok := ret.Get(1).(func(tgbotapi.Chattable) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InternalBot_mock_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type InternalBot_mock_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - _a0 tgbotapi.Chattable
func (_e *InternalBot_mock_Expecter) Send(_a0 interface{}) *InternalBot_mock_Send_Call {
	return &InternalBot_mock_Send_Call{Call: _e.mock.On("Send", _a0)}
}

func (_c *InternalBot_mock_Send_Call) Run(run func(_a0 tgbotapi.Chattable)) *InternalBot_mock_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(tgbotapi.Chattable))
	})
	return _c
}

func (_c *InternalBot_mock_Send_Call) Return(_a0 tgbotapi.Message, _a1 error) *InternalBot_mock_Send_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *InternalBot_mock_Send_Call) RunAndReturn(run func(tgbotapi.Chattable) (tgbotapi.Message, error)) *InternalBot_mock_Send_Call {
	_c.Call.Return(run)
	return _c
}

// NewInternalBot_mock creates a new instance of InternalBot_mock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInternalBot_mock(t interface {
	mock.TestingT
	Cleanup(func())
}) *InternalBot_mock {
	mock := &InternalBot_mock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
