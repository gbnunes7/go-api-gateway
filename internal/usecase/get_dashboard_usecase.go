package usecase

import (
	"api-gateway-go/internal/contract"
	"api-gateway-go/internal/dto"
	"api-gateway-go/internal/observability/logger"
	"api-gateway-go/internal/utils"
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
	logger           logger.Logger
}

func NewGetDashboardUsecase(
	users contract.UsersProvider,
	orders contract.OrdersProvider,
	billings contract.BillingsProvider,
	logger logger.Logger,
) *GetDashboardUsecase {
	return &GetDashboardUsecase{
		usersProvider:    users,
		ordersProvider:   orders,
		billingsProvider: billings,
		logger:           logger,
	}
}

func (u *GetDashboardUsecase) Execute(ctx context.Context) (dto.DashboardResponse, error) {
	u.logger.WithContext(ctx).Info().Msg("Executing GetDashboardUsecase")

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
		u.logger.WithContext(ctx).Error().
			Err(ru.err).
			Str("service", "users").
			Msg("Error getting users")
		return dto.DashboardResponse{}, ru.err
	}

	billingByOrderID := make(map[string]dto.Billing)
	if rb.err == nil {
		for _, b := range rb.billings {
			billingByOrderID[b.OrderID] = b
		}
	}

	var out []dto.UserWithOrders
	for _, u := range ru.users {
		userRow := dto.UserWithOrders{ID: u.ID, Name: u.Name, Email: u.Email}
		ordersList := ro.orders
		if ro.err != nil {
			ordersList = nil
		}
		for _, o := range ordersList {
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

	var errs map[string]string
	if ro.err != nil {
		if errs == nil {
			errs = make(map[string]string)
		}
		_, message := utils.StatusAndMessageFromError(ro.err)
		errs["orders"] = message
		u.logger.WithContext(ctx).Warn().
			Err(ro.err).
			Str("service", "orders").
			Msg("orders provider failed")
	}
	if rb.err != nil {
		if errs == nil {
			errs = make(map[string]string)
		}
		_, message := utils.StatusAndMessageFromError(rb.err)
		errs["billings"] = message
		u.logger.WithContext(ctx).Warn().
			Err(rb.err).
			Str("service", "billings").
			Msg("billings provider failed")
	}

	u.logger.WithContext(ctx).Info().
		Int("users_count", len(out)).
		Msg("dashboard execute completed")

	return dto.DashboardResponse{Users: out, Errors: errs}, nil
}
