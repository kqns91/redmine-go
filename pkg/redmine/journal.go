package redmine

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Journal struct {
	ID        int            `json:"id,omitempty"`
	User      Resource       `json:"user,omitempty"`
	Notes     string         `json:"notes,omitempty"`
	CreatedOn string         `json:"created_on,omitempty"`
	Details   []JournalDetail `json:"details,omitempty"`
}

type JournalDetail struct {
	Property string `json:"property,omitempty"`
	Name     string `json:"name,omitempty"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value,omitempty"`
}

type JournalResponse struct {
	Journal Journal `json:"journal"`
}

// ShowJournal retrieves a specific journal entry
// Note: Journals are typically accessed through issues with include=journals parameter
func (c *Client) ShowJournal(id int) (*JournalResponse, error) {
	endpoint := fmt.Sprintf("%s/journals/%d.json", c.baseURL, id)

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

	var result JournalResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
