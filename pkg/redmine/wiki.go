package redmine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type WikiPage struct {
	Title       string       `json:"title,omitempty"`
	Text        string       `json:"text,omitempty"`
	Version     int          `json:"version,omitempty"`
	Author      Resource     `json:"author,omitempty"`
	Comments    string       `json:"comments,omitempty"`
	CreatedOn   string       `json:"created_on,omitempty"`
	UpdatedOn   string       `json:"updated_on,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type WikiPagesResponse struct {
	WikiPages []WikiPageIndex `json:"wiki_pages"`
}

type WikiPageIndex struct {
	Title     string   `json:"title,omitempty"`
	Version   int      `json:"version,omitempty"`
	CreatedOn string   `json:"created_on,omitempty"`
	UpdatedOn string   `json:"updated_on,omitempty"`
	Parent    Resource `json:"parent,omitempty"`
}

type WikiPageResponse struct {
	WikiPage WikiPage `json:"wiki_page"`
}

type WikiPageRequest struct {
	WikiPage WikiPageUpdate `json:"wiki_page"`
}

type WikiPageUpdate struct {
	Text     string `json:"text,omitempty"`
	Comments string `json:"comments,omitempty"`
	Version  int    `json:"version,omitempty"`
}

// ListWikiPages retrieves wiki pages index for a project
func (c *Client) ListWikiPages(projectIDOrIdentifier string) (*WikiPagesResponse, error) {
	endpoint := fmt.Sprintf("%s/projects/%s/wiki/index.json", c.baseURL, projectIDOrIdentifier)

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

	var result WikiPagesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

type GetWikiPageOptions struct {
	Include string
	Version int
}

// GetWikiPage retrieves a specific wiki page
func (c *Client) GetWikiPage(projectIDOrIdentifier string, pageName string, opts *GetWikiPageOptions) (*WikiPageResponse, error) {
	endpoint := fmt.Sprintf("%s/projects/%s/wiki/%s.json", c.baseURL, projectIDOrIdentifier, pageName)

	if opts != nil {
		params := url.Values{}
		if opts.Include != "" {
			params.Add("include", opts.Include)
		}
		if opts.Version > 0 {
			endpoint = fmt.Sprintf("%s/projects/%s/wiki/%s/%d.json", c.baseURL, projectIDOrIdentifier, pageName, opts.Version)
		} else if len(params) > 0 {
			endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
		}
	}

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

	var result WikiPageResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateOrUpdateWikiPage creates or updates a wiki page
func (c *Client) CreateOrUpdateWikiPage(projectIDOrIdentifier string, pageName string, page WikiPageUpdate) error {
	endpoint := fmt.Sprintf("%s/projects/%s/wiki/%s.json", c.baseURL, projectIDOrIdentifier, pageName)

	reqBody := WikiPageRequest{WikiPage: page}
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

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create/update wiki page: %s", string(body))
	}

	return nil
}

// DeleteWikiPage deletes a wiki page
func (c *Client) DeleteWikiPage(projectIDOrIdentifier string, pageName string) error {
	endpoint := fmt.Sprintf("%s/projects/%s/wiki/%s.json", c.baseURL, projectIDOrIdentifier, pageName)

	resp, err := c.do(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete wiki page: %s", string(body))
	}

	return nil
}
