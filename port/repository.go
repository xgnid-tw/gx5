package port

import (
	"context"
	"time"

	"github.com/xgnid-tw/gx5/domain"
)

type UnpaidSummary struct {
	TotalAmount    float64
	OldestRecordAt time.Time
}

type UserRepository interface {
	GetUsers(ctx context.Context) ([]*domain.User, error)
	GetUserByDiscordID(ctx context.Context, discordID string) (*domain.User, error)
	GetUnpaidSummary(ctx context.Context, userDatabaseID string, currency domain.Currency) (*UnpaidSummary, error)
}
