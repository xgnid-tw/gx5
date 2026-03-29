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

	uc := usecase.NewCreateOrder(repo, tc, nil, nil)

	err := uc.Execute(context.Background(), "ch-1", domain.Order{})

	require.Error(t, err)
	require.ErrorContains(t, err, "orderTitle is required")
}

func TestCreateOrder_ThreadCreationError(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)

	tc.On("CreateThread", mock.Anything, "ch-1", "test order").
		Return("", errors.New("discord error"))

	uc := usecase.NewCreateOrder(repo, tc, nil, nil)

	err := uc.Execute(context.Background(), "ch-1", domain.Order{
		ThreadName: "test order",
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "create thread")
}

func TestCreateOrder_NotionError(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)

	tc.On("CreateThread", mock.Anything, "ch-1", "test order").
		Return("thread-id", nil)
	repo.On("CreateOrder", mock.Anything, domain.Order{ThreadName: "test order"}).
		Return(errors.New("notion error"))

	uc := usecase.NewCreateOrder(repo, tc, nil, nil)

	err := uc.Execute(context.Background(), "ch-1", domain.Order{
		ThreadName: "test order",
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "create order record")
}

func TestCreateOrder_Success_AllFields(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)
	ma := mocks.NewMemberAdder(t)

	order := domain.Order{
		ThreadName: "test order",
		Deadline:   "2026-04-01",
		ShopURL:    "https://shop.example.com",
		Tag:        domain.Tag315Pro,
	}

	tagRoleMap := map[string]string{"315pro": "123456"}
	expectedMessage := "https://shop.example.com\n<@&123456>\n截止時間: 2026-04-01"

	tc.On("CreateThread", mock.Anything, "ch-1", "test order").
		Return("thread-id", nil)
	ma.On("AddRoleMembersToThread", mock.Anything, "thread-id", "123456", expectedMessage).
		Return(nil).Maybe()
	repo.On("CreateOrder", mock.Anything, order).
		Return(nil)

	uc := usecase.NewCreateOrder(repo, tc, ma, tagRoleMap)

	err := uc.Execute(context.Background(), "ch-1", order)

	require.NoError(t, err)
}

func TestCreateOrder_Success_OnlyTitle(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)

	order := domain.Order{
		ThreadName: "minimal order",
	}

	tc.On("CreateThread", mock.Anything, "ch-1", "minimal order").
		Return("thread-id", nil)
	repo.On("CreateOrder", mock.Anything, order).
		Return(nil)

	uc := usecase.NewCreateOrder(repo, tc, nil, nil)

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

	tc.On("CreateThread", mock.Anything, "ch-1", "partial order").
		Return("thread-id", nil)
	tc.On("SendThreadMessage", mock.Anything, "thread-id", expectedMessage).Return(nil)
	repo.On("CreateOrder", mock.Anything, order).
		Return(nil)

	uc := usecase.NewCreateOrder(repo, tc, nil, nil)

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

	tc.On("CreateThread", mock.Anything, "ch-1", "url order").
		Return("thread-id", nil)
	tc.On("SendThreadMessage", mock.Anything, "thread-id", expectedMessage).Return(nil)
	repo.On("CreateOrder", mock.Anything, order).
		Return(nil)

	uc := usecase.NewCreateOrder(repo, tc, nil, nil)

	err := uc.Execute(context.Background(), "ch-1", order)

	require.NoError(t, err)
}

func TestCreateOrder_Success_TagOnly(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)
	ma := mocks.NewMemberAdder(t)

	order := domain.Order{
		ThreadName: "tag order",
		Tag:        domain.TagGakumas,
	}

	tagRoleMap := map[string]string{"学マス": "789012"}
	expectedMessage := "<@&789012>"

	tc.On("CreateThread", mock.Anything, "ch-1", "tag order").
		Return("thread-id", nil)
	ma.On("AddRoleMembersToThread", mock.Anything, "thread-id", "789012", expectedMessage).
		Return(nil).Maybe()
	repo.On("CreateOrder", mock.Anything, order).
		Return(nil)

	uc := usecase.NewCreateOrder(repo, tc, ma, tagRoleMap)

	err := uc.Execute(context.Background(), "ch-1", order)

	require.NoError(t, err)
}

func TestCreateOrder_Success_TagWithoutRoleID(t *testing.T) {
	repo := mocks.NewOrderRepository(t)
	tc := mocks.NewThreadCreator(t)

	order := domain.Order{
		ThreadName: "fallback order",
		Tag:        domain.Tag283Pro,
	}

	tagRoleMap := map[string]string{}
	expectedMessage := "@283pro"

	tc.On("CreateThread", mock.Anything, "ch-1", "fallback order").
		Return("thread-id", nil)
	tc.On("SendThreadMessage", mock.Anything, "thread-id", expectedMessage).Return(nil)
	repo.On("CreateOrder", mock.Anything, order).
		Return(nil)

	uc := usecase.NewCreateOrder(repo, tc, nil, tagRoleMap)

	err := uc.Execute(context.Background(), "ch-1", order)

	require.NoError(t, err)
}
