package redmine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MyAccountResponse struct {
	User User `json:"user"`
}

type MyAccountRequest struct {
	User User `json:"user"`
}

// GetMyAccount retrieves current user's account details
func (c *Client) GetMyAccount(ctx context.Context) (*MyAccountResponse, error) {
	endpoint := c.baseURL + "/my/account.json"

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

	var result MyAccountResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateMyAccount updates current user's account details
func (c *Client) UpdateMyAccount(ctx context.Context, user User) error {
	endpoint := c.baseURL + "/my/account.json"

	reqBody := MyAccountRequest{User: user}
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
		return fmt.Errorf("failed to update my account: %s", string(body))
	}

	return nil
}
