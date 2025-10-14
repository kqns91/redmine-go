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

type Project struct {
	ID                     int           `json:"id,omitempty"`
	Name                   string        `json:"name,omitempty"`
	Identifier             string        `json:"identifier,omitempty"`
	Description            string        `json:"description,omitempty"`
	Homepage               string        `json:"homepage,omitempty"`
	Parent                 Resource      `json:"parent,omitempty"`
	Status                 int           `json:"status,omitempty"`
	IsPublic               bool          `json:"is_public,omitempty"`
	InheritMembers         bool          `json:"inherit_members,omitempty"`
	CustomFields           []CustomField `json:"custom_fields,omitempty"`
	ActiveNewTicketMessage string        `json:"active_new_ticket_message,omitempty"`
	EnableNewTicketMessage int           `json:"enable_new_ticket_message,omitempty"`
	NewTicketMessage       string        `json:"new_ticket_message,omitempty"`
	CreatedOn              string        `json:"created_on,omitempty"`
	UpdatedOn              string        `json:"updated_on,omitempty"`
}

type ProjectsResponse struct {
	Projects   []Project `json:"projects"`
	TotalCount int       `json:"total_count,omitempty"`
	Offset     int       `json:"offset,omitempty"`
	Limit      int       `json:"limit,omitempty"`
}

type ProjectResponse struct {
	Project Project `json:"project"`
}

type ProjectRequest struct {
	Project Project `json:"project"`
}

type ListProjectsOptions struct {
	Include string
	Limit   int
	Offset  int
}

// ListProjects retrieves a list of projects
func (c *Client) ListProjects(ctx context.Context, opts *ListProjectsOptions) (*ProjectsResponse, error) {
	endpoint := c.baseURL + "/projects.json"

	if opts != nil {
		params := url.Values{}
		if opts.Include != "" {
			params.Add("include", opts.Include)
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

	var result ProjectsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

type ShowProjectOptions struct {
	Include string
}

// ShowProject retrieves a single project by ID or identifier
func (c *Client) ShowProject(ctx context.Context, idOrIdentifier string, opts *ShowProjectOptions) (*ProjectResponse, error) {
	endpoint := fmt.Sprintf("%s/projects/%s.json", c.baseURL, idOrIdentifier)

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

	var result ProjectResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateProject creates a new project
func (c *Client) CreateProject(ctx context.Context, project Project) (*ProjectResponse, error) {
	endpoint := c.baseURL + "/projects.json"

	reqBody := ProjectRequest{Project: project}
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
		return nil, fmt.Errorf("failed to create project: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result ProjectResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateProject updates an existing project
func (c *Client) UpdateProject(ctx context.Context, idOrIdentifier string, project Project) error {
	endpoint := fmt.Sprintf("%s/projects/%s.json", c.baseURL, idOrIdentifier)

	reqBody := ProjectRequest{Project: project}
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
		return fmt.Errorf("failed to update project: %s", string(body))
	}

	return nil
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(ctx context.Context, idOrIdentifier string) error {
	endpoint := fmt.Sprintf("%s/projects/%s.json", c.baseURL, idOrIdentifier)

	resp, err := c.do(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete project: %s", string(body))
	}

	return nil
}

// ArchiveProject archives a project (available since Redmine 5.0)
func (c *Client) ArchiveProject(ctx context.Context, idOrIdentifier string) error {
	endpoint := fmt.Sprintf("%s/projects/%s/archive.json", c.baseURL, idOrIdentifier)

	resp, err := c.do(ctx, http.MethodPut, endpoint, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to archive project: %s", string(body))
	}

	return nil
}

// UnarchiveProject unarchives a project (available since Redmine 5.0)
func (c *Client) UnarchiveProject(ctx context.Context, idOrIdentifier string) error {
	endpoint := fmt.Sprintf("%s/projects/%s/unarchive.json", c.baseURL, idOrIdentifier)

	resp, err := c.do(ctx, http.MethodPut, endpoint, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to unarchive project: %s", string(body))
	}

	return nil
}
