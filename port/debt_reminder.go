package port

import "context"

type DebtReminder interface {
	Execute(ctx context.Context, debug bool) error
}
