package redmine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type File struct {
	ID          int      `json:"id,omitempty"`
	Filename    string   `json:"filename,omitempty"`
	Filesize    int      `json:"filesize,omitempty"`
	ContentType string   `json:"content_type,omitempty"`
	Description string   `json:"description,omitempty"`
	ContentURL  string   `json:"content_url,omitempty"`
	Author      Resource `json:"author,omitempty"`
	Version     Resource `json:"version,omitempty"`
	Digest      string   `json:"digest,omitempty"`
	Downloads   int      `json:"downloads,omitempty"`
	CreatedOn   string   `json:"created_on,omitempty"`
}

type FilesResponse struct {
	Files []File `json:"files"`
}

type FileUploadRequest struct {
	File FileUpload `json:"file"`
}

type FileUpload struct {
	Token       string `json:"token"`
	VersionID   int    `json:"version_id,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Description string `json:"description,omitempty"`
}

// ListFiles retrieves files available for a specific project
func (c *Client) ListFiles(ctx context.Context, projectIDOrIdentifier string) (*FilesResponse, error) {
	endpoint := fmt.Sprintf("%s/projects/%s/files.json", c.baseURL, projectIDOrIdentifier)

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

	var result FilesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UploadFile uploads a file to a specific project
func (c *Client) UploadFile(ctx context.Context, projectIDOrIdentifier string, fileUpload FileUpload) error {
	endpoint := fmt.Sprintf("%s/projects/%s/files.json", c.baseURL, projectIDOrIdentifier)

	reqBody := FileUploadRequest{File: fileUpload}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.do(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload file: %s", string(body))
	}

	return nil
}
