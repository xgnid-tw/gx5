package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/xgnid-tw/gx5/domain"
	"github.com/xgnid-tw/gx5/mocks"
	"github.com/xgnid-tw/gx5/port"
	"github.com/xgnid-tw/gx5/usecase"
)

func TestExecute_GetUsersError(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	repo.On("GetUsers", mock.Anything).Return(nil, errors.New("db error"))

	uc := usecase.NewNotifyUnpaid(repo, notifier)

	err := uc.Execute(context.Background(), false)

	require.Error(t, err)
	require.ErrorContains(t, err, "get users")
}

func TestExecute_GetUnpaidSummaryError(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	user := &domain.User{
		DiscordID: "111", Name: "Alice",
		NotionID: "abc", Currency: domain.CurrencyTWD,
	}

	repo.On("GetUsers", mock.Anything).Return([]*domain.User{user}, nil)
	repo.On("GetUnpaidSummary", mock.Anything, "abc", domain.CurrencyTWD).
		Return(nil, errors.New("notion error"))

	uc := usecase.NewNotifyUnpaid(repo, notifier)

	err := uc.Execute(context.Background(), false)

	require.Error(t, err)
	require.ErrorContains(t, err, "get unpaid summary")
}

func TestExecute_AboveThreshold_Notified(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	user := &domain.User{
		DiscordID: "111", Name: "Alice",
		NotionID: "abc", Currency: domain.CurrencyTWD,
	}

	repo.On("GetUsers", mock.Anything).Return([]*domain.User{user}, nil)
	repo.On("GetUnpaidSummary", mock.Anything, "abc", domain.CurrencyTWD).
		Return(&port.UnpaidSummary{TotalAmount: 3000, OldestRecordAt: time.Now()}, nil)
	notifier.On("Notify", mock.Anything, *user, false).Return(nil)

	uc := usecase.NewNotifyUnpaid(repo, notifier)

	err := uc.Execute(context.Background(), false)

	require.NoError(t, err)
}

func TestExecute_BelowThreshold_NotNotified(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	user := &domain.User{
		DiscordID: "111", Name: "Alice",
		NotionID: "abc", Currency: domain.CurrencyTWD,
	}

	repo.On("GetUsers", mock.Anything).Return([]*domain.User{user}, nil)
	repo.On("GetUnpaidSummary", mock.Anything, "abc", domain.CurrencyTWD).
		Return(&port.UnpaidSummary{TotalAmount: 500, OldestRecordAt: time.Now()}, nil)

	uc := usecase.NewNotifyUnpaid(repo, notifier)

	err := uc.Execute(context.Background(), false)

	require.NoError(t, err)
}

func TestExecute_ZeroAmount_NotNotified(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	user := &domain.User{
		DiscordID: "111", Name: "Alice",
		NotionID: "abc", Currency: domain.CurrencyTWD,
	}

	repo.On("GetUsers", mock.Anything).Return([]*domain.User{user}, nil)
	repo.On("GetUnpaidSummary", mock.Anything, "abc", domain.CurrencyTWD).
		Return(&port.UnpaidSummary{}, nil)

	uc := usecase.NewNotifyUnpaid(repo, notifier)

	err := uc.Execute(context.Background(), false)

	require.NoError(t, err)
}

func TestExecute_BelowThreshold_OlderThan3Months_Notified(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	user := &domain.User{
		DiscordID: "111", Name: "Alice",
		NotionID: "abc", Currency: domain.CurrencyTWD,
	}

	fourMonthsAgo := time.Now().AddDate(0, -4, 0)

	repo.On("GetUsers", mock.Anything).Return([]*domain.User{user}, nil)
	repo.On("GetUnpaidSummary", mock.Anything, "abc", domain.CurrencyTWD).
		Return(&port.UnpaidSummary{TotalAmount: 500, OldestRecordAt: fourMonthsAgo}, nil)
	notifier.On("Notify", mock.Anything, *user, false).Return(nil)

	uc := usecase.NewNotifyUnpaid(repo, notifier)

	err := uc.Execute(context.Background(), false)

	require.NoError(t, err)
}

func TestExecute_BelowThreshold_Within3Months_NotNotified(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	user := &domain.User{
		DiscordID: "111", Name: "Alice",
		NotionID: "abc", Currency: domain.CurrencyTWD,
	}

	twoMonthsAgo := time.Now().AddDate(0, -2, 0)

	repo.On("GetUsers", mock.Anything).Return([]*domain.User{user}, nil)
	repo.On("GetUnpaidSummary", mock.Anything, "abc", domain.CurrencyTWD).
		Return(&port.UnpaidSummary{TotalAmount: 500, OldestRecordAt: twoMonthsAgo}, nil)

	uc := usecase.NewNotifyUnpaid(repo, notifier)

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
	repo.On("GetUnpaidSummary", mock.Anything, "abc", domain.CurrencyTWD).
		Return(&port.UnpaidSummary{TotalAmount: 3000, OldestRecordAt: time.Now()}, nil)
	repo.On("GetUnpaidSummary", mock.Anything, "def", domain.CurrencyTWD).
		Return(&port.UnpaidSummary{TotalAmount: 3000, OldestRecordAt: time.Now()}, nil)
	notifier.On("Notify", mock.Anything, *user1, false).
		Return(errors.New("discord error"))
	notifier.On("Notify", mock.Anything, *user2, false).Return(nil)

	uc := usecase.NewNotifyUnpaid(repo, notifier)

	err := uc.Execute(context.Background(), false)

	require.NoError(t, err)
}
