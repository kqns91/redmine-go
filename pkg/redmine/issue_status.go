package redmine

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IssueStatus struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IsClosed bool   `json:"is_closed,omitempty"`
}

type IssueStatusesResponse struct {
	IssueStatuses []IssueStatus `json:"issue_statuses"`
}

// ListIssueStatuses retrieves the list of all issue statuses
func (c *Client) ListIssueStatuses() (*IssueStatusesResponse, error) {
	endpoint := fmt.Sprintf("%s/issue_statuses.json", c.baseURL)

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

	var result IssueStatusesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
