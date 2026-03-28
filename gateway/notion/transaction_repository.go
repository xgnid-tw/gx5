package notion

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"

	"github.com/xgnid-tw/gx5/domain"
)

// TransactionRepository implements port.TransactionRepository using the Notion API.
type TransactionRepository struct {
	page notionapi.PageService
}

func NewTransactionRepository(page notionapi.PageService) *TransactionRepository {
	return &TransactionRepository{page: page}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, tx domain.Transaction) error {
	req := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: notionapi.DatabaseID(tx.DatabaseID),
		},
		Properties: notionapi.Properties{
			"品項": notionapi.TitleProperty{
				Type: notionapi.PropertyTypeTitle,
				Title: []notionapi.RichText{
					{Type: notionapi.ObjectTypeText, Text: &notionapi.Text{Content: tx.ItemName}},
				},
			},
			"日幣": notionapi.NumberProperty{
				Type:   notionapi.PropertyTypeNumber,
				Number: tx.JPYAmount,
			},
			"台幣": notionapi.NumberProperty{
				Type:   notionapi.PropertyTypeNumber,
				Number: tx.TWDAmount,
			},
			"付款狀況": notionapi.SelectProperty{
				Type:   notionapi.PropertyTypeSelect,
				Select: notionapi.Option{Name: "尚未付款"},
			},
		},
	}

	_, err := r.page.Create(ctx, req)
	if err != nil {
		return fmt.Errorf("notion page create failed: %w", err)
	}

	return nil
}
