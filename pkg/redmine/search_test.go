package redmine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/search.json" {
			t.Errorf("Expected path /search.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-Api-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		query := r.URL.Query().Get("q")
		if query != "test query" {
			t.Errorf("Expected query 'test query', got %s", query)
		}

		response := SearchResponse{
			Results: []SearchResult{
				{ID: 1, Title: "Result 1", Type: "issue"},
				{ID: 2, Title: "Result 2", Type: "wiki_page"},
			},
			TotalCount: 2,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	opts := &SearchOptions{
		Query: []string{"test query"},
	}
	result, err := client.Search(context.Background(), opts)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(result.Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(result.Results))
	}
	if result.Results[0].Title != "Result 1" {
		t.Errorf("Expected title 'Result 1', got %s", result.Results[0].Title)
	}
}
