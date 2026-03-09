package utils

import (
	"api-gateway-go/internal/constants"
	"net/http"

	"github.com/google/uuid"
)

func GetOrGenerateTraceID(r *http.Request) string {
	traceID := r.Header.Get(string(constants.TraceIDKey))
	if traceID == "" {
		traceID = uuid.New().String()
	}
	return traceID
}
