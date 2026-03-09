package main

import (
	"fmt"
	"net/http"

	"api-gateway-go/internal/config"
	"api-gateway-go/internal/container"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	cfg := config.NewConfig()
	container := container.New(cfg)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", container.Mux); err != nil {
		fmt.Println(err)
	}
}
