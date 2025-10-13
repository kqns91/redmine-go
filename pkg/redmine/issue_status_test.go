package redmine

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListIssueStatuses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/issue_statuses.json" {
			t.Errorf("Expected path /issue_statuses.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-API-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		response := IssueStatusesResponse{
			IssueStatuses: []IssueStatus{
				{ID: 1, Name: "New", IsClosed: false},
				{ID: 2, Name: "Closed", IsClosed: true},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ListIssueStatuses()
	if err != nil {
		t.Fatalf("ListIssueStatuses failed: %v", err)
	}

	if len(result.IssueStatuses) != 2 {
		t.Errorf("Expected 2 issue statuses, got %d", len(result.IssueStatuses))
	}
	if result.IssueStatuses[0].Name != "New" {
		t.Errorf("Expected name 'New', got %s", result.IssueStatuses[0].Name)
	}
	if result.IssueStatuses[1].IsClosed != true {
		t.Errorf("Expected IsClosed to be true for 'Closed' status")
	}
}
