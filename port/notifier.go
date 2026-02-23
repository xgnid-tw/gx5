package port

import (
	"context"

	"github.com/xgnid-tw/gx5/domain"
)

type Notifier interface {
	Notify(ctx context.Context, user domain.User) error
}
