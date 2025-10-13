package redmine

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListIssueCategories(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/projects/test-project/issue_categories.json" {
			t.Errorf("Expected path /projects/test-project/issue_categories.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-API-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		response := IssueCategoriesResponse{
			IssueCategories: []IssueCategory{
				{ID: 1, Name: "Backend"},
				{ID: 2, Name: "Frontend"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ListIssueCategories("test-project")
	if err != nil {
		t.Fatalf("ListIssueCategories failed: %v", err)
	}

	if len(result.IssueCategories) != 2 {
		t.Errorf("Expected 2 issue categories, got %d", len(result.IssueCategories))
	}
	if result.IssueCategories[0].Name != "Backend" {
		t.Errorf("Expected name 'Backend', got %s", result.IssueCategories[0].Name)
	}
}

func TestCreateIssueCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/projects/test-project/issue_categories.json" {
			t.Errorf("Expected path /projects/test-project/issue_categories.json, got %s", r.URL.Path)
		}

		var req IssueCategoryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.IssueCategory.Name != "New Category" {
			t.Errorf("Expected name 'New Category', got %s", req.IssueCategory.Name)
		}

		w.WriteHeader(http.StatusCreated)
		response := IssueCategoryResponse{
			IssueCategory: IssueCategory{
				ID:   3,
				Name: req.IssueCategory.Name,
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	category := IssueCategory{
		Name: "New Category",
	}
	result, err := client.CreateIssueCategory("test-project", category)
	if err != nil {
		t.Fatalf("CreateIssueCategory failed: %v", err)
	}

	if result.IssueCategory.Name != "New Category" {
		t.Errorf("Expected name 'New Category', got %s", result.IssueCategory.Name)
	}
}

func TestDeleteIssueCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/issue_categories/1.json" {
			t.Errorf("Expected path /issue_categories/1.json, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	err := client.DeleteIssueCategory(1, nil)
	if err != nil {
		t.Fatalf("DeleteIssueCategory failed: %v", err)
	}
}
