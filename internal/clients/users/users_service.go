package users

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

func (c *Client) GetUsers(ctx context.Context) ([]dto.User, error) {
	u, err := url.JoinPath(c.baseURL, "/users")
	if err != nil {
		return nil, fmt.Errorf("users client: build url error: %w", err)
	}

	if delay, _ := ctx.Value(constants.DelayKey).(string); delay != "" {
		u = u + "?delay=" + url.QueryEscape(delay)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)

	if err != nil {
		return nil, fmt.Errorf("users client: create request error: %w", err)
	}

	traceID, _ := ctx.Value(constants.TraceIDKey).(string)

	if traceID != "" {
		req.Header.Set(string(constants.TraceIDKey), traceID)
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("users client: do request error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("users client: status code error: %d", resp.StatusCode)
	}

	var users []dto.User

	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("users client: decode response error: %w", err)
	}

	return users, nil
}
