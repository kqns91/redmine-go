package redmine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IssueRelation struct {
	ID           int    `json:"id,omitempty"`
	IssueID      int    `json:"issue_id,omitempty"`
	IssueToID    int    `json:"issue_to_id,omitempty"`
	RelationType string `json:"relation_type,omitempty"`
	Delay        int    `json:"delay,omitempty"`
}

type IssueRelationsResponse struct {
	Relations []IssueRelation `json:"relations"`
}

type IssueRelationResponse struct {
	Relation IssueRelation `json:"relation"`
}

type IssueRelationRequest struct {
	Relation IssueRelation `json:"relation"`
}

// ListIssueRelations retrieves all relations for a specific issue
func (c *Client) ListIssueRelations(issueID int) (*IssueRelationsResponse, error) {
	endpoint := fmt.Sprintf("%s/issues/%d/relations.json", c.baseURL, issueID)

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

	var result IssueRelationsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ShowIssueRelation retrieves details of a specific relation
func (c *Client) ShowIssueRelation(relationID int) (*IssueRelationResponse, error) {
	endpoint := fmt.Sprintf("%s/relations/%d.json", c.baseURL, relationID)

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

	var result IssueRelationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateIssueRelation creates a new issue relation
func (c *Client) CreateIssueRelation(issueID int, relation IssueRelation) (*IssueRelationResponse, error) {
	endpoint := fmt.Sprintf("%s/issues/%d/relations.json", c.baseURL, issueID)

	reqBody := IssueRelationRequest{Relation: relation}
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
		return nil, fmt.Errorf("failed to create issue relation: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result IssueRelationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteIssueRelation deletes a specific issue relation
func (c *Client) DeleteIssueRelation(relationID int) error {
	endpoint := fmt.Sprintf("%s/relations/%d.json", c.baseURL, relationID)

	resp, err := c.do(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete issue relation: %s", string(body))
	}

	return nil
}
