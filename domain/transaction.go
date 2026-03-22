package domain

// Transaction represents a buy record to be inserted into a member's TBL-002.
type Transaction struct {
	ItemName  string  // 品項: thread title
	JPYAmount float64 // 日幣: user-input JPY amount
	TWDAmount float64 // 台幣: JPY × exchange rate
	DatabaseID string // target member's TBL-002 database ID (from TBL-001 notion_id)
}
