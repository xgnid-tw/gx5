// Code generated manually. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ThreadCreator is a mock type for the ThreadCreator type.
type ThreadCreator struct {
	mock.Mock
}

// CreateThread provides a mock function with given fields: ctx, channelID, name.
func (_m *ThreadCreator) CreateThread(
	ctx context.Context, channelID string, name string,
) (string, error) {
	ret := _m.Called(ctx, channelID, name)

	if len(ret) == 0 {
		panic("no return value specified for CreateThread")
	}

	var r0 string
	var r1 error

	if rf, ok := ret.Get(0).(func(context.Context, string, string) (string, error)); ok {
		return rf(ctx, channelID, name)
	}

	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, channelID, name)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, channelID, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendThreadMessage provides a mock function with given fields: ctx, threadID, message.
func (_m *ThreadCreator) SendThreadMessage(
	ctx context.Context, threadID string, message string,
) error {
	ret := _m.Called(ctx, threadID, message)

	if len(ret) == 0 {
		panic("no return value specified for SendThreadMessage")
	}

	var r0 error

	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, threadID, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewThreadCreator creates a new instance of ThreadCreator.
func NewThreadCreator(t interface {
	mock.TestingT
	Cleanup(func())
}) *ThreadCreator {
	mock := &ThreadCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
