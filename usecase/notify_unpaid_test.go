package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/xgnid-tw/gx5/domain"
	"github.com/xgnid-tw/gx5/mocks"
	"github.com/xgnid-tw/gx5/usecase"
)

const testOthersDBID = "others-db"

var testLocation = time.UTC

func mockClock(t time.Time) *clock.Mock {
	clk := clock.NewMock()
	clk.Set(t)

	return clk
}

func TestExecute_NotFirstOfMonth_SkipsAll(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID, testLocation)
	uc.Clock = mockClock(time.Date(2026, 1, 10, 9, 0, 0, 0, time.UTC))

	err := uc.Execute(context.Background())

	require.NoError(t, err)
}

func TestExecute_GetUsersError(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	notifier := mocks.NewNotifier(t)

	repo.On("GetUsers", mock.Anything).Return(nil, errors.New("db error"))

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID, testLocation)
	uc.Clock = mockClock(time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC))

	err := uc.Execute(context.Background())

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

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID, testLocation)
	uc.Clock = mockClock(time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC))

	err := uc.Execute(context.Background())

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

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID, testLocation)
	uc.Clock = mockClock(time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC))

	err := uc.Execute(context.Background())

	require.Error(t, err)
	require.ErrorContains(t, err, "get others unpaid amount")
}

func TestExecute_PersonalDB_AboveThreshold_Notified(t *testing.T) {
	tests := []struct {
		name string
		day  int
	}{
		{"on 1st", 1},
		{"on 15th", 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			notifier := mocks.NewNotifier(t)

			user := &domain.User{
				DiscordID: "111", Name: "Alice",
				NotionID: "abc", Currency: domain.CurrencyTWD,
			}

			repo.On("GetUsers", mock.Anything).Return([]*domain.User{user}, nil)
			repo.On("GetUnpaidAmount", mock.Anything, "abc", domain.CurrencyTWD).
				Return(float64(3000), nil)
			notifier.On("Notify", mock.Anything, *user).Return(nil)

			uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID, testLocation)
			uc.Clock = mockClock(time.Date(2026, 1, tt.day, 9, 0, 0, 0, time.UTC))

			err := uc.Execute(context.Background())

			require.NoError(t, err)
		})
	}
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
	notifier.On("Notify", mock.Anything, *user).Return(nil)

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID, testLocation)
	uc.Clock = mockClock(time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC))

	err := uc.Execute(context.Background())

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

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID, testLocation)
	uc.Clock = mockClock(time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC))

	err := uc.Execute(context.Background())

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

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID, testLocation)
	uc.Clock = mockClock(time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC))

	err := uc.Execute(context.Background())

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
	notifier.On("Notify", mock.Anything, *user1).
		Return(errors.New("discord error"))
	notifier.On("Notify", mock.Anything, *user2).Return(nil)

	uc := usecase.NewNotifyUnpaid(repo, notifier, testOthersDBID, testLocation)
	uc.Clock = mockClock(time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC))

	err := uc.Execute(context.Background())

	require.NoError(t, err)
}
