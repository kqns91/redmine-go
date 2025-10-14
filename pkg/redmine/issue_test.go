package redmine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListIssues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		response := IssuesResponse{
			Issues: []Issue{
				{ID: 1, Subject: "Test Issue 1"},
				{ID: 2, Subject: "Test Issue 2"},
			},
			TotalCount: 2,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	result, err := client.ListIssues(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListIssues failed: %v", err)
	}

	if len(result.Issues) != 2 {
		t.Errorf("Expected 2 issues, got %d", len(result.Issues))
	}
}

func TestListIssuesWithFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("project_id") != "1" {
			t.Errorf("Expected project_id=1, got %s", query.Get("project_id"))
		}
		if query.Get("status_id") != "open" {
			t.Errorf("Expected status_id=open, got %s", query.Get("status_id"))
		}

		response := IssuesResponse{Issues: []Issue{}}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	opts := &ListIssuesOptions{
		ProjectID: 1,
		StatusID:  "open",
	}
	_, err := client.ListIssues(context.Background(), opts)
	if err != nil {
		t.Fatalf("ListIssues with filters failed: %v", err)
	}
}

func TestCreateIssue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var req IssueRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Issue.Subject != "New Issue" {
			t.Errorf("Expected subject 'New Issue', got '%s'", req.Issue.Subject)
		}

		w.WriteHeader(http.StatusCreated)
		response := IssueResponse{
			Issue: Issue{
				ID:      123,
				Subject: req.Issue.Subject,
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	issue := Issue{
		Subject: "New Issue",
	}
	result, err := client.CreateIssue(context.Background(), issue)
	if err != nil {
		t.Fatalf("CreateIssue failed: %v", err)
	}

	if result.Issue.ID != 123 {
		t.Errorf("Expected issue ID 123, got %d", result.Issue.ID)
	}
}

func TestAddWatcher(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/issues/1/watchers.json" {
			t.Errorf("Expected path /issues/1/watchers.json, got %s", r.URL.Path)
		}

		var req WatcherRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.UserID != 5 {
			t.Errorf("Expected user_id 5, got %d", req.UserID)
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	err := client.AddWatcher(context.Background(), 1, 5)
	if err != nil {
		t.Fatalf("AddWatcher failed: %v", err)
	}
}
