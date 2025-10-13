package redmine

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListTrackers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/trackers.json" {
			t.Errorf("Expected path /trackers.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-API-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		response := TrackersResponse{
			Trackers: []Tracker{
				{ID: 1, Name: "Bug"},
				{ID: 2, Name: "Feature"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ListTrackers()
	if err != nil {
		t.Fatalf("ListTrackers failed: %v", err)
	}

	if len(result.Trackers) != 2 {
		t.Errorf("Expected 2 trackers, got %d", len(result.Trackers))
	}
	if result.Trackers[0].Name != "Bug" {
		t.Errorf("Expected name 'Bug', got %s", result.Trackers[0].Name)
	}
}
