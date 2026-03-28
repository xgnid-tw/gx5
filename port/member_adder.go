package port

import "context"

// MemberAdder adds guild members with a specific role to a thread.
type MemberAdder interface {
	AddRoleMembersToThread(ctx context.Context, threadID string, roleID string) error
}
