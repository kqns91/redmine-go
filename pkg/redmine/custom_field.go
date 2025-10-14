package redmine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CustomFieldDefinition struct {
	ID             int      `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	CustomizedType string   `json:"customized_type,omitempty"`
	FieldFormat    string   `json:"field_format,omitempty"`
	Regexp         string   `json:"regexp,omitempty"`
	MinLength      int      `json:"min_length,omitempty"`
	MaxLength      int      `json:"max_length,omitempty"`
	IsRequired     bool     `json:"is_required,omitempty"`
	IsFilter       bool     `json:"is_filter,omitempty"`
	Searchable     bool     `json:"searchable,omitempty"`
	Multiple       bool     `json:"multiple,omitempty"`
	DefaultValue   string   `json:"default_value,omitempty"`
	Visible        bool     `json:"visible,omitempty"`
	PossibleValues []string `json:"possible_values,omitempty"`
}

type CustomFieldsResponse struct {
	CustomFields []CustomFieldDefinition `json:"custom_fields"`
}

// ListCustomFields retrieves all custom fields definitions (requires admin privileges)
func (c *Client) ListCustomFields(ctx context.Context) (*CustomFieldsResponse, error) {
	endpoint := c.baseURL + "/custom_fields.json"

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

	var result CustomFieldsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
