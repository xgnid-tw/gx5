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

const testOthersDBID = "others-db"

func TestExecute_GetUsersError(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	repo.On("GetUsers", mock.Anything).Return(nil, errors.New("db error"))

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID)

	err := uc.Execute(context.Background(), false)

	require.Error(t, err)
	require.ErrorContains(t, err, "get users")
}

func TestExecute_GetUnpaidAmountError(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	user := &domain.User{
		DiscordID: "111", Name: "Alice",
		NotionID: "abc", Currency: domain.CurrencyTWD,
	}

	repo.On("GetUsers", mock.Anything).Return([]*domain.User{user}, nil)
	repo.On("GetUnpaidAmount", mock.Anything, "abc", domain.CurrencyTWD).
		Return(float64(0), errors.New("notion error"))

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID)

	err := uc.Execute(context.Background(), false)

	require.Error(t, err)
	require.ErrorContains(t, err, "get unpaid amount")
}

func TestExecute_GetOthersUnpaidAmountError(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	user := &domain.User{
		DiscordID: "111", Name: "Alice",
		NotionID: testOthersDBID, Currency: domain.CurrencyTWD,
	}

	repo.On("GetUsers", mock.Anything).Return([]*domain.User{user}, nil)
	repo.On("GetOthersUnpaidAmount", mock.Anything, "Alice", domain.CurrencyTWD).
		Return(float64(0), errors.New("notion error"))

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID)

	err := uc.Execute(context.Background(), false)

	require.Error(t, err)
	require.ErrorContains(t, err, "get others unpaid amount")
}

func TestExecute_PersonalDB_AboveThreshold_Notified(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	user := &domain.User{
		DiscordID: "111", Name: "Alice",
		NotionID: "abc", Currency: domain.CurrencyTWD,
	}

	repo.On("GetUsers", mock.Anything).Return([]*domain.User{user}, nil)
	repo.On("GetUnpaidAmount", mock.Anything, "abc", domain.CurrencyTWD).
		Return(float64(3000), nil)
	notifier.On("Notify", mock.Anything, *user, false).Return(nil)

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID)

	err := uc.Execute(context.Background(), false)

	require.NoError(t, err)
}

func TestExecute_OthersDB_AboveThreshold_Notified(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	user := &domain.User{
		DiscordID: "333", Name: "Carol",
		NotionID: testOthersDBID, Currency: domain.CurrencyTWD,
	}

	repo.On("GetUsers", mock.Anything).Return([]*domain.User{user}, nil)
	repo.On("GetOthersUnpaidAmount", mock.Anything, "Carol", domain.CurrencyTWD).
		Return(float64(2500), nil)
	notifier.On("Notify", mock.Anything, *user, false).Return(nil)

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID)

	err := uc.Execute(context.Background(), false)

	require.NoError(t, err)
}

func TestExecute_PersonalDB_ZeroAmount_NotNotified(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	user := &domain.User{
		DiscordID: "111", Name: "Alice",
		NotionID: "abc", Currency: domain.CurrencyTWD,
	}

	repo.On("GetUsers", mock.Anything).Return([]*domain.User{user}, nil)
	repo.On("GetUnpaidAmount", mock.Anything, "abc", domain.CurrencyTWD).
		Return(float64(0), nil)

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID)

	err := uc.Execute(context.Background(), false)

	require.NoError(t, err)
}

func TestExecute_OthersDB_ZeroAmount_NotNotified(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	user := &domain.User{
		DiscordID: "333", Name: "Carol",
		NotionID: testOthersDBID, Currency: domain.CurrencyTWD,
	}

	repo.On("GetUsers", mock.Anything).Return([]*domain.User{user}, nil)
	repo.On("GetOthersUnpaidAmount", mock.Anything, "Carol", domain.CurrencyTWD).
		Return(float64(0), nil)

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID)

	err := uc.Execute(context.Background(), false)

	require.NoError(t, err)
}

func TestExecute_NotifyError_ContinuesNextUser(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	user1 := &domain.User{
		DiscordID: "111", Name: "Alice",
		NotionID: "abc", Currency: domain.CurrencyTWD,
	}
	user2 := &domain.User{
		DiscordID: "222", Name: "Bob",
		NotionID: "def", Currency: domain.CurrencyTWD,
	}

	repo.On("GetUsers", mock.Anything).Return([]*domain.User{user1, user2}, nil)
	repo.On("GetUnpaidAmount", mock.Anything, "abc", domain.CurrencyTWD).
		Return(float64(3000), nil)
	repo.On("GetUnpaidAmount", mock.Anything, "def", domain.CurrencyTWD).
		Return(float64(3000), nil)
	notifier.On("Notify", mock.Anything, *user1, false).
		Return(errors.New("discord error"))
	notifier.On("Notify", mock.Anything, *user2, false).Return(nil)

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID)

	err := uc.Execute(context.Background(), false)

	require.NoError(t, err)
}
