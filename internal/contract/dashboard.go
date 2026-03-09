package contract

import (
	"api-gateway-go/internal/dto"
	"context"
)

type UsersProvider interface {
	GetUsers(ctx context.Context) ([]dto.User, error)
}

type OrdersProvider interface {
	GetOrders(ctx context.Context) ([]dto.Order, error)
}

type BillingsProvider interface {
	GetBillings(ctx context.Context) ([]dto.Billing, error)
}
