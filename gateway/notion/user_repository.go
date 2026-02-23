package notion

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"

	"github.com/xgnid-tw/gx5/domain"
)

var currencyColumnMap = map[domain.Currency]string{
	domain.CurrencyTWD: "台幣",
	domain.CurrencyJPY: "日幣",
}

// Repository implements port.UserRepository using the Notion API.
type Repository struct {
	db         notionapi.DatabaseService
	userDBID   notionapi.DatabaseID
	othersDBID notionapi.DatabaseID
}

func NewRepository(token string, userDBID string, othersDBID string) *Repository {
	c := notionapi.NewClient(notionapi.Token(token))

	return &Repository{
		db:         c.Database,
		userDBID:   notionapi.DatabaseID(userDBID),
		othersDBID: notionapi.DatabaseID(othersDBID),
	}
}

func (r *Repository) GetUsers(ctx context.Context) ([]*domain.User, error) {
	result, err := r.db.Query(ctx, r.userDBID, &notionapi.DatabaseQueryRequest{})
	if err != nil {
		return nil, fmt.Errorf("notion database query failed: %w", err)
	}

	users := make([]*domain.User, 0, len(result.Results))

	for _, v := range result.Results {
		discordID, ok := getTitleContent(v.Properties["discord_id"])
		if !ok {
			return nil, fmt.Errorf("failed to fetch discord column")
		}

		name, ok := getRichTextContent(v.Properties["name"])
		if !ok {
			return nil, fmt.Errorf("failed to fetch name column")
		}

		notionID, ok := getRichTextContent(v.Properties["notion_id"])
		if !ok {
			return nil, fmt.Errorf("failed to fetch notion column")
		}

		currency, ok := getSelectContent(v.Properties["currency"])
		if !ok {
			return nil, fmt.Errorf("failed to fetch currency column")
		}

		users = append(users, &domain.User{
			DiscordID: discordID,
			Name:      name,
			NotionID:  notionID,
			Currency:  domain.Currency(currency),
		})
	}

	return users, nil
}

func (r *Repository) GetUnpaidAmount(
	ctx context.Context, userDatabaseID string, currency domain.Currency,
) (float64, error) {
	col, ok := currencyColumnMap[currency]
	if !ok {
		return 0, fmt.Errorf("unsupported currency: %s", currency)
	}

	filter := &notionapi.PropertyFilter{
		Property: "付款狀況",
		Select:   &notionapi.SelectFilterCondition{Equals: "尚未付款"},
	}

	res, err := r.db.Query(ctx, notionapi.DatabaseID(userDatabaseID), &notionapi.DatabaseQueryRequest{
		Filter: filter,
	})
	if err != nil {
		return 0, fmt.Errorf("notion database query failed: %w", err)
	}

	total := float64(0)

	for _, p := range res.Results {
		amount, ok := getNumberContent(p.Properties[col])
		if !ok {
			return 0, fmt.Errorf("failed to fetch amount column")
		}

		total += amount
	}

	return total, nil
}

func (r *Repository) GetOthersUnpaidAmount(
	ctx context.Context, buyerName string, currency domain.Currency,
) (float64, error) {
	col, ok := currencyColumnMap[currency]
	if !ok {
		return 0, fmt.Errorf("unsupported currency: %s", currency)
	}

	filter := notionapi.AndCompoundFilter{
		notionapi.PropertyFilter{
			Property: "購買人",
			Select:   &notionapi.SelectFilterCondition{Equals: buyerName},
		},
		notionapi.PropertyFilter{
			Property: "付款狀況",
			Select:   &notionapi.SelectFilterCondition{Equals: "尚未付款"},
		},
	}

	res, err := r.db.Query(ctx, r.othersDBID, &notionapi.DatabaseQueryRequest{
		Filter: filter,
	})
	if err != nil {
		return 0, fmt.Errorf("notion database query failed: %w", err)
	}

	total := float64(0)

	for _, p := range res.Results {
		amount, ok := getNumberContent(p.Properties[col])
		if !ok {
			return 0, fmt.Errorf("failed to fetch amount column")
		}

		total += amount
	}

	return total, nil
}

func getTitleContent(p notionapi.Property) (string, bool) {
	tp, ok := p.(*notionapi.TitleProperty)
	if ok && len(tp.Title) > 0 {
		return tp.Title[0].Text.Content, true
	}

	return "", false
}

func getRichTextContent(p notionapi.Property) (string, bool) {
	rtp, ok := p.(*notionapi.RichTextProperty)
	if ok && len(rtp.RichText) > 0 {
		return rtp.RichText[0].Text.Content, true
	}

	return "", false
}

func getSelectContent(p notionapi.Property) (string, bool) {
	sp, ok := p.(*notionapi.SelectProperty)
	if ok && sp.Select.Name != "" {
		return sp.Select.Name, true
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
