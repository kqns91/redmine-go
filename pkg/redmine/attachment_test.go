package redmine

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShowAttachment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/attachments/123.json" {
			t.Errorf("Expected path /attachments/123.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-API-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		response := AttachmentResponse{
			Attachment: Attachment{
				ID:       123,
				Filename: "test.pdf",
				Filesize: 1024,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ShowAttachment(123)
	if err != nil {
		t.Fatalf("ShowAttachment failed: %v", err)
	}

	if result.Attachment.ID != 123 {
		t.Errorf("Expected ID 123, got %d", result.Attachment.ID)
	}
	if result.Attachment.Filename != "test.pdf" {
		t.Errorf("Expected filename 'test.pdf', got %s", result.Attachment.Filename)
	}
}
