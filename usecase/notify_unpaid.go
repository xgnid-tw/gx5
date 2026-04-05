package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/xgnid-tw/gx5/domain"
	"github.com/xgnid-tw/gx5/port"
)

const (
	twdNotificationThreshold = 2000
	jpyNotificationThreshold = 8000
)

var notificationAmountLimit = map[domain.Currency]float64{
	domain.CurrencyTWD: twdNotificationThreshold,
	domain.CurrencyJPY: jpyNotificationThreshold,
}

type NotifyUnpaid struct {
	repo       port.UserRepository
	notifier   port.Notifier
	othersDBID string
}

func NewNotifyUnpaid(
	repo port.UserRepository, notifier port.Notifier,
	othersDBID string,
) *NotifyUnpaid {
	return &NotifyUnpaid{
		repo: repo, notifier: notifier,
		othersDBID: othersDBID,
	}
}

func (uc *NotifyUnpaid) Execute(ctx context.Context, debug bool) error {
	users, err := uc.repo.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("get users: %w", err)
	}

	for _, u := range users {
		shouldNotify, err := uc.shouldNotifyUser(ctx, u)
		if err != nil {
			return err
		}

		if shouldNotify {
			err = uc.notifier.Notify(ctx, *u, debug)
			if err != nil {
				log.Printf("notify %s: %s", u.Name, err)
			}
		}
	}

	return nil
}

func (uc *NotifyUnpaid) shouldNotifyUser(
	ctx context.Context, u *domain.User,
) (bool, error) {
	if u.NotionID != uc.othersDBID {
		a, err := uc.repo.GetUnpaidAmount(ctx, u.NotionID, u.Currency)
		if err != nil {
			return false, fmt.Errorf("get unpaid amount for %s: %w", u.Name, err)
		}

		return a > notificationAmountLimit[u.Currency], nil
	}

	a, err := uc.repo.GetOthersUnpaidAmount(ctx, u.Name, u.Currency)
	if err != nil {
		return false, fmt.Errorf("get others unpaid amount for %s: %w", u.Name, err)
	}

	return a > 0, nil
}
