package dto

type Order struct {
	ID         string `json:"id"`
	TotalPrice int    `json:"total_price"`
	CreatedAt  string `json:"created_at"`

	BillingID string `json:"billing_id"`
	UserID    string `json:"user_id"`
}
