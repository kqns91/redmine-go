package redmine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// TimeEntry represents a time entry returned by GET endpoints
type TimeEntry struct {
	ID           int           `json:"id,omitempty"`
	Project      Resource      `json:"project,omitempty"`
	Issue        Resource      `json:"issue,omitempty"`
	User         Resource      `json:"user,omitempty"`
	Activity     Resource      `json:"activity,omitempty"`
	Hours        float64       `json:"hours,omitempty"`
	Comments     string        `json:"comments,omitempty"`
	SpentOn      string        `json:"spent_on,omitempty"`
	CreatedOn    string        `json:"created_on,omitempty"`
	UpdatedOn    string        `json:"updated_on,omitempty"`
	CustomFields []CustomField `json:"custom_fields,omitempty"`
}

// TimeEntryCreateRequest represents the request body for creating a new time entry
type TimeEntryCreateRequest struct {
	IssueID      int           `json:"issue_id,omitempty"`
	ProjectID    int           `json:"project_id,omitempty"`
	SpentOn      string        `json:"spent_on,omitempty"`
	Hours        float64       `json:"hours"`
	ActivityID   int           `json:"activity_id,omitempty"`
	Comments     string        `json:"comments,omitempty"`
	UserID       int           `json:"user_id,omitempty"`
	CustomFields []CustomField `json:"custom_fields,omitempty"`
}

// TimeEntryUpdateRequest represents the request body for updating an existing time entry
type TimeEntryUpdateRequest struct {
	IssueID      int           `json:"issue_id,omitempty"`
	ProjectID    int           `json:"project_id,omitempty"`
	SpentOn      string        `json:"spent_on,omitempty"`
	Hours        float64       `json:"hours,omitempty"`
	ActivityID   int           `json:"activity_id,omitempty"`
	Comments     string        `json:"comments,omitempty"`
	UserID       int           `json:"user_id,omitempty"`
	CustomFields []CustomField `json:"custom_fields,omitempty"`
}

type TimeEntriesResponse struct {
	TimeEntries []TimeEntry `json:"time_entries"`
	TotalCount  int         `json:"total_count,omitempty"`
	Offset      int         `json:"offset,omitempty"`
	Limit       int         `json:"limit,omitempty"`
}

type TimeEntryResponse struct {
	TimeEntry TimeEntry `json:"time_entry"`
}

type TimeEntryCreateRequestWrapper struct {
	TimeEntry TimeEntryCreateRequest `json:"time_entry"`
}

type TimeEntryUpdateRequestWrapper struct {
	TimeEntry TimeEntryUpdateRequest `json:"time_entry"`
}

type ListTimeEntriesOptions struct {
	UserID    int
	ProjectID string
	SpentOn   string
	From      string
	To        string
	Limit     int
	Offset    int
}

// ListTimeEntries retrieves a list of time entries
func (c *Client) ListTimeEntries(ctx context.Context, opts *ListTimeEntriesOptions) (*TimeEntriesResponse, error) {
	endpoint := c.baseURL + "/time_entries.json"

	if opts != nil {
		params := url.Values{}
		if opts.UserID > 0 {
			params.Add("user_id", strconv.Itoa(opts.UserID))
		}
		if opts.ProjectID != "" {
			params.Add("project_id", opts.ProjectID)
		}
		if opts.SpentOn != "" {
			params.Add("spent_on", opts.SpentOn)
		}
		if opts.From != "" {
			params.Add("from", opts.From)
		}
		if opts.To != "" {
			params.Add("to", opts.To)
		}
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

	var result TimeEntriesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ShowTimeEntry retrieves a single time entry by ID
func (c *Client) ShowTimeEntry(ctx context.Context, id int) (*TimeEntryResponse, error) {
	endpoint := fmt.Sprintf("%s/time_entries/%d.json", c.baseURL, id)

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

	var result TimeEntryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateTimeEntry creates a new time entry
func (c *Client) CreateTimeEntry(ctx context.Context, req TimeEntryCreateRequest) (*TimeEntryResponse, error) {
	endpoint := c.baseURL + "/time_entries.json"

	reqBody := TimeEntryCreateRequestWrapper{TimeEntry: req}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.do(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create time entry: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result TimeEntryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateTimeEntry updates an existing time entry
func (c *Client) UpdateTimeEntry(ctx context.Context, id int, req TimeEntryUpdateRequest) error {
	endpoint := fmt.Sprintf("%s/time_entries/%d.json", c.baseURL, id)

	reqBody := TimeEntryUpdateRequestWrapper{TimeEntry: req}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.do(ctx, http.MethodPut, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update time entry: %s", string(body))
	}

	return nil
}

// DeleteTimeEntry deletes a time entry
func (c *Client) DeleteTimeEntry(ctx context.Context, id int) error {
	endpoint := fmt.Sprintf("%s/time_entries/%d.json", c.baseURL, id)

	resp, err := c.do(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete time entry: %s", string(body))
	}

	return nil
}
