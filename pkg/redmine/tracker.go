package redmine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Tracker struct {
	ID                    int      `json:"id,omitempty"`
	Name                  string   `json:"name,omitempty"`
	DefaultStatus         Resource `json:"default_status,omitempty"`
	Description           string   `json:"description,omitempty"`
	EnabledStandardFields []string `json:"enabled_standard_fields,omitempty"`
}

type TrackersResponse struct {
	Trackers []Tracker `json:"trackers"`
}

// ListTrackers retrieves the list of all trackers
func (c *Client) ListTrackers(ctx context.Context) (*TrackersResponse, error) {
	endpoint := c.baseURL + "/trackers.json"

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

	var result TrackersResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
