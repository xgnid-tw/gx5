package port

import "context"

// BuyRecordRegisterer abstracts the register-buy-record use case for the gateway layer.
type BuyRecordRegisterer interface {
	Execute(ctx context.Context, targetDiscordID string, jpyAmount float64, itemName string) error
}
