package utils_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"api-gateway-go/internal/utils"
)

type mockNetError struct {
	message string
}

func (e *mockNetError) Error() string   { return e.message }
func (e *mockNetError) Timeout() bool   { return false }
func (e *mockNetError) Temporary() bool { return false }

func TestStatusAndMessageFromError_DeadlineExceeded(t *testing.T) {
	err := context.DeadlineExceeded
	status, message := utils.StatusAndMessageFromError(err)
	if status != http.StatusGatewayTimeout {
		t.Errorf("status = %d, want %d", status, http.StatusGatewayTimeout)
	}
	if message != "Request timed out" {
		t.Errorf("message = %q, want %q", message, "Request timed out")
	}
}

func TestStatusAndMessageFromError_DeadlineExceededWrapped(t *testing.T) {
	err := fmt.Errorf("users client: %w", context.DeadlineExceeded)
	status, message := utils.StatusAndMessageFromError(err)
	if status != http.StatusGatewayTimeout {
		t.Errorf("status = %d, want %d", status, http.StatusGatewayTimeout)
	}
	if message != "Request timed out" {
		t.Errorf("message = %q, want %q", message, "Request timed out")
	}
}

func TestStatusAndMessageFromError_Canceled(t *testing.T) {
	err := context.Canceled
	status, message := utils.StatusAndMessageFromError(err)
	if status != 499 {
		t.Errorf("status = %d, want 499", status)
	}
	if message != "Request canceled" {
		t.Errorf("message = %q, want %q", message, "Request canceled")
	}
}

func TestStatusAndMessageFromError_CanceledWrapped(t *testing.T) {
	err := fmt.Errorf("request aborted: %w", context.Canceled)
	status, message := utils.StatusAndMessageFromError(err)
	if status != 499 {
		t.Errorf("status = %d, want 499", status)
	}
	if message != "Request canceled" {
		t.Errorf("message = %q, want %q", message, "Request canceled")
	}
}

func TestStatusAndMessageFromError_NetError(t *testing.T) {
	err := &mockNetError{message: "connection refused"}
	status, message := utils.StatusAndMessageFromError(err)
	if status != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", status, http.StatusBadGateway)
	}
	if message != "Upstream unavailable" {
		t.Errorf("message = %q, want %q", message, "Upstream unavailable")
	}
}

func TestStatusAndMessageFromError_NetErrorWrapped(t *testing.T) {
	err := fmt.Errorf("users client: %w", &mockNetError{message: "dial timeout"})
	status, message := utils.StatusAndMessageFromError(err)
	if status != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", status, http.StatusBadGateway)
	}
	if message != "Upstream unavailable" {
		t.Errorf("message = %q, want %q", message, "Upstream unavailable")
	}
}

func TestStatusAndMessageFromError_GenericError(t *testing.T) {
	err := errors.New("something went wrong")
	status, message := utils.StatusAndMessageFromError(err)
	if status != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", status, http.StatusInternalServerError)
	}
	if message != "Internal server error" {
		t.Errorf("message = %q, want %q", message, "Internal server error")
	}
}

func TestStatusAndMessageFromError_Nil(t *testing.T) {
	status, message := utils.StatusAndMessageFromError(nil)
	if status != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d (nil treated as generic)", status, http.StatusInternalServerError)
	}
	if message != "Internal server error" {
		t.Errorf("message = %q, want %q", message, "Internal server error")
	}
}
