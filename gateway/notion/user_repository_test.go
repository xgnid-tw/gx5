package notion

import (
	"context"
	"errors"
	"testing"

	"github.com/jomei/notionapi"
	"github.com/stretchr/testify/require"

	"github.com/xgnid-tw/gx5/domain"
)

// --- GetUsers tests ---

func TestGetUsers_Success(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(_ context.Context, _ notionapi.DatabaseID, _ *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			return &notionapi.DatabaseQueryResponse{
				Results: []notionapi.Page{makeUserPage("111", "Alice", "abc", "TWD")},
			}, nil
		},
	}

	repo := newTestRepository(db, "user-db")
	users, err := repo.GetUsers(context.Background())

	require.NoError(t, err)
	require.Len(t, users, 1)
	require.Equal(t, "111", users[0].DiscordID)
	require.Equal(t, "Alice", users[0].Name)
	require.Equal(t, "abc", users[0].NotionID)
	require.Equal(t, domain.CurrencyTWD, users[0].Currency)
}

func TestGetUsers_MultipleUsers(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(_ context.Context, _ notionapi.DatabaseID, _ *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			return &notionapi.DatabaseQueryResponse{
				Results: []notionapi.Page{
					makeUserPage("111", "Alice", "abc", "TWD"),
					makeUserPage("222", "Bob", "def", "JPY"),
				},
			}, nil
		},
	}

	repo := newTestRepository(db, "user-db")
	users, err := repo.GetUsers(context.Background())

	require.NoError(t, err)
	require.Len(t, users, 2)
	require.Equal(t, "Alice", users[0].Name)
	require.Equal(t, domain.CurrencyTWD, users[0].Currency)
	require.Equal(t, "Bob", users[1].Name)
	require.Equal(t, domain.CurrencyJPY, users[1].Currency)
}

func TestGetUsers_QueryError(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(context.Context, notionapi.DatabaseID, *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			return nil, errors.New("api down")
		},
	}

	repo := newTestRepository(db, "user-db")
	_, err := repo.GetUsers(context.Background())

	require.Error(t, err)
	require.ErrorContains(t, err, "notion database query failed")
}

func TestGetUsers_MissingDiscordID(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(context.Context, notionapi.DatabaseID, *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			page := notionapi.Page{
				Properties: notionapi.Properties{
					"discord_id": &notionapi.TitleProperty{Title: []notionapi.RichText{}},
					"name":       &notionapi.RichTextProperty{RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: "Alice"}}}},
					"notion_id":  &notionapi.RichTextProperty{RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: "abc"}}}},
					"currency":   &notionapi.SelectProperty{Select: notionapi.Option{Name: "TWD"}},
				},
			}
			return &notionapi.DatabaseQueryResponse{Results: []notionapi.Page{page}}, nil
		},
	}

	repo := newTestRepository(db, "user-db")
	_, err := repo.GetUsers(context.Background())

	require.Error(t, err)
	require.ErrorContains(t, err, "failed to fetch discord column")
}

func TestGetUsers_MissingName(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(context.Context, notionapi.DatabaseID, *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			page := notionapi.Page{
				Properties: notionapi.Properties{
					"discord_id": &notionapi.TitleProperty{Title: []notionapi.RichText{{Text: &notionapi.Text{Content: "111"}}}},
					"name":       &notionapi.RichTextProperty{RichText: []notionapi.RichText{}},
					"notion_id":  &notionapi.RichTextProperty{RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: "abc"}}}},
					"currency":   &notionapi.SelectProperty{Select: notionapi.Option{Name: "TWD"}},
				},
			}
			return &notionapi.DatabaseQueryResponse{Results: []notionapi.Page{page}}, nil
		},
	}

	repo := newTestRepository(db, "user-db")
	_, err := repo.GetUsers(context.Background())

	require.Error(t, err)
	require.ErrorContains(t, err, "failed to fetch name column")
}

func TestGetUsers_MissingNotionID(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(context.Context, notionapi.DatabaseID, *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			page := notionapi.Page{
				Properties: notionapi.Properties{
					"discord_id": &notionapi.TitleProperty{Title: []notionapi.RichText{{Text: &notionapi.Text{Content: "111"}}}},
					"name":       &notionapi.RichTextProperty{RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: "Alice"}}}},
					"notion_id":  &notionapi.RichTextProperty{RichText: []notionapi.RichText{}},
					"currency":   &notionapi.SelectProperty{Select: notionapi.Option{Name: "TWD"}},
				},
			}
			return &notionapi.DatabaseQueryResponse{Results: []notionapi.Page{page}}, nil
		},
	}

	repo := newTestRepository(db, "user-db")
	_, err := repo.GetUsers(context.Background())

	require.Error(t, err)
	require.ErrorContains(t, err, "failed to fetch notion column")
}

func TestGetUsers_MissingCurrency(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(context.Context, notionapi.DatabaseID, *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			page := notionapi.Page{
				Properties: notionapi.Properties{
					"discord_id": &notionapi.TitleProperty{Title: []notionapi.RichText{{Text: &notionapi.Text{Content: "111"}}}},
					"name":       &notionapi.RichTextProperty{RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: "Alice"}}}},
					"notion_id":  &notionapi.RichTextProperty{RichText: []notionapi.RichText{{Text: &notionapi.Text{Content: "abc"}}}},
					"currency":   &notionapi.SelectProperty{},
				},
			}
			return &notionapi.DatabaseQueryResponse{Results: []notionapi.Page{page}}, nil
		},
	}

	repo := newTestRepository(db, "user-db")
	_, err := repo.GetUsers(context.Background())

	require.Error(t, err)
	require.ErrorContains(t, err, "failed to fetch currency column")
}

// --- GetUnpaidAmount tests ---

func TestGetUnpaidAmount_Success(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(_ context.Context, id notionapi.DatabaseID, _ *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			require.Equal(t, notionapi.DatabaseID("tx-db"), id)
			return &notionapi.DatabaseQueryResponse{
				Results: []notionapi.Page{
					makeAmountPage("台幣", 1000),
					makeAmountPage("台幣", 500),
				},
			}, nil
		},
	}

	repo := newTestRepository(db, "user-db")
	total, err := repo.GetUnpaidAmount(context.Background(), "tx-db", domain.CurrencyTWD)

	require.NoError(t, err)
	require.Equal(t, 1500.0, total)
}

func TestGetUnpaidAmount_JPY(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(_ context.Context, _ notionapi.DatabaseID, _ *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			return &notionapi.DatabaseQueryResponse{
				Results: []notionapi.Page{
					makeAmountPage("日幣", 3000),
					makeAmountPage("日幣", 5000),
				},
			}, nil
		},
	}

	repo := newTestRepository(db, "user-db")
	total, err := repo.GetUnpaidAmount(context.Background(), "tx-db", domain.CurrencyJPY)

	require.NoError(t, err)
	require.Equal(t, 8000.0, total)
}

func TestGetUnpaidAmount_EmptyResult(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(context.Context, notionapi.DatabaseID, *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			return &notionapi.DatabaseQueryResponse{Results: []notionapi.Page{}}, nil
		},
	}

	repo := newTestRepository(db, "user-db")
	total, err := repo.GetUnpaidAmount(context.Background(), "tx-db", domain.CurrencyTWD)

	require.NoError(t, err)
	require.Equal(t, 0.0, total)
}

func TestGetUnpaidAmount_QueryError(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(context.Context, notionapi.DatabaseID, *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			return nil, errors.New("api down")
		},
	}

	repo := newTestRepository(db, "user-db")
	_, err := repo.GetUnpaidAmount(context.Background(), "tx-db", domain.CurrencyTWD)

	require.Error(t, err)
	require.ErrorContains(t, err, "notion database query failed")
}

func TestGetUnpaidAmount_UnsupportedCurrency(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(context.Context, notionapi.DatabaseID, *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			return &notionapi.DatabaseQueryResponse{Results: []notionapi.Page{}}, nil
		},
	}

	repo := newTestRepository(db, "user-db")
	_, err := repo.GetUnpaidAmount(context.Background(), "tx-db", domain.Currency("USD"))

	require.Error(t, err)
	require.ErrorContains(t, err, "unsupported currency")
}

func TestGetUnpaidAmount_MissingAmountColumn(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(context.Context, notionapi.DatabaseID, *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			page := notionapi.Page{
				Properties: notionapi.Properties{},
			}
			return &notionapi.DatabaseQueryResponse{Results: []notionapi.Page{page}}, nil
		},
	}

	repo := newTestRepository(db, "user-db")
	_, err := repo.GetUnpaidAmount(context.Background(), "tx-db", domain.CurrencyTWD)

	require.Error(t, err)
	require.ErrorContains(t, err, "failed to fetch amount column")
}

// --- GetOthersUnpaidAmount tests ---

func TestGetOthersUnpaidAmount_Success(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(_ context.Context, id notionapi.DatabaseID, _ *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			require.Equal(t, notionapi.DatabaseID("others-db"), id)
			return &notionapi.DatabaseQueryResponse{
				Results: []notionapi.Page{
					makeAmountPage("台幣", 800),
					makeAmountPage("台幣", 200),
				},
			}, nil
		},
	}

	repo := newTestRepository(db, "user-db")
	total, err := repo.GetOthersUnpaidAmount(context.Background(), "Alice", domain.CurrencyTWD)

	require.NoError(t, err)
	require.Equal(t, 1000.0, total)
}

func TestGetOthersUnpaidAmount_JPY(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(_ context.Context, _ notionapi.DatabaseID, _ *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			return &notionapi.DatabaseQueryResponse{
				Results: []notionapi.Page{
					makeAmountPage("日幣", 4000),
					makeAmountPage("日幣", 5000),
				},
			}, nil
		},
	}

	repo := newTestRepository(db, "user-db")
	total, err := repo.GetOthersUnpaidAmount(context.Background(), "Bob", domain.CurrencyJPY)

	require.NoError(t, err)
	require.Equal(t, 9000.0, total)
}

func TestGetOthersUnpaidAmount_EmptyResult(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(context.Context, notionapi.DatabaseID, *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			return &notionapi.DatabaseQueryResponse{Results: []notionapi.Page{}}, nil
		},
	}

	repo := newTestRepository(db, "user-db")
	total, err := repo.GetOthersUnpaidAmount(context.Background(), "Alice", domain.CurrencyTWD)

	require.NoError(t, err)
	require.Equal(t, 0.0, total)
}

func TestGetOthersUnpaidAmount_QueryError(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(context.Context, notionapi.DatabaseID, *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			return nil, errors.New("api down")
		},
	}

	repo := newTestRepository(db, "user-db")
	_, err := repo.GetOthersUnpaidAmount(context.Background(), "Alice", domain.CurrencyTWD)

	require.Error(t, err)
	require.ErrorContains(t, err, "notion database query failed")
}

func TestGetOthersUnpaidAmount_UnsupportedCurrency(t *testing.T) {
	db := &mockDatabaseService{
		queryFn: func(context.Context, notionapi.DatabaseID, *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
			return &notionapi.DatabaseQueryResponse{Results: []notionapi.Page{}}, nil
		},
	}

	repo := newTestRepository(db, "user-db")
	_, err := repo.GetOthersUnpaidAmount(context.Background(), "Alice", domain.Currency("USD"))

	require.Error(t, err)
	require.ErrorContains(t, err, "unsupported currency")
}
