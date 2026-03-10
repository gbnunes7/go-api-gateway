package usecase

import (
	"api-gateway-go/internal/contract"
	"api-gateway-go/internal/dto"
	"context"
)

type usersResult struct {
	users []dto.User
	err   error
}

type ordersResult struct {
	orders []dto.Order
	err    error
}

type billingsResult struct {
	billings []dto.Billing
	err      error
}

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
	usersChan := make(chan usersResult, 1)
	ordersChan := make(chan ordersResult, 1)
	billingsChan := make(chan billingsResult, 1)

	go func() {
		users, err := u.usersProvider.GetUsers(ctx)
		usersChan <- usersResult{users: users, err: err}
	}()

	go func() {
		orders, err := u.ordersProvider.GetOrders(ctx)
		ordersChan <- ordersResult{orders: orders, err: err}
	}()

	go func() {
		billings, err := u.billingsProvider.GetBillings(ctx)
		billingsChan <- billingsResult{billings: billings, err: err}
	}()

	ru := <-usersChan
	ro := <-ordersChan
	rb := <-billingsChan

	if ru.err != nil {
		return dto.DashboardResponse{}, ru.err
	}
	if ro.err != nil {
		return dto.DashboardResponse{}, ro.err
	}
	if rb.err != nil {
		return dto.DashboardResponse{}, rb.err
	}

	billingByOrderID := make(map[string]dto.Billing)
	for _, b := range rb.billings {
		billingByOrderID[b.OrderID] = b
	}

	var out []dto.UserWithOrders
	for _, u := range ru.users {
		userRow := dto.UserWithOrders{ID: u.ID, Name: u.Name, Email: u.Email}
		for _, o := range ro.orders {
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
