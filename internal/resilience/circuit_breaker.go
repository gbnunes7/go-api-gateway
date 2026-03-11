package resilience

import (
	"api-gateway-go/internal/contract"
	"api-gateway-go/internal/dto"
	"api-gateway-go/internal/observability/logger"
	"context"
	"time"

	"github.com/sony/gobreaker"
)

type UsersProviderWithCircuitBreaker struct {
	userProvider   contract.UsersProvider
	circuitBreaker *gobreaker.CircuitBreaker
}

type OrdersProviderWithCircuitBreaker struct {
	ordersProvider contract.OrdersProvider
	circuitBreaker *gobreaker.CircuitBreaker
}

type BillingsProviderWithCircuitBreaker struct {
	billingsProvider contract.BillingsProvider
	circuitBreaker   *gobreaker.CircuitBreaker
}

func (p *UsersProviderWithCircuitBreaker) GetUsers(ctx context.Context) ([]dto.User, error) {
	v, err := p.circuitBreaker.Execute(func() (interface{}, error) {
		return p.userProvider.GetUsers(ctx)
	})
	if err != nil {
		return nil, err
	}
	return v.([]dto.User), nil
}

func (p *BillingsProviderWithCircuitBreaker) GetBillings(ctx context.Context) ([]dto.Billing, error) {
	v, err := p.circuitBreaker.Execute(func() (interface{}, error) {
		return p.billingsProvider.GetBillings(ctx)
	})
	if err != nil {
		return nil, err
	}
	return v.([]dto.Billing), nil
}
func (p *OrdersProviderWithCircuitBreaker) GetOrders(ctx context.Context) ([]dto.Order, error) {

	v, err := p.circuitBreaker.Execute(func() (interface{}, error) {
		return p.ordersProvider.GetOrders(ctx)
	})

	if err != nil {
		return nil, err
	}

	return v.([]dto.Order), nil
}

func NewUsersProviderWithCircuitBreaker(userProvider contract.UsersProvider, circuitBreaker *gobreaker.CircuitBreaker) contract.UsersProvider {
	return &UsersProviderWithCircuitBreaker{userProvider: userProvider, circuitBreaker: circuitBreaker}
}

func NewOrdersProviderWithCircuitBreaker(ordersProvider contract.OrdersProvider, circuitBreaker *gobreaker.CircuitBreaker) contract.OrdersProvider {
	return &OrdersProviderWithCircuitBreaker{ordersProvider: ordersProvider, circuitBreaker: circuitBreaker}
}

func NewBillingsProviderWithCircuitBreaker(billingsProvider contract.BillingsProvider, circuitBreaker *gobreaker.CircuitBreaker) contract.BillingsProvider {
	return &BillingsProviderWithCircuitBreaker{billingsProvider: billingsProvider, circuitBreaker: circuitBreaker}
}

func DefaultCircuitBreakerSettings(name string, logger logger.Logger) gobreaker.Settings {
	return gobreaker.Settings{
		Name:        name,
		MaxRequests: 2,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.WithContext(context.Background()).Info().
				Str("name", name).
				Str("from", from.String()).
				Str("to", to.String()).
				Msg("Circuit breaker state changed")
		},
	}
}

func NewCircuitBreaker(name string, logger logger.Logger) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(DefaultCircuitBreakerSettings(name, logger))
}
