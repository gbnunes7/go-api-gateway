package container

import (
	"net/http"

	"api-gateway-go/internal/clients/billings"
	"api-gateway-go/internal/clients/orders"
	"api-gateway-go/internal/clients/users"
	"api-gateway-go/internal/config"
	"api-gateway-go/internal/handler"
	"api-gateway-go/internal/router"
	"api-gateway-go/internal/usecase"
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

	httpClient := &http.Client{}

	usersClient := users.NewClient(cfg.UsersURL, httpClient)
	ordersClient := orders.NewClient(cfg.OrdersURL, httpClient)
	billingsClient := billings.NewClient(cfg.BillingsURL, httpClient)

	dashboardUsecase := usecase.NewGetDashboardUsecase(usersClient, ordersClient, billingsClient)
	dashboardHandler := handler.NewDashboardHandler(dashboardUsecase)

	mux := http.NewServeMux()
	router.BindRoutes(mux, dashboardHandler)

	return &Container{Mux: mux}
}
