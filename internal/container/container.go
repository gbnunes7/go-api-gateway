package container

import (
	"net/http"

	"api-gateway-go/internal/clients/billings"
	"api-gateway-go/internal/clients/orders"
	"api-gateway-go/internal/clients/users"
	"api-gateway-go/internal/config"
	"api-gateway-go/internal/handler"
	"api-gateway-go/internal/observability/logger"
	"api-gateway-go/internal/router"
	"api-gateway-go/internal/usecase"

	"api-gateway-go/internal/resilience"
)

type Config struct {
	UsersURL    string
	OrdersURL   string
	BillingsURL string
}

type Container struct {
	Mux *http.ServeMux
}

func New(cfg *config.Config) *Container {
	cfg = config.NewConfig()

	logger := logger.New()

	httpClient := &http.Client{}

	usersClient := users.NewClient(cfg.UsersURL, httpClient)
	ordersClient := orders.NewClient(cfg.OrdersURL, httpClient)
	billingsClient := billings.NewClient(cfg.BillingsURL, httpClient)

	cbUsers := resilience.NewCircuitBreaker("users", logger)
	cbOrders := resilience.NewCircuitBreaker("orders", logger)
	cbBillings := resilience.NewCircuitBreaker("billings", logger)

	usersProvider := resilience.NewUsersProviderWithCircuitBreaker(usersClient, cbUsers)
	ordersProvider := resilience.NewOrdersProviderWithCircuitBreaker(ordersClient, cbOrders)
	billingsProvider := resilience.NewBillingsProviderWithCircuitBreaker(billingsClient, cbBillings)

	dashboardUsecase := usecase.NewGetDashboardUsecase(usersProvider, ordersProvider, billingsProvider, logger)
	dashboardHandler := handler.NewDashboardHandler(dashboardUsecase, logger)

	mux := http.NewServeMux()
	router.BindRoutes(mux, dashboardHandler)

	return &Container{Mux: mux}
}
