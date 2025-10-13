package redmine

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListNews(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/news.json" {
			t.Errorf("Expected path /news.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-API-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		response := NewsResponse{
			News: []News{
				{ID: 1, Title: "News 1", Summary: "Summary 1"},
				{ID: 2, Title: "News 2", Summary: "Summary 2"},
			},
			TotalCount: 2,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ListNews(nil)
	if err != nil {
		t.Fatalf("ListNews failed: %v", err)
	}

	if len(result.News) != 2 {
		t.Errorf("Expected 2 news items, got %d", len(result.News))
	}
	if result.News[0].Title != "News 1" {
		t.Errorf("Expected title 'News 1', got %s", result.News[0].Title)
	}
}

func TestListProjectNews(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/projects/test-project/news.json" {
			t.Errorf("Expected path /projects/test-project/news.json, got %s", r.URL.Path)
		}

		response := NewsResponse{
			News: []News{
				{ID: 1, Title: "Project News", Summary: "Project summary"},
			},
			TotalCount: 1,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ListProjectNews("test-project", nil)
	if err != nil {
		t.Fatalf("ListProjectNews failed: %v", err)
	}

	if len(result.News) != 1 {
		t.Errorf("Expected 1 news item, got %d", len(result.News))
	}
}
