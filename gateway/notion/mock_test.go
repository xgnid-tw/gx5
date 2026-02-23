package notion

import (
	"context"

	"github.com/jomei/notionapi"
)

type mockDatabaseService struct {
	queryFn func(
		ctx context.Context, id notionapi.DatabaseID, req *notionapi.DatabaseQueryRequest,
	) (*notionapi.DatabaseQueryResponse, error)
}

func (m *mockDatabaseService) Query(
	ctx context.Context, id notionapi.DatabaseID, req *notionapi.DatabaseQueryRequest,
) (*notionapi.DatabaseQueryResponse, error) {
	return m.queryFn(ctx, id, req)
}

func (m *mockDatabaseService) Create(
	context.Context, *notionapi.DatabaseCreateRequest,
) (*notionapi.Database, error) {
	panic("not implemented")
}

func (m *mockDatabaseService) Get(
	context.Context, notionapi.DatabaseID,
) (*notionapi.Database, error) {
	panic("not implemented")
}

func (m *mockDatabaseService) Update(
	context.Context, notionapi.DatabaseID, *notionapi.DatabaseUpdateRequest,
) (*notionapi.Database, error) {
	panic("not implemented")
}

func newTestRepository(db notionapi.DatabaseService, userDBID string) *Repository {
	return &Repository{
		db:         db,
		userDBID:   notionapi.DatabaseID(userDBID),
		othersDBID: notionapi.DatabaseID("others-db"),
	}
}

func makeUserPage(discordID, name, notionID, currency string) notionapi.Page {
	return notionapi.Page{
		Properties: notionapi.Properties{
			"discord_id": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{Text: &notionapi.Text{Content: discordID}},
				},
			},
			"name": &notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{Text: &notionapi.Text{Content: name}},
				},
			},
			"notion_id": &notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{Text: &notionapi.Text{Content: notionID}},
				},
			},
			"currency": &notionapi.SelectProperty{
				Select: notionapi.Option{Name: currency},
			},
		},
	}
}

func makeAmountPage(column string, amount float64) notionapi.Page {
	return notionapi.Page{
		Properties: notionapi.Properties{
			column: &notionapi.NumberProperty{Number: amount},
		},
	}
}
