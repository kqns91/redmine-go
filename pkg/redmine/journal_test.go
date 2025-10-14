package redmine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShowJournal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/journals/123.json" {
			t.Errorf("Expected path /journals/123.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-Api-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		response := JournalResponse{
			Journal: Journal{
				ID:        123,
				Notes:     "Test journal note",
				CreatedOn: "2024-01-01T00:00:00Z",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ShowJournal(context.Background(), 123)
	if err != nil {
		t.Fatalf("ShowJournal failed: %v", err)
	}

	if result.Journal.ID != 123 {
		t.Errorf("Expected ID 123, got %d", result.Journal.ID)
	}
	if result.Journal.Notes != "Test journal note" {
		t.Errorf("Expected notes 'Test journal note', got %s", result.Journal.Notes)
	}
}
