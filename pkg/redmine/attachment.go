package redmine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Attachment struct {
	ID          int      `json:"id,omitempty"`
	Filename    string   `json:"filename,omitempty"`
	Filesize    int      `json:"filesize,omitempty"`
	ContentType string   `json:"content_type,omitempty"`
	Description string   `json:"description,omitempty"`
	ContentURL  string   `json:"content_url,omitempty"`
	Author      Resource `json:"author,omitempty"`
	CreatedOn   string   `json:"created_on,omitempty"`
}

type AttachmentResponse struct {
	Attachment Attachment `json:"attachment"`
}

type AttachmentRequest struct {
	Attachment Attachment `json:"attachment"`
}

// ShowAttachment retrieves details of a specific attachment
func (c *Client) ShowAttachment(ctx context.Context, id int) (*AttachmentResponse, error) {
	endpoint := fmt.Sprintf("%s/attachments/%d.json", c.baseURL, id)

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

	var result AttachmentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateAttachment updates an existing attachment
func (c *Client) UpdateAttachment(ctx context.Context, id int, attachment Attachment) error {
	endpoint := fmt.Sprintf("%s/attachments/%d.json", c.baseURL, id)

	reqBody := AttachmentRequest{Attachment: attachment}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.do(ctx, http.MethodPatch, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update attachment: %s", string(body))
	}

	return nil
}

// DeleteAttachment deletes a specific attachment
func (c *Client) DeleteAttachment(ctx context.Context, id int) error {
	endpoint := fmt.Sprintf("%s/attachments/%d.json", c.baseURL, id)

	resp, err := c.do(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete attachment: %s", string(body))
	}

	return nil
}
