package redmine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Query struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	IsPublic  bool   `json:"is_public,omitempty"`
	ProjectID int    `json:"project_id,omitempty"`
}

type QueriesResponse struct {
	Queries    []Query `json:"queries"`
	TotalCount int     `json:"total_count,omitempty"`
	Offset     int     `json:"offset,omitempty"`
	Limit      int     `json:"limit,omitempty"`
}

// ListQueries retrieves all custom queries visible by the user
func (c *Client) ListQueries(ctx context.Context) (*QueriesResponse, error) {
	endpoint := c.baseURL + "/queries.json"

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

	var result QueriesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
