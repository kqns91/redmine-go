package redmine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Enumeration struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	IsDefault bool   `json:"is_default,omitempty"`
}

type EnumerationsResponse struct {
	Enumerations []Enumeration
}

// UnmarshalJSON implements custom unmarshaling for EnumerationsResponse
// to handle different JSON field names (issue_priorities, time_entry_activities, document_categories)
func (e *EnumerationsResponse) UnmarshalJSON(data []byte) error {
	var temp map[string]json.RawMessage
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Try each possible field name
	for _, key := range []string{"issue_priorities", "time_entry_activities", "document_categories"} {
		if raw, ok := temp[key]; ok {
			return json.Unmarshal(raw, &e.Enumerations)
		}
	}

	return errors.New("no recognized enumeration field found")
}

// ListIssuePriorities retrieves the list of issue priorities
func (c *Client) ListIssuePriorities(ctx context.Context) (*EnumerationsResponse, error) {
	endpoint := c.baseURL + "/enumerations/issue_priorities.json"

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

	var result EnumerationsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ListTimeEntryActivities retrieves the list of time entry activities
func (c *Client) ListTimeEntryActivities(ctx context.Context) (*EnumerationsResponse, error) {
	endpoint := c.baseURL + "/enumerations/time_entry_activities.json"

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

	var result EnumerationsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ListDocumentCategories retrieves the list of document categories
func (c *Client) ListDocumentCategories(ctx context.Context) (*EnumerationsResponse, error) {
	endpoint := c.baseURL + "/enumerations/document_categories.json"

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

	var result EnumerationsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
