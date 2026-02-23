package domain

type Currency string

const (
	CurrencyTWD Currency = "TWD"
	CurrencyJPY Currency = "JPY"
)

type User struct {
	DiscordID string
	Name      string
	NotionID  string
	Currency  Currency
}
