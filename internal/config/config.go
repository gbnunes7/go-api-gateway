package config

import "os"

type Config struct {
	UsersURL    string
	OrdersURL   string
	BillingsURL string
}

func NewConfig() *Config {
	usersURL := os.Getenv("USERS_URL")
	ordersURL := os.Getenv("ORDERS_URL")
	billingsURL := os.Getenv("BILLINGS_URL")

	return &Config{
		UsersURL:    usersURL,
		OrdersURL:   ordersURL,
		BillingsURL: billingsURL,
	}
}
