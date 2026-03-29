package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/xgnid-tw/gx5/domain"
	"github.com/xgnid-tw/gx5/port"
)

type RegisterBuyRecord struct {
	userRepo     port.UserRepository
	txRepo       port.TransactionRepository
	jpyToTWDRate float64
}

func NewRegisterBuyRecord(
	userRepo port.UserRepository, txRepo port.TransactionRepository, jpyToTWDRate float64,
) *RegisterBuyRecord {
	return &RegisterBuyRecord{userRepo: userRepo, txRepo: txRepo, jpyToTWDRate: jpyToTWDRate}
}

func (uc *RegisterBuyRecord) Execute(
	ctx context.Context, targetDiscordID string, jpyAmount float64, itemName string,
) (*domain.BuyResult, error) {
	user, err := uc.userRepo.GetUserByDiscordID(ctx, targetDiscordID)
	if err != nil {
		return nil, fmt.Errorf("get user by discord id: %w", err)
	}

	twdAmount := math.Round(jpyAmount * uc.jpyToTWDRate)

	tx := domain.Transaction{
		ItemName:   itemName,
		JPYAmount:  jpyAmount,
		TWDAmount:  twdAmount,
		DatabaseID: user.NotionID,
	}

	err = uc.txRepo.CreateTransaction(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}

	displayAmount := twdAmount
	if user.Currency == domain.CurrencyJPY {
		displayAmount = jpyAmount
	}

	return &domain.BuyResult{
		DisplayAmount: displayAmount,
		Currency:      user.Currency,
		ItemName:      itemName,
	}, nil
}
