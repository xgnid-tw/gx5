package port

import (
	"context"

	"github.com/xgnid-tw/gx5/domain"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order domain.Order) error
}
