package redmine

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListTimeEntries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/time_entries.json" {
			t.Errorf("Expected path /time_entries.json, got %s", r.URL.Path)
		}

		response := TimeEntriesResponse{
			TimeEntries: []TimeEntry{
				{ID: 1, Hours: 5.0, Comments: "Development work"},
				{ID: 2, Hours: 3.5, Comments: "Testing"},
			},
			TotalCount: 2,
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	result, err := client.ListTimeEntries(nil)
	if err != nil {
		t.Fatalf("ListTimeEntries failed: %v", err)
	}

	if len(result.TimeEntries) != 2 {
		t.Errorf("Expected 2 time entries, got %d", len(result.TimeEntries))
	}
	if result.TimeEntries[0].Hours != 5.0 {
		t.Errorf("Expected hours 5.0, got %f", result.TimeEntries[0].Hours)
	}
}

func TestListTimeEntriesWithFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("user_id") != "5" {
			t.Errorf("Expected user_id=5, got %s", query.Get("user_id"))
		}
		if query.Get("project_id") != "test-project" {
			t.Errorf("Expected project_id=test-project, got %s", query.Get("project_id"))
		}
		if query.Get("from") != "2024-01-01" {
			t.Errorf("Expected from=2024-01-01, got %s", query.Get("from"))
		}
		if query.Get("to") != "2024-12-31" {
			t.Errorf("Expected to=2024-12-31, got %s", query.Get("to"))
		}

		response := TimeEntriesResponse{TimeEntries: []TimeEntry{}}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	opts := &ListTimeEntriesOptions{
		UserID:    5,
		ProjectID: "test-project",
		From:      "2024-01-01",
		To:        "2024-12-31",
	}
	_, err := client.ListTimeEntries(opts)
	if err != nil {
		t.Fatalf("ListTimeEntries with filters failed: %v", err)
	}
}

func TestShowTimeEntry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/time_entries/123.json" {
			t.Errorf("Expected path /time_entries/123.json, got %s", r.URL.Path)
		}

		response := TimeEntryResponse{
			TimeEntry: TimeEntry{
				ID:       123,
				Hours:    8.0,
				Comments: "Full day work",
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	result, err := client.ShowTimeEntry(123)
	if err != nil {
		t.Fatalf("ShowTimeEntry failed: %v", err)
	}

	if result.TimeEntry.ID != 123 {
		t.Errorf("Expected ID 123, got %d", result.TimeEntry.ID)
	}
	if result.TimeEntry.Hours != 8.0 {
		t.Errorf("Expected hours 8.0, got %f", result.TimeEntry.Hours)
	}
}

func TestCreateTimeEntry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var req TimeEntryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.TimeEntry.Hours != 4.5 {
			t.Errorf("Expected hours 4.5, got %f", req.TimeEntry.Hours)
		}

		w.WriteHeader(http.StatusCreated)
		response := TimeEntryResponse{
			TimeEntry: TimeEntry{
				ID:    999,
				Hours: req.TimeEntry.Hours,
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	timeEntry := TimeEntry{
		Hours:    4.5,
		Comments: "Bug fixing",
	}
	result, err := client.CreateTimeEntry(timeEntry)
	if err != nil {
		t.Fatalf("CreateTimeEntry failed: %v", err)
	}

	if result.TimeEntry.ID != 999 {
		t.Errorf("Expected ID 999, got %d", result.TimeEntry.ID)
	}
}

func TestUpdateTimeEntry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/time_entries/123.json" {
			t.Errorf("Expected path /time_entries/123.json, got %s", r.URL.Path)
		}

		var req TimeEntryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	timeEntry := TimeEntry{
		Hours:    6.0,
		Comments: "Updated hours",
	}
	err := client.UpdateTimeEntry(123, timeEntry)
	if err != nil {
		t.Fatalf("UpdateTimeEntry failed: %v", err)
	}
}

func TestDeleteTimeEntry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/time_entries/123.json" {
			t.Errorf("Expected path /time_entries/123.json, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	err := client.DeleteTimeEntry(123)
	if err != nil {
		t.Fatalf("DeleteTimeEntry failed: %v", err)
	}
}
