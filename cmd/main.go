package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-gateway-go/internal/config"
	"api-gateway-go/internal/container"
	"api-gateway-go/internal/observability/telemetry"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	ctx := context.Background()
	if err := telemetry.Init(ctx); err != nil {
		log.Fatal(err)
	}

	cfg := config.NewConfig()
	c := container.New(cfg)

	server := &http.Server{
		Addr:    ":8080",
		Handler: c.Mux,
	}

	go func() {
		fmt.Println("Server is running on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("http server shutdown: %v", err)
	}
	if err := telemetry.Shutdown(shutdownCtx); err != nil {
		log.Printf("telemetry shutdown: %v", err)
	}

	log.Println("shutdown complete")
}
