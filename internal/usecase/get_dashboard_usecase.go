package usecase

import (
	"api-gateway-go/internal/contract"
	"api-gateway-go/internal/dto"
	"context"
)

type GetDashboardUsecase struct {
	usersProvider    contract.UsersProvider
	ordersProvider   contract.OrdersProvider
	billingsProvider contract.BillingsProvider
}

func NewGetDashboardUsecase(
	users contract.UsersProvider,
	orders contract.OrdersProvider,
	billings contract.BillingsProvider,
) *GetDashboardUsecase {
	return &GetDashboardUsecase{
		usersProvider:    users,
		ordersProvider:   orders,
		billingsProvider: billings,
	}
}

func (u *GetDashboardUsecase) Execute(ctx context.Context) (dto.DashboardResponse, error) {
	users, err := u.usersProvider.GetUsers(ctx)
	if err != nil {
		return dto.DashboardResponse{}, err
	}
	orders, err := u.ordersProvider.GetOrders(ctx)
	if err != nil {
		return dto.DashboardResponse{}, err
	}
	billings, err := u.billingsProvider.GetBillings(ctx)
	if err != nil {
		return dto.DashboardResponse{}, err
	}

	billingByOrderID := make(map[string]dto.Billing)
	for _, b := range billings {
		billingByOrderID[b.OrderID] = b
	}

	var out []dto.UserWithOrders
	for _, u := range users {
		userRow := dto.UserWithOrders{ID: u.ID, Name: u.Name, Email: u.Email}
		for _, o := range orders {
			if o.UserID != u.ID {
				continue
			}
			ob := dto.OrderWithBilling{ID: o.ID, TotalPrice: o.TotalPrice, CreatedAt: o.CreatedAt}
			if b, ok := billingByOrderID[o.ID]; ok {
				ob.Billing = &b
			}
			userRow.Orders = append(userRow.Orders, ob)
		}
		out = append(out, userRow)
	}
	return dto.DashboardResponse{Users: out}, nil
}
