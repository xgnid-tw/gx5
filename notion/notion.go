package notion

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jomei/notionapi"

	"github.com/xgnid-tw/gx5/model"
)

type Notion interface {
	SendNotPaidInformation(ctx context.Context) error

	GetDiscordIDList(ctx context.Context) ([]*model.User, error)
}

const (
	// tolerance value
	discordNotificationAmountLimit = 2000
)

type notion struct {
	client   *notionapi.Client
	ch       chan model.User
	userDBID notionapi.DatabaseID
}

func NewNotion(token string, ch chan model.User, userDBID string) (Notion, error) {
	return &notion{
		client:   notionapi.NewClient(notionapi.Token(token)),
		ch:       ch,
		userDBID: notionapi.DatabaseID(userDBID),
	}, nil
}

func (s *notion) SendNotPaidInformation(ctx context.Context) error {
	// get all user
	users, err := s.GetDiscordIDList(ctx)
	if err != nil {
		log.Fatalf("%s", err)
	}

	for _, u := range users {
		// once in a day for high freq user
		if !u.High && time.Now().Day() != 1 {
			continue
		}

		amount, err := s.getUserNotPaidAmount(ctx, notionapi.DatabaseID(u.NotionID), u.High)
		if err != nil {
			return fmt.Errorf("get user not paid amount: %w", err)
		}

		if amount > discordNotificationAmountLimit {
			// send with channel
			s.ch <- *u
		}
	}

	return nil
}

func (s *notion) GetDiscordIDList(ctx context.Context) ([]*model.User, error) {
	result, err := s.client.Database.Query(ctx, s.userDBID,
		&notionapi.DatabaseQueryRequest{},
	)
	if err != nil {
		return nil, fmt.Errorf("notion database query failed: %w", err)
	}

	users := make([]*model.User, 0)

	for _, v := range result.Results {
		discordID, discordOk := getTitleContent(v.Properties["discord_id"])
		if !discordOk {
			return nil, fmt.Errorf("failed to fetch discord column")
		}

		name, nameOk := getRichTextContent(v.Properties["name"])
		if !nameOk {
			return nil, fmt.Errorf("failed to fetch name column")
		}

		notionID, notionOk := getRichTextContent(v.Properties["notion_id"])
		if !notionOk {
			return nil, fmt.Errorf("failed to fetch notion column")
		}

		high, highOk := getCheckboxContent(v.Properties["high"])
		if !highOk {
			return nil, fmt.Errorf("failed to fetch high column")
		}

		users = append(users, &model.User{
			DiscordID: discordID,
			Name:      name,
			NotionID:  notionID,
			High:      high,
		})
	}

	return users, nil
}

func (s *notion) getUserNotPaidAmount(ctx context.Context,
	userDatabaseID notionapi.DatabaseID, high bool,
) (float64, error) {
	expiredDateObj := notionapi.Date(time.Now().AddDate(0, -2, 0))

	var filter notionapi.Filter

	if high {
		filter = &notionapi.AndCompoundFilter{
			&notionapi.PropertyFilter{
				Property: "付款狀況",
				Select: &notionapi.SelectFilterCondition{
					Equals: "尚未付款",
				},
			},
		}
	} else {
		filter = &notionapi.AndCompoundFilter{
			&notionapi.PropertyFilter{
				Property: "付款狀況",
				Select: &notionapi.SelectFilterCondition{
					Equals: "尚未付款",
				},
			},
			&notionapi.TimestampFilter{
				Timestamp: notionapi.TimestampCreated,
				CreatedTime: &notionapi.DateFilterCondition{
					Before: &expiredDateObj,
				},
			},
		}
	}

	res, err := s.client.Database.Query(ctx, userDatabaseID, &notionapi.DatabaseQueryRequest{
		Filter: filter,
	})
	if err != nil {
		return 0, fmt.Errorf("notion database query failed: %w", err)
	}

	total := float64(0)

	for _, p := range res.Results {
		amount, ok := getNumberContent(p.Properties["台幣"])
		if !ok {
			return 0, fmt.Errorf("fail to fetch amount column")
		}

		total += amount
	}

	return total, nil
}

func getTitleContent(p notionapi.Property) (string, bool) {
	tp, ok := p.(*notionapi.TitleProperty)

	if ok {
		return tp.Title[0].Text.Content, true
	}

	return "", false
}

func getRichTextContent(p notionapi.Property) (string, bool) {
	rtp, ok := p.(*notionapi.RichTextProperty)
	if ok {
		return rtp.RichText[0].Text.Content, true
	}

	return "", false
}

func getNumberContent(p notionapi.Property) (float64, bool) {
	np, ok := p.(*notionapi.NumberProperty)
	if ok {
		return np.Number, true
	}

	return 0, false
}

func getCheckboxContent(p notionapi.Property) (bool, bool) {
	np, ok := p.(*notionapi.CheckboxProperty)
	if ok {
		return np.Checkbox, true
	}

	return false, false
}
