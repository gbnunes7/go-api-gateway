package utils

import (
	"context"
	"errors"
	"net"
	"net/http"
)

func StatusAndMessageFromError(err error) (status int, message string) {
	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusGatewayTimeout, "Request timed out"
	}
	if errors.Is(err, context.Canceled) {
		return 499, "Request canceled"
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return http.StatusBadGateway, "Upstream unavailable"
	}
	return http.StatusInternalServerError, "Internal server error"
}
