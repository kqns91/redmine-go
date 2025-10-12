package redmine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Group struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	UserIDs []int  `json:"user_ids,omitempty"`
}

type GroupsResponse struct {
	Groups     []Group `json:"groups"`
	TotalCount int     `json:"total_count,omitempty"`
}

type GroupResponse struct {
	Group Group `json:"group"`
}

type GroupRequest struct {
	Group Group `json:"group"`
}

type ListGroupsOptions struct {
	Include string
}

// ListGroups retrieves the list of all groups (admin only)
func (c *Client) ListGroups(opts *ListGroupsOptions) (*GroupsResponse, error) {
	endpoint := fmt.Sprintf("%s/groups.json", c.baseURL)

	if opts != nil && opts.Include != "" {
		params := url.Values{}
		params.Add("include", opts.Include)
		endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
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

	var result GroupsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

type ShowGroupOptions struct {
	Include string
}

// ShowGroup retrieves group details (admin only)
func (c *Client) ShowGroup(id int, opts *ShowGroupOptions) (*GroupResponse, error) {
	endpoint := fmt.Sprintf("%s/groups/%d.json", c.baseURL, id)

	if opts != nil && opts.Include != "" {
		params := url.Values{}
		params.Add("include", opts.Include)
		endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
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

	var result GroupResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateGroup creates a new group (admin only)
func (c *Client) CreateGroup(group Group) (*GroupResponse, error) {
	endpoint := fmt.Sprintf("%s/groups.json", c.baseURL)

	reqBody := GroupRequest{Group: group}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.do(http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create group: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result GroupResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateGroup updates an existing group (admin only)
func (c *Client) UpdateGroup(id int, group Group) error {
	endpoint := fmt.Sprintf("%s/groups/%d.json", c.baseURL, id)

	reqBody := GroupRequest{Group: group}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.do(http.MethodPut, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update group: %s", string(body))
	}

	return nil
}

// DeleteGroup deletes a group (admin only)
func (c *Client) DeleteGroup(id int) error {
	endpoint := fmt.Sprintf("%s/groups/%d.json", c.baseURL, id)

	resp, err := c.do(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete group: %s", string(body))
	}

	return nil
}

type AddUserToGroupRequest struct {
	UserID int `json:"user_id"`
}

// AddUserToGroup adds a user to a group (admin only)
func (c *Client) AddUserToGroup(groupID int, userID int) error {
	endpoint := fmt.Sprintf("%s/groups/%d/users.json", c.baseURL, groupID)

	reqBody := AddUserToGroupRequest{UserID: userID}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.do(http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add user to group: %s", string(body))
	}

	return nil
}

// RemoveUserFromGroup removes a user from a group (admin only)
func (c *Client) RemoveUserFromGroup(groupID int, userID int) error {
	endpoint := fmt.Sprintf("%s/groups/%d/users/%d.json", c.baseURL, groupID, userID)

	resp, err := c.do(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to remove user from group: %s", string(body))
	}

	return nil
}
