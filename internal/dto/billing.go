package dto

type Billing struct {
	ID          string `json:"id"`
	PaymentType string `json:"payment_type"`
	PaidValue   int    `json:"paid_value"`

	OrderID string `json:"order_id"`
}
