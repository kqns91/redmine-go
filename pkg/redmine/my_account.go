package redmine

import (
	"bytes"
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
func (c *Client) GetMyAccount() (*MyAccountResponse, error) {
	endpoint := fmt.Sprintf("%s/my/account.json", c.baseURL)

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

	var result MyAccountResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateMyAccount updates current user's account details
func (c *Client) UpdateMyAccount(user User) error {
	endpoint := fmt.Sprintf("%s/my/account.json", c.baseURL)

	reqBody := MyAccountRequest{User: user}
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
		return fmt.Errorf("failed to update my account: %s", string(body))
	}

	return nil
}
