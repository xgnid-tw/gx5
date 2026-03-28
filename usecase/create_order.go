package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/xgnid-tw/gx5/domain"
	"github.com/xgnid-tw/gx5/port"
)

type CreateOrder struct {
	repo          port.OrderRepository
	threadCreator port.ThreadCreator
	tagRoleMap    map[string]string
}

func NewCreateOrder(
	repo port.OrderRepository, threadCreator port.ThreadCreator, tagRoleMap map[string]string,
) *CreateOrder {
	return &CreateOrder{
		repo:          repo,
		threadCreator: threadCreator,
		tagRoleMap:    tagRoleMap,
	}
}

func (uc *CreateOrder) Execute(
	ctx context.Context, channelID string, order domain.Order,
) error {
	if order.ThreadName == "" {
		return fmt.Errorf("orderTitle is required")
	}

	message := buildThreadMessage(order, uc.tagRoleMap)

	_, err := uc.threadCreator.CreateThread(ctx, channelID, order.ThreadName, message)
	if err != nil {
		return fmt.Errorf("create thread: %w", err)
	}

	err = uc.repo.CreateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("create order record: %w", err)
	}

	return nil
}

func buildThreadMessage(order domain.Order, tagRoleMap map[string]string) string {
	var lines []string

	if order.ShopURL != "" {
		lines = append(lines, order.ShopURL)
	}

	if order.Tag != "" {
		if roleID, ok := tagRoleMap[string(order.Tag)]; ok {
			lines = append(lines, fmt.Sprintf("<@&%s>", roleID))
		} else {
			lines = append(lines, fmt.Sprintf("@%s", order.Tag))
		}
	}

	if order.Deadline != "" {
		lines = append(lines, fmt.Sprintf("截止時間: %s", order.Deadline))
	}

	return strings.Join(lines, "\n")
}
