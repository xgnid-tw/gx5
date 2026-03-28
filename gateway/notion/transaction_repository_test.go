package notion

import (
	"context"
	"errors"
	"testing"

	"github.com/jomei/notionapi"
	"github.com/stretchr/testify/require"

	"github.com/xgnid-tw/gx5/domain"
)

func TestCreateTransaction_Success(t *testing.T) {
	var capturedReq *notionapi.PageCreateRequest

	page := &mockPageService{
		createFn: func(_ context.Context, req *notionapi.PageCreateRequest) (*notionapi.Page, error) {
			capturedReq = req
			return &notionapi.Page{}, nil
		},
	}

	repo := NewTransactionRepository(page)
	tx := domain.Transaction{
		ItemName:   "Test Item",
		JPYAmount:  3000,
		TWDAmount:  720,
		DatabaseID: "target-db",
	}

	err := repo.CreateTransaction(context.Background(), tx)

	require.NoError(t, err)
	require.Equal(t, notionapi.DatabaseID("target-db"), capturedReq.Parent.DatabaseID)

	title, ok := capturedReq.Properties["品項"].(notionapi.TitleProperty)
	require.True(t, ok)
	require.Equal(t, "Test Item", title.Title[0].Text.Content)

	jpy, ok := capturedReq.Properties["日幣"].(notionapi.NumberProperty)
	require.True(t, ok)
	require.Equal(t, 3000.0, jpy.Number)

	twd, ok := capturedReq.Properties["台幣"].(notionapi.NumberProperty)
	require.True(t, ok)
	require.Equal(t, 720.0, twd.Number)

	status, ok := capturedReq.Properties["付款狀況"].(notionapi.SelectProperty)
	require.True(t, ok)
	require.Equal(t, "尚未付款", status.Select.Name)
}

func TestCreateTransaction_Error(t *testing.T) {
	page := &mockPageService{
		createFn: func(context.Context, *notionapi.PageCreateRequest) (*notionapi.Page, error) {
			return nil, errors.New("api down")
		},
	}

	repo := NewTransactionRepository(page)
	tx := domain.Transaction{
		ItemName:   "Test Item",
		JPYAmount:  3000,
		TWDAmount:  720,
		DatabaseID: "target-db",
	}

	err := repo.CreateTransaction(context.Background(), tx)

	require.Error(t, err)
	require.ErrorContains(t, err, "notion page create failed")
}
