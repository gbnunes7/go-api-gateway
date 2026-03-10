package dto

type DashboardResponse struct {
	Users  []UserWithOrders  `json:"users"`
	Errors map[string]string `json:"errors,omitempty"`
}

type UserWithOrders struct {
	ID     string             `json:"id"`
	Name   string             `json:"name"`
	Email  string             `json:"email"`
	Orders []OrderWithBilling `json:"orders"`
}

type OrderWithBilling struct {
	ID         string   `json:"id"`
	TotalPrice int      `json:"totalPrice"`
	CreatedAt  string   `json:"createdAt"`
	Billing    *Billing `json:"billing"`
}
