package handler

import (
	"context"
	"net/http"

	"api-gateway-go/internal/constants"
	"api-gateway-go/internal/utils"
)

func RequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := utils.GetOrGenerateTraceID(r)
		ctx := context.WithValue(r.Context(), constants.TraceIDKey, traceID)
		ctx, cancel := context.WithTimeout(ctx, constants.Timeout)
		defer cancel()
		if delay := r.URL.Query().Get("delay"); delay != "" {
			ctx = context.WithValue(ctx, constants.DelayKey, delay)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
