package handler

import (
	"context"
	"encoding/json"
	"errors"
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
		if errors.Is(err, context.DeadlineExceeded) {
			w.WriteHeader(http.StatusGatewayTimeout)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Trace-ID", traceID)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
