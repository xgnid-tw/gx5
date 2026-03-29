package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/xgnid-tw/gx5/domain"
	"github.com/xgnid-tw/gx5/mocks"
	"github.com/xgnid-tw/gx5/usecase"
)

func TestRegisterBuyRecord_Success(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	txRepo := mocks.NewTransactionRepository(t)

	user := &domain.User{
		DiscordID: "111", Name: "Alice",
		NotionID: "abc-db", Currency: domain.CurrencyJPY,
	}

	userRepo.On("GetUserByDiscordID", mock.Anything, "111").Return(user, nil)
	txRepo.On("CreateTransaction", mock.Anything, domain.Transaction{
		ItemName:   "Thread Title",
		JPYAmount:  3000,
		TWDAmount:  720,
		DatabaseID: "abc-db",
	}).Return(nil)

	uc := usecase.NewRegisterBuyRecord(userRepo, txRepo, 0.24)
	result, err := uc.Execute(context.Background(), "111", 3000, "Thread Title")

	require.NoError(t, err)
	require.Equal(t, float64(3000), result.DisplayAmount)
	require.Equal(t, domain.CurrencyJPY, result.Currency)
	require.Equal(t, "Thread Title", result.ItemName)
}

func TestRegisterBuyRecord_Success_TWDUser(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	txRepo := mocks.NewTransactionRepository(t)

	user := &domain.User{
		DiscordID: "222", Name: "Bob",
		NotionID: "bob-db", Currency: domain.CurrencyTWD,
	}

	userRepo.On("GetUserByDiscordID", mock.Anything, "222").Return(user, nil)
	txRepo.On("CreateTransaction", mock.Anything, domain.Transaction{
		ItemName:   "Item",
		JPYAmount:  3000,
		TWDAmount:  720,
		DatabaseID: "bob-db",
	}).Return(nil)

	uc := usecase.NewRegisterBuyRecord(userRepo, txRepo, 0.24)
	result, err := uc.Execute(context.Background(), "222", 3000, "Item")

	require.NoError(t, err)
	require.Equal(t, float64(720), result.DisplayAmount)
	require.Equal(t, domain.CurrencyTWD, result.Currency)
}

func TestRegisterBuyRecord_UserNotFound(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	txRepo := mocks.NewTransactionRepository(t)

	userRepo.On("GetUserByDiscordID", mock.Anything, "999").
		Return(nil, errors.New("user not found"))

	uc := usecase.NewRegisterBuyRecord(userRepo, txRepo, 0.24)
	result, err := uc.Execute(context.Background(), "999", 3000, "Thread Title")

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorContains(t, err, "get user by discord id")
}

func TestRegisterBuyRecord_CreateTransactionError(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	txRepo := mocks.NewTransactionRepository(t)

	user := &domain.User{
		DiscordID: "111", Name: "Alice",
		NotionID: "abc-db", Currency: domain.CurrencyJPY,
	}

	userRepo.On("GetUserByDiscordID", mock.Anything, "111").Return(user, nil)
	txRepo.On("CreateTransaction", mock.Anything, mock.Anything).
		Return(errors.New("notion error"))

	uc := usecase.NewRegisterBuyRecord(userRepo, txRepo, 0.24)
	result, err := uc.Execute(context.Background(), "111", 3000, "Thread Title")

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorContains(t, err, "create transaction")
}

func TestRegisterBuyRecord_ExchangeRate(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	txRepo := mocks.NewTransactionRepository(t)

	user := &domain.User{
		DiscordID: "111", Name: "Alice",
		NotionID: "abc-db", Currency: domain.CurrencyJPY,
	}

	userRepo.On("GetUserByDiscordID", mock.Anything, "111").Return(user, nil)
	txRepo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx domain.Transaction) bool {
		return tx.JPYAmount == 10000 && tx.TWDAmount == 2400
	})).Return(nil)

	uc := usecase.NewRegisterBuyRecord(userRepo, txRepo, 0.24)
	result, err := uc.Execute(context.Background(), "111", 10000, "Item")

	require.NoError(t, err)
	require.Equal(t, float64(10000), result.DisplayAmount)
}

func TestRegisterBuyRecord_TWDRounded(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	txRepo := mocks.NewTransactionRepository(t)

	user := &domain.User{
		DiscordID: "111", Name: "Alice",
		NotionID: "abc-db", Currency: domain.CurrencyJPY,
	}

	// 3500 * 0.217 = 759.5 → rounds to 760
	userRepo.On("GetUserByDiscordID", mock.Anything, "111").Return(user, nil)
	txRepo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx domain.Transaction) bool {
		return tx.JPYAmount == 3500 && tx.TWDAmount == 760
	})).Return(nil)

	uc := usecase.NewRegisterBuyRecord(userRepo, txRepo, 0.217)
	result, err := uc.Execute(context.Background(), "111", 3500, "Item")

	require.NoError(t, err)
	require.Equal(t, float64(3500), result.DisplayAmount)
}
