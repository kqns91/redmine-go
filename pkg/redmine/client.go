package redmine

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	baseURL string
	apiKey  string

	HTTPClient *http.Client
}

func New(endpoint string, apiKey string) *Client {
	return &Client{
		baseURL:    strings.TrimSuffix(endpoint, "/"),
		apiKey:     apiKey,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) do(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Redmine-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request: %w", err)
	}

	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		//nolint:errcheck
		defer resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		return nil, fmt.Errorf("failed to request: %s", b)
	}

	return resp, nil
}
