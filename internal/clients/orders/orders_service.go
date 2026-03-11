package orders

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"api-gateway-go/internal/constants"
	"api-gateway-go/internal/dto"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{baseURL: baseURL, httpClient: httpClient}
}

func (c *Client) GetOrders(ctx context.Context) ([]dto.Order, error) {
	u, err := url.JoinPath(c.baseURL, "/orders")
	if err != nil {
		return nil, fmt.Errorf("orders client: build url error: %w", err)
	}

	tracer := otel.Tracer("api-gateway")
	ctx, span := tracer.Start(ctx, "http.client.orders")
	defer span.End()
	span.SetAttributes(attribute.String("http.url", u))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, fmt.Errorf("orders client: create request error: %w", err)
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	traceID, _ := ctx.Value(constants.TraceIDKey).(string)
	if traceID != "" {
		req.Header.Set(string(constants.TraceIDKey), traceID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, fmt.Errorf("orders client: do request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		span.SetStatus(codes.Error, fmt.Sprintf("status %d", resp.StatusCode))
		return nil, fmt.Errorf("orders client: status code error: %d", resp.StatusCode)
	}

	var orders []dto.Order

	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("orders client: decode response error: %w", err)
	}

	return orders, nil
}
