package port

import (
	"context"

	"github.com/xgnid-tw/gx5/domain"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]*domain.User, error)
	GetUserByDiscordID(ctx context.Context, discordID string) (*domain.User, error)
	GetUnpaidAmount(ctx context.Context, userDatabaseID string, currency domain.Currency) (float64, error)
	GetOthersUnpaidAmount(ctx context.Context, buyerName string, currency domain.Currency) (float64, error)
}
