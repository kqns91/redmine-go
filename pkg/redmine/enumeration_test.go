package redmine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListIssuePriorities(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/enumerations/issue_priorities.json" {
			t.Errorf("Expected path /enumerations/issue_priorities.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-Api-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		// Response should be {"issue_priorities": [...]}
		response := map[string][]Enumeration{
			"issue_priorities": {
				{ID: 1, Name: "Low"},
				{ID: 2, Name: "High"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ListIssuePriorities(context.Background())
	if err != nil {
		t.Fatalf("ListIssuePriorities failed: %v", err)
	}

	if len(result.Enumerations) != 2 {
		t.Errorf("Expected 2 enumerations, got %d", len(result.Enumerations))
	}
	if result.Enumerations[0].Name != "Low" {
		t.Errorf("Expected name 'Low', got %s", result.Enumerations[0].Name)
	}
}
