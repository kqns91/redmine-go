package redmine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Membership struct {
	ID      int        `json:"id,omitempty"`
	Project Resource   `json:"project,omitempty"`
	User    Resource   `json:"user,omitempty"`
	Group   Resource   `json:"group,omitempty"`
	Roles   []Resource `json:"roles,omitempty"`
}

type MembershipsResponse struct {
	Memberships []Membership `json:"memberships"`
	TotalCount  int          `json:"total_count,omitempty"`
	Offset      int          `json:"offset,omitempty"`
	Limit       int          `json:"limit,omitempty"`
}

type MembershipResponse struct {
	Membership Membership `json:"membership"`
}

type MembershipRequest struct {
	Membership MembershipCreateUpdate `json:"membership"`
}

type MembershipCreateUpdate struct {
	UserID  int   `json:"user_id,omitempty"`
	RoleIDs []int `json:"role_ids,omitempty"`
}

// ListMemberships retrieves paginated list of project memberships
func (c *Client) ListMemberships(ctx context.Context, projectIDOrIdentifier string) (*MembershipsResponse, error) {
	endpoint := fmt.Sprintf("%s/projects/%s/memberships.json", c.baseURL, projectIDOrIdentifier)

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

	var result MembershipsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ShowMembership retrieves specific membership details
func (c *Client) ShowMembership(ctx context.Context, id int) (*MembershipResponse, error) {
	endpoint := fmt.Sprintf("%s/memberships/%d.json", c.baseURL, id)

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

	var result MembershipResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateMembership adds a new project member
func (c *Client) CreateMembership(ctx context.Context, projectIDOrIdentifier string, membership MembershipCreateUpdate) (*MembershipResponse, error) {
	endpoint := fmt.Sprintf("%s/projects/%s/memberships.json", c.baseURL, projectIDOrIdentifier)

	reqBody := MembershipRequest{Membership: membership}
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
		return nil, fmt.Errorf("failed to create membership: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result MembershipResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateMembership updates membership roles
func (c *Client) UpdateMembership(ctx context.Context, id int, roleIDs []int) error {
	endpoint := fmt.Sprintf("%s/memberships/%d.json", c.baseURL, id)

	reqBody := MembershipRequest{
		Membership: MembershipCreateUpdate{
			RoleIDs: roleIDs,
		},
	}
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
		return fmt.Errorf("failed to update membership: %s", string(body))
	}

	return nil
}

// DeleteMembership deletes a membership
func (c *Client) DeleteMembership(ctx context.Context, id int) error {
	endpoint := fmt.Sprintf("%s/memberships/%d.json", c.baseURL, id)

	resp, err := c.do(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete membership: %s", string(body))
	}

	return nil
}
