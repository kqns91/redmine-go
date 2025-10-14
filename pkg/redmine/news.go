package redmine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type News struct {
	ID          int      `json:"id,omitempty"`
	Project     Resource `json:"project,omitempty"`
	Author      Resource `json:"author,omitempty"`
	Title       string   `json:"title,omitempty"`
	Summary     string   `json:"summary,omitempty"`
	Description string   `json:"description,omitempty"`
	CreatedOn   string   `json:"created_on,omitempty"`
}

type NewsResponse struct {
	News       []News `json:"news"`
	TotalCount int    `json:"total_count,omitempty"`
	Offset     int    `json:"offset,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

type ListNewsOptions struct {
	Limit  int
	Offset int
}

// ListNews retrieves all news across all projects with pagination
func (c *Client) ListNews(ctx context.Context, opts *ListNewsOptions) (*NewsResponse, error) {
	endpoint := c.baseURL + "/news.json"

	if opts != nil {
		params := url.Values{}
		if opts.Limit > 0 {
			params.Add("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			params.Add("offset", strconv.Itoa(opts.Offset))
		}
		if len(params) > 0 {
			endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
		}
	}

	resp, err := c.do(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result NewsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ListProjectNews retrieves all news from a specific project with pagination
func (c *Client) ListProjectNews(ctx context.Context, projectIDOrIdentifier string, opts *ListNewsOptions) (*NewsResponse, error) {
	endpoint := fmt.Sprintf("%s/projects/%s/news.json", c.baseURL, projectIDOrIdentifier)

	if opts != nil {
		params := url.Values{}
		if opts.Limit > 0 {
			params.Add("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			params.Add("offset", strconv.Itoa(opts.Offset))
		}
		if len(params) > 0 {
			endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
		}
	}

	resp, err := c.do(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result NewsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
