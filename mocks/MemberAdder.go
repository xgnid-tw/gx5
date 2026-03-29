// Code generated manually. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MemberAdder is a mock type for the MemberAdder type.
type MemberAdder struct {
	mock.Mock
}

// AddRoleMembersToThread provides a mock function with given fields: ctx, threadID, roleID.
func (_m *MemberAdder) AddRoleMembersToThread(
	ctx context.Context, threadID string, roleID string,
) error {
	ret := _m.Called(ctx, threadID, roleID)

	if len(ret) == 0 {
		panic("no return value specified for AddRoleMembersToThread")
	}

	var r0 error

	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, threadID, roleID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMemberAdder creates a new instance of MemberAdder.
func NewMemberAdder(t interface {
	mock.TestingT
	Cleanup(func())
}) *MemberAdder {
	mock := &MemberAdder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
