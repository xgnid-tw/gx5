package usecase

import (
	"context"
	"fmt"

	"github.com/xgnid-tw/gx5/domain"
	"github.com/xgnid-tw/gx5/port"
)

const JPYToTWDRate = 0.24

type RegisterBuyRecord struct {
	userRepo port.UserRepository
	txRepo   port.TransactionRepository
}

func NewRegisterBuyRecord(
	userRepo port.UserRepository, txRepo port.TransactionRepository,
) *RegisterBuyRecord {
	return &RegisterBuyRecord{userRepo: userRepo, txRepo: txRepo}
}

func (uc *RegisterBuyRecord) Execute(
	ctx context.Context, targetDiscordID string, jpyAmount float64, itemName string,
) error {
	user, err := uc.userRepo.GetUserByDiscordID(ctx, targetDiscordID)
	if err != nil {
		return fmt.Errorf("get user by discord id: %w", err)
	}

	twdAmount := jpyAmount * JPYToTWDRate

	tx := domain.Transaction{
		ItemName:   itemName,
		JPYAmount:  jpyAmount,
		TWDAmount:  twdAmount,
		DatabaseID: user.NotionID,
	}

	err = uc.txRepo.CreateTransaction(ctx, tx)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	return nil
}
