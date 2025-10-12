package redmine

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type SearchResult struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Type        string `json:"type,omitempty"`
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
	Datetime    string `json:"datetime,omitempty"`
}

type SearchResponse struct {
	Results    []SearchResult `json:"results"`
	TotalCount int            `json:"total_count,omitempty"`
	Offset     int            `json:"offset,omitempty"`
	Limit      int            `json:"limit,omitempty"`
}

type SearchOptions struct {
	Query       []string
	Scope       string
	Issues      bool
	WikiPages   bool
	Attachments bool
	Offset      int
	Limit       int
}

// Search performs a search with specified conditions
func (c *Client) Search(opts *SearchOptions) (*SearchResponse, error) {
	endpoint := fmt.Sprintf("%s/search.json", c.baseURL)

	if opts != nil {
		params := url.Values{}
		if len(opts.Query) > 0 {
			for _, q := range opts.Query {
				params.Add("q", q)
			}
		}
		if opts.Scope != "" {
			params.Add("scope", opts.Scope)
		}
		if opts.Issues {
			params.Add("issues", "1")
		}
		if opts.WikiPages {
			params.Add("wiki_pages", "1")
		}
		if opts.Attachments {
			params.Add("attachments", "1")
		}
		if opts.Offset > 0 {
			params.Add("offset", strconv.Itoa(opts.Offset))
		}
		if opts.Limit > 0 {
			params.Add("limit", strconv.Itoa(opts.Limit))
		}
		if len(params) > 0 {
			endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
		}
	}

	resp, err := c.do(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// SearchSimple performs a simple search with a query string
func (c *Client) SearchSimple(query string) (*SearchResponse, error) {
	return c.Search(&SearchOptions{
		Query: strings.Split(query, " "),
	})
}
