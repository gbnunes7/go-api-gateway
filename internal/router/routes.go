package router

import (
	"encoding/json"
	"net/http"

	"api-gateway-go/internal/handler"
)

func BindRoutes(mux *http.ServeMux, dashboard *handler.DashboardHandler) {
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "Server is healthy"})
	}))
	mux.Handle("/dashboard", dashboard)
}
