package redmine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// IssueCategory represents an issue category returned by GET endpoints
type IssueCategory struct {
	ID         int      `json:"id,omitempty"`
	Project    Resource `json:"project,omitempty"`
	Name       string   `json:"name,omitempty"`
	AssignedTo Resource `json:"assigned_to,omitempty"`
}

// IssueCategoryCreateRequest represents the request body for creating a new issue category
type IssueCategoryCreateRequest struct {
	Name         string `json:"name"`
	AssignedToID int    `json:"assigned_to_id,omitempty"`
}

// IssueCategoryUpdateRequest represents the request body for updating an existing issue category
type IssueCategoryUpdateRequest struct {
	Name         string `json:"name,omitempty"`
	AssignedToID int    `json:"assigned_to_id,omitempty"`
}

type IssueCategoriesResponse struct {
	IssueCategories []IssueCategory `json:"issue_categories"`
}

type IssueCategoryResponse struct {
	IssueCategory IssueCategory `json:"issue_category"`
}

type IssueCategoryCreateRequestWrapper struct {
	IssueCategory IssueCategoryCreateRequest `json:"issue_category"`
}

type IssueCategoryUpdateRequestWrapper struct {
	IssueCategory IssueCategoryUpdateRequest `json:"issue_category"`
}

// ListIssueCategories retrieves all issue categories for a specific project
func (c *Client) ListIssueCategories(ctx context.Context, projectIDOrIdentifier string) (*IssueCategoriesResponse, error) {
	endpoint := fmt.Sprintf("%s/projects/%s/issue_categories.json", c.baseURL, projectIDOrIdentifier)

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

	var result IssueCategoriesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ShowIssueCategory retrieves a specific issue category by ID
func (c *Client) ShowIssueCategory(ctx context.Context, id int) (*IssueCategoryResponse, error) {
	endpoint := fmt.Sprintf("%s/issue_categories/%d.json", c.baseURL, id)

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

	var result IssueCategoryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateIssueCategory creates a new issue category for a project
func (c *Client) CreateIssueCategory(ctx context.Context, projectIDOrIdentifier string, req IssueCategoryCreateRequest) (*IssueCategoryResponse, error) {
	endpoint := fmt.Sprintf("%s/projects/%s/issue_categories.json", c.baseURL, projectIDOrIdentifier)

	reqBody := IssueCategoryCreateRequestWrapper{IssueCategory: req}
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
		return nil, fmt.Errorf("failed to create issue category: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result IssueCategoryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateIssueCategory updates an existing issue category
func (c *Client) UpdateIssueCategory(ctx context.Context, id int, req IssueCategoryUpdateRequest) error {
	endpoint := fmt.Sprintf("%s/issue_categories/%d.json", c.baseURL, id)

	reqBody := IssueCategoryUpdateRequestWrapper{IssueCategory: req}
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
		return fmt.Errorf("failed to update issue category: %s", string(body))
	}

	return nil
}

type DeleteIssueCategoryOptions struct {
	ReassignToID int
}

// DeleteIssueCategory deletes an issue category
func (c *Client) DeleteIssueCategory(ctx context.Context, id int, opts *DeleteIssueCategoryOptions) error {
	endpoint := fmt.Sprintf("%s/issue_categories/%d.json", c.baseURL, id)

	if opts != nil && opts.ReassignToID > 0 {
		endpoint = fmt.Sprintf("%s?reassign_to_id=%s", endpoint, strconv.Itoa(opts.ReassignToID))
	}

	resp, err := c.do(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete issue category: %s", string(body))
	}

	return nil
}
