package notion

import (
	"context"
	"fmt"
	"time"

	"github.com/jomei/notionapi"

	"github.com/xgnid-tw/gx5/domain"
)

// OrderRepository implements port.OrderRepository using the Notion API.
type OrderRepository struct {
	page      notionapi.PageService
	orderDBID notionapi.DatabaseID
}

func NewOrderRepository(page notionapi.PageService, orderDBID string) *OrderRepository {
	return &OrderRepository{
		page:      page,
		orderDBID: notionapi.DatabaseID(orderDBID),
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order domain.Order) error {
	props := notionapi.Properties{
		"threadName": notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{Text: &notionapi.Text{Content: order.ThreadName}},
			},
		},
	}

	if order.Deadline != "" {
		t, err := time.Parse("2006-01-02", order.Deadline)
		if err != nil {
			return fmt.Errorf("invalid deadline format: %w", err)
		}

		d := notionapi.Date(t)
		props["deadline"] = notionapi.DateProperty{
			Date: &notionapi.DateObject{Start: &d},
		}
	}

	if order.Tag != "" {
		props["tags"] = notionapi.SelectProperty{
			Select: notionapi.Option{Name: string(order.Tag)},
		}
	}

	_, err := r.page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: r.orderDBID,
		},
		Properties: props,
	})
	if err != nil {
		return fmt.Errorf("notion page create failed: %w", err)
	}

	return nil
}
