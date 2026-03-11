package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"api-gateway-go/internal/dto"
	"api-gateway-go/internal/observability/logger"
	"api-gateway-go/internal/usecase"
)

type mockUsersProvider struct {
	users []dto.User
	err   error
}

type mockLogger struct {
}

func (m *mockLogger) WithContext(ctx context.Context) logger.Logger {
	return m
}

func (m *mockLogger) Info() logger.Event {
	return m
}

func (m *mockLogger) Error() logger.Event {
	return m
}

func (m *mockLogger) Warn() logger.Event {
	return m
}

func (m *mockLogger) Debug() logger.Event {
	return m
}

func (m *mockLogger) Str(key string, val string) logger.Event {
	return m
}

func (m *mockLogger) Int(key string, val int) logger.Event {
	return m
}

func (m *mockLogger) Err(err error) logger.Event {
	return m
}

func (m *mockLogger) Msg(msg string) {}

func (m *mockLogger) Msgf(format string, args ...interface{}) {}

func (m *mockUsersProvider) GetUsers(ctx context.Context) ([]dto.User, error) {
	return m.users, m.err
}

type mockOrdersProvider struct {
	orders []dto.Order
	err    error
}

func (m *mockOrdersProvider) GetOrders(ctx context.Context) ([]dto.Order, error) {
	return m.orders, m.err
}

type mockBillingsProvider struct {
	billings []dto.Billing
	err      error
}

func (m *mockBillingsProvider) GetBillings(ctx context.Context) ([]dto.Billing, error) {
	return m.billings, m.err
}

type mockUsersProviderContextAware struct{}

func (m *mockUsersProviderContextAware) GetUsers(ctx context.Context) ([]dto.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return []dto.User{}, nil
	}
}

type mockUsersProviderSlow struct {
	delay time.Duration
	users []dto.User
	err   error
}

func (m *mockUsersProviderSlow) GetUsers(ctx context.Context) ([]dto.User, error) {
	time.Sleep(m.delay)
	return m.users, m.err
}

type mockOrdersProviderSlow struct {
	delay  time.Duration
	orders []dto.Order
	err    error
}

func (m *mockOrdersProviderSlow) GetOrders(ctx context.Context) ([]dto.Order, error) {
	time.Sleep(m.delay)
	return m.orders, m.err
}

type mockBillingsProviderSlow struct {
	delay    time.Duration
	billings []dto.Billing
	err      error
}

func (m *mockBillingsProviderSlow) GetBillings(ctx context.Context) ([]dto.Billing, error) {
	time.Sleep(m.delay)
	return m.billings, m.err
}

func TestExecute_Success(t *testing.T) {
	users := []dto.User{
		{ID: "user-1", Name: "John", Email: "john@example.com", OrderID: "order-1"},
	}
	orders := []dto.Order{
		{ID: "order-1", TotalPrice: 100, CreatedAt: "2025-01-01", BillingID: "b1", UserID: "user-1"},
	}
	billings := []dto.Billing{
		{ID: "b1", PaymentType: "credit", PaidValue: 100, OrderID: "order-1"},
	}

	uc := usecase.NewGetDashboardUsecase(
		&mockUsersProvider{users: users},
		&mockOrdersProvider{orders: orders},
		&mockBillingsProvider{billings: billings},
		&mockLogger{},
	)

	ctx := context.Background()
	resp, err := uc.Execute(ctx)
	if err != nil {
		t.Fatalf("Execute() err = %v, want nil", err)
	}
	if len(resp.Users) != 1 {
		t.Fatalf("len(resp.Users) = %d, want 1", len(resp.Users))
	}
	u := resp.Users[0]
	if u.ID != "user-1" || u.Name != "John" || u.Email != "john@example.com" {
		t.Errorf("resp.Users[0] = %+v, want ID=user-1 Name=John", u)
	}
	if len(u.Orders) != 1 {
		t.Fatalf("len(resp.Users[0].Orders) = %d, want 1", len(u.Orders))
	}
	o := u.Orders[0]
	if o.ID != "order-1" || o.TotalPrice != 100 || o.CreatedAt != "2025-01-01" {
		t.Errorf("order = %+v", o)
	}
	if o.Billing == nil || o.Billing.ID != "b1" || o.Billing.PaymentType != "credit" {
		t.Errorf("order.Billing = %+v, want ID=b1 PaymentType=credit", o.Billing)
	}
}

func TestExecute_UsersError(t *testing.T) {
	wantErr := errors.New("users service unavailable")
	uc := usecase.NewGetDashboardUsecase(
		&mockUsersProvider{err: wantErr},
		&mockOrdersProvider{orders: []dto.Order{}},
		&mockBillingsProvider{billings: []dto.Billing{}},
		&mockLogger{},
	)

	ctx := context.Background()
	_, err := uc.Execute(ctx)
	if err != wantErr {
		t.Errorf("Execute() err = %v, want %v", err, wantErr)
	}
}

func TestExecute_OrdersError(t *testing.T) {
	wantErr := errors.New("orders service unavailable")
	uc := usecase.NewGetDashboardUsecase(
		&mockUsersProvider{users: []dto.User{{ID: "user-1", Name: "John", Email: "j@x.com", OrderID: ""}}},
		&mockOrdersProvider{err: wantErr},
		&mockBillingsProvider{billings: []dto.Billing{}},
		&mockLogger{},
	)

	ctx := context.Background()
	resp, err := uc.Execute(ctx)
	if err != nil {
		t.Fatalf("Execute() err = %v, want nil (graceful degradation)", err)
	}
	if len(resp.Users) != 1 {
		t.Fatalf("len(resp.Users) = %d, want 1", len(resp.Users))
	}
	if len(resp.Users[0].Orders) != 0 {
		t.Errorf("len(resp.Users[0].Orders) = %d, want 0 (orders failed)", len(resp.Users[0].Orders))
	}
	if resp.Errors == nil || resp.Errors["orders"] == "" {
		t.Errorf("resp.Errors[orders] want set, got %v", resp.Errors)
	}
}

func TestExecute_BillingsError(t *testing.T) {
	wantErr := errors.New("billings service unavailable")
	uc := usecase.NewGetDashboardUsecase(
		&mockUsersProvider{users: []dto.User{{ID: "user-1", Name: "John", Email: "j@x.com", OrderID: "order-1"}}},
		&mockOrdersProvider{orders: []dto.Order{{ID: "order-1", TotalPrice: 100, CreatedAt: "2025-01-01", BillingID: "b1", UserID: "user-1"}}},
		&mockBillingsProvider{err: wantErr},
		&mockLogger{},
	)

	ctx := context.Background()
	resp, err := uc.Execute(ctx)
	if err != nil {
		t.Fatalf("Execute() err = %v, want nil (graceful degradation)", err)
	}
	if len(resp.Users) != 1 || len(resp.Users[0].Orders) != 1 {
		t.Fatalf("want 1 user with 1 order, got %d users", len(resp.Users))
	}
	o := resp.Users[0].Orders[0]
	if o.Billing != nil {
		t.Errorf("order.Billing = %v, want nil (billings failed)", o.Billing)
	}
	if resp.Errors == nil || resp.Errors["billings"] == "" {
		t.Errorf("resp.Errors[billings] want set, got %v", resp.Errors)
	}
}

func TestExecute_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	uc := usecase.NewGetDashboardUsecase(
		&mockUsersProviderContextAware{},
		&mockOrdersProvider{orders: []dto.Order{}},
		&mockBillingsProvider{billings: []dto.Billing{}},
		&mockLogger{},
	)

	_, err := uc.Execute(ctx)
	if err == nil {
		t.Error("Execute() err = nil, want context.Canceled")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Execute() err = %v, want context.Canceled", err)
	}
}

func TestExecute_EmptyData(t *testing.T) {
	uc := usecase.NewGetDashboardUsecase(
		&mockUsersProvider{users: []dto.User{}},
		&mockOrdersProvider{orders: []dto.Order{}},
		&mockBillingsProvider{billings: []dto.Billing{}},
		&mockLogger{},
	)

	ctx := context.Background()
	resp, err := uc.Execute(ctx)
	if err != nil {
		t.Fatalf("Execute() err = %v, want nil", err)
	}
	if resp.Users != nil {
		t.Errorf("resp.Users = %v, want nil or empty slice", resp.Users)
	}
}

func TestExecute_ParallelTiming(t *testing.T) {
	t.Log("1. Criando use case com mocks lentos: users=100ms, orders=200ms, billings=150ms")
	uc := usecase.NewGetDashboardUsecase(
		&mockUsersProviderSlow{delay: 100 * time.Millisecond, users: []dto.User{}},
		&mockOrdersProviderSlow{delay: 200 * time.Millisecond, orders: []dto.Order{}},
		&mockBillingsProviderSlow{delay: 150 * time.Millisecond, billings: []dto.Billing{}},
		&mockLogger{},
	)

	t.Log("2. Context sem timeout para não cancelar os sleeps")
	ctx := context.Background()

	t.Log("3. Iniciando cronômetro e chamando Execute(ctx)")
	start := time.Now()
	_, err := uc.Execute(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Execute() err = %v", err)
	}

	t.Logf("4. Execute() retornou em %v", elapsed)
	t.Log("5. Em paralelo esperamos ~200ms (maior sleep); sequencial seria ~450ms (soma)")

	if elapsed > 300*time.Millisecond {
		t.Errorf("Tempo %v > 300ms: parece sequencial (esperado ~200ms se paralelo)", elapsed)
	} else {
		t.Logf("6. OK: tempo %v está dentro do esperado para execução paralela", elapsed)
	}
}
