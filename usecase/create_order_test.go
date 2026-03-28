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

func TestCreateOrder_MissingOrderTitle(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)

	uc := usecase.NewCreateOrder(repo, tc)

	err := uc.Execute(context.Background(), "ch-1", domain.Order{})

	require.Error(t, err)
	require.ErrorContains(t, err, "orderTitle is required")
}

func TestCreateOrder_ThreadCreationError(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)

	tc.On("CreateThread", mock.Anything, "ch-1", "test order", mock.Anything).
		Return(errors.New("discord error"))

	uc := usecase.NewCreateOrder(repo, tc)

	err := uc.Execute(context.Background(), "ch-1", domain.Order{
		ThreadName: "test order",
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "create thread")
}

func TestCreateOrder_NotionError(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)

	tc.On("CreateThread", mock.Anything, "ch-1", "test order", mock.Anything).
		Return(nil)
	repo.On("CreateOrder", mock.Anything, domain.Order{ThreadName: "test order"}).
		Return(errors.New("notion error"))

	uc := usecase.NewCreateOrder(repo, tc)

	err := uc.Execute(context.Background(), "ch-1", domain.Order{
		ThreadName: "test order",
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "create order record")
}

func TestCreateOrder_Success_AllFields(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)

	order := domain.Order{
		ThreadName: "test order",
		Deadline:   "2026-04-01",
		ShopURL:    "https://shop.example.com",
		Tag:        domain.Tag315Pro,
	}

	expectedMessage := "https://shop.example.com\n@315pro\n截止時間: 2026-04-01"

	tc.On("CreateThread", mock.Anything, "ch-1", "test order", expectedMessage).
		Return(nil)
	repo.On("CreateOrder", mock.Anything, order).
		Return(nil)

	uc := usecase.NewCreateOrder(repo, tc)

	err := uc.Execute(context.Background(), "ch-1", order)

	require.NoError(t, err)
}

func TestCreateOrder_Success_OnlyTitle(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)

	order := domain.Order{
		ThreadName: "minimal order",
	}

	tc.On("CreateThread", mock.Anything, "ch-1", "minimal order", "").
		Return(nil)
	repo.On("CreateOrder", mock.Anything, order).
		Return(nil)

	uc := usecase.NewCreateOrder(repo, tc)

	err := uc.Execute(context.Background(), "ch-1", order)

	require.NoError(t, err)
}

func TestCreateOrder_Success_PartialFields(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)

	order := domain.Order{
		ThreadName: "partial order",
		Deadline:   "2026-05-15",
	}

	expectedMessage := "截止時間: 2026-05-15"

	tc.On("CreateThread", mock.Anything, "ch-1", "partial order", expectedMessage).
		Return(nil)
	repo.On("CreateOrder", mock.Anything, order).
		Return(nil)

	uc := usecase.NewCreateOrder(repo, tc)

	err := uc.Execute(context.Background(), "ch-1", order)

	require.NoError(t, err)
}

func TestCreateOrder_Success_ShopURLOnly(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)

	order := domain.Order{
		ThreadName: "url order",
		ShopURL:    "https://example.com",
	}

	expectedMessage := "https://example.com"

	tc.On("CreateThread", mock.Anything, "ch-1", "url order", expectedMessage).
		Return(nil)
	repo.On("CreateOrder", mock.Anything, order).
		Return(nil)

	uc := usecase.NewCreateOrder(repo, tc)

	err := uc.Execute(context.Background(), "ch-1", order)

	require.NoError(t, err)
}

func TestCreateOrder_Success_TagOnly(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)

	order := domain.Order{
		ThreadName: "tag order",
		Tag:        domain.TagGakumas,
	}

	expectedMessage := "@学マス"

	tc.On("CreateThread", mock.Anything, "ch-1", "tag order", expectedMessage).
		Return(nil)
	repo.On("CreateOrder", mock.Anything, order).
		Return(nil)

	uc := usecase.NewCreateOrder(repo, tc)

	err := uc.Execute(context.Background(), "ch-1", order)

	require.NoError(t, err)
}
