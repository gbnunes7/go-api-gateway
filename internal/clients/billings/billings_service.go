package billings

import (
	"api-gateway-go/internal/dto"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"api-gateway-go/internal/constants"
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

func (c *Client) GetBillings(ctx context.Context) ([]dto.Billing, error) {
	u, err := url.JoinPath(c.baseURL, "/billings")

	if err != nil {
		return nil, fmt.Errorf("billings client: build url error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	if err != nil {
		return nil, fmt.Errorf("billings client: create request error: %w", err)
	}

	traceID, _ := ctx.Value(constants.TraceIDKey).(string)

	if traceID != "" {
		req.Header.Set(string(constants.TraceIDKey), traceID)
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("billings client: do request error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("billings client: status code error: %d", resp.StatusCode)
	}

	var billings []dto.Billing

	if err := json.NewDecoder(resp.Body).Decode(&billings); err != nil {
		return nil, fmt.Errorf("billings client: decode response error: %w", err)
	}

	return billings, nil
}
