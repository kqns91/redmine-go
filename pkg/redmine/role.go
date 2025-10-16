package redmine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Role struct {
	ID                    int      `json:"id,omitempty"`
	Name                  string   `json:"name,omitempty"`
	Assignable            bool     `json:"assignable,omitempty"`
	IssuesVisibility      string   `json:"issues_visibility,omitempty"`
	UsersVisibility       string   `json:"users_visibility,omitempty"`
	TimeEntriesVisibility string   `json:"time_entries_visibility,omitempty"`
	Permissions           []string `json:"permissions,omitempty"`
}

type RolesResponse struct {
	Roles []Role `json:"roles"`
}

type RoleResponse struct {
	Role Role `json:"role"`
}

// ListRoles retrieves the list of all roles
func (c *Client) ListRoles(ctx context.Context) (*RolesResponse, error) {
	endpoint := c.baseURL + "/roles.json"

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

	var result RolesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ShowRole retrieves permissions for a specific role
func (c *Client) ShowRole(ctx context.Context, id int) (*RoleResponse, error) {
	endpoint := fmt.Sprintf("%s/roles/%d.json", c.baseURL, id)

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

	var result RoleResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
