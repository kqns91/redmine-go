package redmine

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	apiKeyHeader = "X-Redmine-Api-Key" //nolint:gosec // header name, not a credential

	// maxRedirects matches the default policy of net/http.
	maxRedirects = 10
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

func (c *Client) do(ctx context.Context, method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set(apiKeyHeader, c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient().Do(req)
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

// httpClient returns a shallow copy of HTTPClient whose redirect policy never
// forwards the API key to another host or over plain HTTP.
//
// net/http strips only well-known credential headers such as Authorization on
// cross-domain redirects, and keeps custom headers like X-Redmine-API-Key.
// The policy is applied per request so that it also covers an HTTPClient
// replaced by the caller, whose own CheckRedirect still runs afterwards.
func (c *Client) httpClient() *http.Client {
	client := *c.HTTPClient
	next := client.CheckRedirect

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		orig := via[0].URL
		if req.URL.Host != orig.Host || (orig.Scheme == "https" && req.URL.Scheme != "https") {
			req.Header.Del(apiKeyHeader)
		}

		if next != nil {
			return next(req, via)
		}
		if len(via) >= maxRedirects {
			return errors.New("stopped after 10 redirects")
		}

		return nil
	}

	return &client
}
