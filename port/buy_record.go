package port

import (
	"context"

	"github.com/xgnid-tw/gx5/domain"
)

// BuyRecordRegisterer abstracts the register-buy-record use case for the gateway layer.
type BuyRecordRegisterer interface {
	Execute(
		ctx context.Context, targetDiscordID string, jpyAmount float64, itemName string,
	) (*domain.BuyResult, error)
}
