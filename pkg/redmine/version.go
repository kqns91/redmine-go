package redmine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Version struct {
	ID             int      `json:"id,omitempty"`
	Project        Resource `json:"project,omitempty"`
	Name           string   `json:"name,omitempty"`
	Description    string   `json:"description,omitempty"`
	Status         string   `json:"status,omitempty"`
	DueDate        string   `json:"due_date,omitempty"`
	Sharing        string   `json:"sharing,omitempty"`
	WikiPageTitle  string   `json:"wiki_page_title,omitempty"`
	EstimatedHours float64  `json:"estimated_hours,omitempty"`
	SpentHours     float64  `json:"spent_hours,omitempty"`
	CreatedOn      string   `json:"created_on,omitempty"`
	UpdatedOn      string   `json:"updated_on,omitempty"`
}

type VersionsResponse struct {
	Versions   []Version `json:"versions"`
	TotalCount int       `json:"total_count,omitempty"`
}

type VersionResponse struct {
	Version Version `json:"version"`
}

type VersionRequest struct {
	Version Version `json:"version"`
}

// ListVersions retrieves versions for a specific project
func (c *Client) ListVersions(ctx context.Context, projectIDOrIdentifier string) (*VersionsResponse, error) {
	endpoint := fmt.Sprintf("%s/projects/%s/versions.json", c.baseURL, projectIDOrIdentifier)

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

	var result VersionsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ShowVersion retrieves a specific version by ID
func (c *Client) ShowVersion(ctx context.Context, id int) (*VersionResponse, error) {
	endpoint := fmt.Sprintf("%s/versions/%d.json", c.baseURL, id)

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

	var result VersionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateVersion creates a new version for a project
func (c *Client) CreateVersion(ctx context.Context, projectIDOrIdentifier string, version Version) (*VersionResponse, error) {
	endpoint := fmt.Sprintf("%s/projects/%s/versions.json", c.baseURL, projectIDOrIdentifier)

	reqBody := VersionRequest{Version: version}
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
		return nil, fmt.Errorf("failed to create version: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result VersionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateVersion updates an existing version
func (c *Client) UpdateVersion(ctx context.Context, id int, version Version) error {
	endpoint := fmt.Sprintf("%s/versions/%d.json", c.baseURL, id)

	reqBody := VersionRequest{Version: version}
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
		return fmt.Errorf("failed to update version: %s", string(body))
	}

	return nil
}

// DeleteVersion deletes a version
func (c *Client) DeleteVersion(ctx context.Context, id int) error {
	endpoint := fmt.Sprintf("%s/versions/%d.json", c.baseURL, id)

	resp, err := c.do(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete version: %s", string(body))
	}

	return nil
}
