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

type User struct {
	ID               int           `json:"id,omitempty"`
	Login            string        `json:"login,omitempty"`
	Admin            bool          `json:"admin,omitempty"`
	Firstname        string        `json:"firstname,omitempty"`
	Lastname         string        `json:"lastname,omitempty"`
	Mail             string        `json:"mail,omitempty"`
	CreatedOn        string        `json:"created_on,omitempty"`
	UpdatedOn        string        `json:"updated_on,omitempty"`
	LastLoginOn      string        `json:"last_login_on,omitempty"`
	PasswdChangedOn  string        `json:"passwd_changed_on,omitempty"`
	TwofaScheme      string        `json:"twofa_scheme,omitempty"`
	APIKey           string        `json:"api_key,omitempty"`
	Status           int           `json:"status,omitempty"`
	AuthSourceID     int           `json:"auth_source_id,omitempty"`
	Password         string        `json:"password,omitempty"`
	MailNotification string        `json:"mail_notification,omitempty"`
	MustChangePasswd bool          `json:"must_change_passwd,omitempty"`
	GeneratePassword bool          `json:"generate_password,omitempty"`
	SendInformation  bool          `json:"send_information,omitempty"`
	CustomFields     []CustomField `json:"custom_fields,omitempty"`
}

type UsersResponse struct {
	Users      []User `json:"users"`
	TotalCount int    `json:"total_count,omitempty"`
	Offset     int    `json:"offset,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

type UserResponse struct {
	User User `json:"user"`
}

type UserRequest struct {
	User User `json:"user"`
}

type ListUsersOptions struct {
	Status  string
	Name    string
	GroupID int
	Include string
	Limit   int
	Offset  int
}

// ListUsers retrieves a list of users (requires admin privileges)
func (c *Client) ListUsers(ctx context.Context, opts *ListUsersOptions) (*UsersResponse, error) {
	endpoint := c.baseURL + "/users.json"

	if opts != nil {
		params := url.Values{}
		if opts.Status != "" {
			params.Add("status", opts.Status)
		}
		if opts.Name != "" {
			params.Add("name", opts.Name)
		}
		if opts.GroupID > 0 {
			params.Add("group_id", strconv.Itoa(opts.GroupID))
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

	var result UsersResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

type ShowUserOptions struct {
	Include string
}

// ShowUser retrieves a single user by ID
func (c *Client) ShowUser(ctx context.Context, id int, opts *ShowUserOptions) (*UserResponse, error) {
	endpoint := fmt.Sprintf("%s/users/%d.json", c.baseURL, id)

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

	var result UserResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetCurrentUser retrieves the currently authenticated user
func (c *Client) GetCurrentUser(ctx context.Context, opts *ShowUserOptions) (*UserResponse, error) {
	endpoint := c.baseURL + "/users/current.json"

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

	var result UserResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateUser creates a new user (requires admin privileges)
func (c *Client) CreateUser(ctx context.Context, user User) (*UserResponse, error) {
	endpoint := c.baseURL + "/users.json"

	reqBody := UserRequest{User: user}
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
		return nil, fmt.Errorf("failed to create user: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result UserResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateUser updates an existing user (requires admin privileges)
func (c *Client) UpdateUser(ctx context.Context, id int, user User) error {
	endpoint := fmt.Sprintf("%s/users/%d.json", c.baseURL, id)

	reqBody := UserRequest{User: user}
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
		return fmt.Errorf("failed to update user: %s", string(body))
	}

	return nil
}

// DeleteUser deletes a user (requires admin privileges)
func (c *Client) DeleteUser(ctx context.Context, id int) error {
	endpoint := fmt.Sprintf("%s/users/%d.json", c.baseURL, id)

	resp, err := c.do(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete user: %s", string(body))
	}

	return nil
}
