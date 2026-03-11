package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"api-gateway-go/internal/constants"
	"api-gateway-go/internal/observability/logger"
	"api-gateway-go/internal/observability/metrics"
	"api-gateway-go/internal/usecase"
	"api-gateway-go/internal/utils"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type DashboardHandler struct {
	getDashboard *usecase.GetDashboardUsecase
	logger       logger.Logger
}

func NewDashboardHandler(getDashboard *usecase.GetDashboardUsecase, logger logger.Logger) *DashboardHandler {
	return &DashboardHandler{getDashboard: getDashboard, logger: logger}
}

func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx := r.Context()
	tracer := otel.Tracer("api-gateway")
	ctx, span := tracer.Start(ctx, "GET /dashboard")
	defer span.End()

	start := time.Now()
	h.logger.WithContext(ctx).Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("[GET /dashboard] request started")

	resp, err := h.getDashboard.Execute(ctx)
	durationMs := int(time.Since(start).Milliseconds())

	if err != nil {
		status, message := utils.StatusAndMessageFromError(err)
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		h.logger.WithContext(ctx).Error().
			Err(err).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", status).
			Int("duration_ms", durationMs).
			Msg("[GET /dashboard] request failed")
		WriteError(w, status, message)
		return
	}

	statusStr := strconv.Itoa(http.StatusOK)
	metrics.HTTPRequestsTotal.WithLabelValues(r.Method, r.URL.Path, statusStr).Inc()
	metrics.HTTPRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(float64(durationMs) / 1000.0)

	span.SetAttributes(
		attribute.Int("http.status_code", http.StatusOK),
		attribute.Int64("duration_ms", int64(durationMs)),
	)

	h.logger.WithContext(ctx).Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Int("status", http.StatusOK).
		Int("duration_ms", durationMs).
		Msg("[GET /dashboard] request completed")

	traceID, _ := ctx.Value(constants.TraceIDKey).(string)
	w.Header().Set("Content-Type", "application/json")
	if traceID != "" {
		w.Header().Set("X-Trace-ID", traceID)
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
