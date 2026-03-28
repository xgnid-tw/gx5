package port

import (
	"context"

	"github.com/xgnid-tw/gx5/domain"
)

// OrderCreator creates a new order with a Discord thread and Notion record.
type OrderCreator interface {
	Execute(ctx context.Context, channelID string, order domain.Order) error
}
