package redmine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListQueries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/queries.json" {
			t.Errorf("Expected path /queries.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-Api-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		response := QueriesResponse{
			Queries: []Query{
				{ID: 1, Name: "Query 1", IsPublic: true},
				{ID: 2, Name: "Query 2", IsPublic: false},
			},
			TotalCount: 2,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ListQueries(context.Background())
	if err != nil {
		t.Fatalf("ListQueries failed: %v", err)
	}

	if len(result.Queries) != 2 {
		t.Errorf("Expected 2 queries, got %d", len(result.Queries))
	}
	if result.Queries[0].Name != "Query 1" {
		t.Errorf("Expected name 'Query 1', got %s", result.Queries[0].Name)
	}
}
