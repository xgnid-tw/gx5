package port

import "context"

// MemberAdder adds guild members with a specific role to a thread,
// then sends the formatted message after all members are added.
type MemberAdder interface {
	AddRoleMembersToThread(ctx context.Context, threadID string, roleID string, message string) error
}
