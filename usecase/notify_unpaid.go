package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/xgnid-tw/gx5/domain"
	"github.com/xgnid-tw/gx5/port"
)

const (
	twdNotificationThreshold = 2000
	jpyNotificationThreshold = 8000
	overdueMonths            = 3
)

var notificationAmountLimit = map[domain.Currency]float64{
	domain.CurrencyTWD: twdNotificationThreshold,
	domain.CurrencyJPY: jpyNotificationThreshold,
}

type NotifyUnpaid struct {
	repo     port.UserRepository
	notifier port.Notifier
	now      func() time.Time
}

func NewNotifyUnpaid(
	repo port.UserRepository, notifier port.Notifier,
) *NotifyUnpaid {
	return &NotifyUnpaid{
		repo: repo, notifier: notifier,
		now: time.Now,
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
	summary, err := uc.repo.GetUnpaidSummary(ctx, u.NotionID, u.Currency)
	if err != nil {
		return false, fmt.Errorf("get unpaid summary for %s: %w", u.Name, err)
	}

	if summary.TotalAmount > notificationAmountLimit[u.Currency] {
		return true, nil
	}

	if summary.TotalAmount > 0 && !summary.OldestRecordAt.IsZero() {
		cutoff := uc.now().AddDate(0, -overdueMonths, 0)
		if summary.OldestRecordAt.Before(cutoff) {
			return true, nil
		}
	}

	return false, nil
}
