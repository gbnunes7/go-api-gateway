package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"api-gateway-go/internal/constants"
	"api-gateway-go/internal/usecase"
	"api-gateway-go/internal/utils"
)

type DashboardHandler struct {
	getDashboard *usecase.GetDashboardUsecase
}

func NewDashboardHandler(getDashboard *usecase.GetDashboardUsecase) *DashboardHandler {
	return &DashboardHandler{getDashboard: getDashboard}
}

func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	traceID := utils.GetOrGenerateTraceID(r)

	ctx := context.WithValue(r.Context(), constants.TraceIDKey, traceID)

	ctx, cancel := context.WithTimeout(ctx, constants.Timeout)

	defer cancel()

	delay := r.URL.Query().Get("delay")
	if delay != "" {
		ctx = context.WithValue(ctx, constants.DelayKey, delay)
	}

	resp, err := h.getDashboard.Execute(ctx)
	if err != nil {
		status, message := utils.StatusAndMessageFromError(err)
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Trace-ID", traceID)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
