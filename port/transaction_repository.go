package port

import (
	"context"

	"github.com/xgnid-tw/gx5/domain"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx domain.Transaction) error
}
