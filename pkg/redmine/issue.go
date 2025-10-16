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

// Issue represents an issue returned by GET endpoints
type Issue struct {
	ID              int             `json:"id,omitempty"`
	Project         Resource        `json:"project,omitempty"`
	Tracker         Resource        `json:"tracker,omitempty"`
	Status          Resource        `json:"status,omitempty"`
	Priority        Resource        `json:"priority,omitempty"`
	Author          Resource        `json:"author,omitempty"`
	AssignedTo      Resource        `json:"assigned_to,omitempty"`
	Category        Resource        `json:"category,omitempty"`
	Subject         string          `json:"subject,omitempty"`
	Description     string          `json:"description,omitempty"`
	StartDate       string          `json:"start_date,omitempty"`
	DueDate         string          `json:"due_date,omitempty"`
	DoneRatio       int             `json:"done_ratio,omitempty"`
	IsPrivate       bool            `json:"is_private,omitempty"`
	EstimatedHours  float64         `json:"estimated_hours,omitempty"`
	CustomFields    []CustomField   `json:"custom_fields,omitempty"`
	CreatedOn       string          `json:"created_on,omitempty"`
	UpdatedOn       string          `json:"updated_on,omitempty"`
	ClosedOn        string          `json:"closed_on,omitempty"`
	Journals        []Journal       `json:"journals,omitempty"`
	Children        []Issue         `json:"children,omitempty"`
	Attachments     []Attachment    `json:"attachments,omitempty"`
	Relations       []IssueRelation `json:"relations,omitempty"`
	Changesets      []Changeset     `json:"changesets,omitempty"`
	Watchers        []Watcher       `json:"watchers,omitempty"`
	AllowedStatuses []IssueStatus   `json:"allowed_statuses,omitempty"`
}

// IssueCreateRequest represents the request body for creating a new issue
type IssueCreateRequest struct {
	ProjectID      int           `json:"project_id"`
	TrackerID      int           `json:"tracker_id"`
	Subject        string        `json:"subject"`
	StatusID       int           `json:"status_id,omitempty"`
	PriorityID     int           `json:"priority_id,omitempty"`
	CategoryID     int           `json:"category_id,omitempty"`
	FixedVersionID int           `json:"fixed_version_id,omitempty"`
	AssignedToID   int           `json:"assigned_to_id,omitempty"`
	ParentIssueID  int           `json:"parent_issue_id,omitempty"`
	Description    string        `json:"description,omitempty"`
	StartDate      string        `json:"start_date,omitempty"`
	DueDate        string        `json:"due_date,omitempty"`
	DoneRatio      int           `json:"done_ratio,omitempty"`
	EstimatedHours float64       `json:"estimated_hours,omitempty"`
	IsPrivate      bool          `json:"is_private,omitempty"`
	WatcherUserIDs []int         `json:"watcher_user_ids,omitempty"`
	CustomFields   []CustomField `json:"custom_fields,omitempty"`
	Uploads        []Upload      `json:"uploads,omitempty"`
}

// IssueUpdateRequest represents the request body for updating an existing issue
type IssueUpdateRequest struct {
	ProjectID      int           `json:"project_id,omitempty"`
	TrackerID      int           `json:"tracker_id,omitempty"`
	Subject        string        `json:"subject,omitempty"`
	StatusID       int           `json:"status_id,omitempty"`
	PriorityID     int           `json:"priority_id,omitempty"`
	CategoryID     int           `json:"category_id,omitempty"`
	FixedVersionID int           `json:"fixed_version_id,omitempty"`
	AssignedToID   int           `json:"assigned_to_id,omitempty"`
	ParentIssueID  int           `json:"parent_issue_id,omitempty"`
	Description    string        `json:"description,omitempty"`
	StartDate      string        `json:"start_date,omitempty"`
	DueDate        string        `json:"due_date,omitempty"`
	DoneRatio      int           `json:"done_ratio,omitempty"`
	EstimatedHours float64       `json:"estimated_hours,omitempty"`
	IsPrivate      bool          `json:"is_private,omitempty"`
	Notes          string        `json:"notes,omitempty"`
	PrivateNotes   bool          `json:"private_notes,omitempty"`
	CustomFields   []CustomField `json:"custom_fields,omitempty"`
	Uploads        []Upload      `json:"uploads,omitempty"`
}

type IssuesResponse struct {
	Issues     []Issue `json:"issues"`
	TotalCount int     `json:"total_count,omitempty"`
	Offset     int     `json:"offset,omitempty"`
	Limit      int     `json:"limit,omitempty"`
}

type IssueResponse struct {
	Issue Issue `json:"issue"`
}

type IssueCreateRequestWrapper struct {
	Issue IssueCreateRequest `json:"issue"`
}

type IssueUpdateRequestWrapper struct {
	Issue IssueUpdateRequest `json:"issue"`
}

type ListIssuesOptions struct {
	ProjectID      int
	SubprojectID   string
	TrackerID      int
	StatusID       string
	AssignedToID   string
	PriorityID     int
	CategoryID     int
	FixedVersionID int
	IssueID        string
	ParentID       int
	Subject        string
	Description    string
	CreatedOn      string
	UpdatedOn      string
	ClosedOn       string
	StartDate      string
	DueDate        string
	EstimatedHours string
	DoneRatio      string
	Include        string
	Limit          int
	Offset         int
	Sort           string
}

// ListIssues retrieves a list of issues
func (c *Client) ListIssues(ctx context.Context, opts *ListIssuesOptions) (*IssuesResponse, error) {
	endpoint := c.baseURL + "/issues.json"

	if opts != nil {
		params := url.Values{}
		if opts.ProjectID > 0 {
			params.Add("project_id", strconv.Itoa(opts.ProjectID))
		}
		if opts.SubprojectID != "" {
			params.Add("subproject_id", opts.SubprojectID)
		}
		if opts.TrackerID > 0 {
			params.Add("tracker_id", strconv.Itoa(opts.TrackerID))
		}
		if opts.StatusID != "" {
			params.Add("status_id", opts.StatusID)
		}
		if opts.AssignedToID != "" {
			params.Add("assigned_to_id", opts.AssignedToID)
		}
		if opts.PriorityID > 0 {
			params.Add("priority_id", strconv.Itoa(opts.PriorityID))
		}
		if opts.CategoryID > 0 {
			params.Add("category_id", strconv.Itoa(opts.CategoryID))
		}
		if opts.FixedVersionID > 0 {
			params.Add("fixed_version_id", strconv.Itoa(opts.FixedVersionID))
		}
		if opts.IssueID != "" {
			params.Add("issue_id", opts.IssueID)
		}
		if opts.ParentID > 0 {
			params.Add("parent_id", strconv.Itoa(opts.ParentID))
		}
		if opts.Subject != "" {
			params.Add("subject", opts.Subject)
		}
		if opts.Description != "" {
			params.Add("description", opts.Description)
		}
		if opts.CreatedOn != "" {
			params.Add("created_on", opts.CreatedOn)
		}
		if opts.UpdatedOn != "" {
			params.Add("updated_on", opts.UpdatedOn)
		}
		if opts.ClosedOn != "" {
			params.Add("closed_on", opts.ClosedOn)
		}
		if opts.StartDate != "" {
			params.Add("start_date", opts.StartDate)
		}
		if opts.DueDate != "" {
			params.Add("due_date", opts.DueDate)
		}
		if opts.EstimatedHours != "" {
			params.Add("estimated_hours", opts.EstimatedHours)
		}
		if opts.DoneRatio != "" {
			params.Add("done_ratio", opts.DoneRatio)
		}
		if opts.Include != "" {
			params.Add("include", opts.Include)
		}
		if opts.Limit > 0 {
			params.Add("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			params.Add("offset", strconv.Itoa(opts.Offset))
		}
		if opts.Sort != "" {
			params.Add("sort", opts.Sort)
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

	var result IssuesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

type ShowIssueOptions struct {
	Include string
}

// ShowIssue retrieves a single issue by ID
func (c *Client) ShowIssue(ctx context.Context, id int, opts *ShowIssueOptions) (*IssueResponse, error) {
	endpoint := fmt.Sprintf("%s/issues/%d.json", c.baseURL, id)

	if opts != nil && opts.Include != "" {
		params := url.Values{}
		params.Add("include", opts.Include)
		endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
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

	var result IssueResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateIssue creates a new issue
func (c *Client) CreateIssue(ctx context.Context, req IssueCreateRequest) (*IssueResponse, error) {
	endpoint := c.baseURL + "/issues.json"

	reqBody := IssueCreateRequestWrapper{Issue: req}
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
		return nil, fmt.Errorf("failed to create issue: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result IssueResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateIssue updates an existing issue
func (c *Client) UpdateIssue(ctx context.Context, id int, req IssueUpdateRequest) error {
	endpoint := fmt.Sprintf("%s/issues/%d.json", c.baseURL, id)

	reqBody := IssueUpdateRequestWrapper{Issue: req}
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
		return fmt.Errorf("failed to update issue: %s", string(body))
	}

	return nil
}

// DeleteIssue deletes an issue
func (c *Client) DeleteIssue(ctx context.Context, id int) error {
	endpoint := fmt.Sprintf("%s/issues/%d.json", c.baseURL, id)

	resp, err := c.do(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete issue: %s", string(body))
	}

	return nil
}

type WatcherRequest struct {
	UserID int `json:"user_id"`
}

// AddWatcher adds a watcher to an issue
func (c *Client) AddWatcher(ctx context.Context, issueID int, userID int) error {
	endpoint := fmt.Sprintf("%s/issues/%d/watchers.json", c.baseURL, issueID)

	reqBody := WatcherRequest{UserID: userID}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.do(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add watcher: %s", string(body))
	}

	return nil
}

// RemoveWatcher removes a watcher from an issue
func (c *Client) RemoveWatcher(ctx context.Context, issueID int, userID int) error {
	endpoint := fmt.Sprintf("%s/issues/%d/watchers/%d.json", c.baseURL, issueID, userID)

	resp, err := c.do(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to remove watcher: %s", string(body))
	}

	return nil
}
