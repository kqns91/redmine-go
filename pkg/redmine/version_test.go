package redmine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListVersions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test-project/versions.json" {
			t.Errorf("Expected path /projects/test-project/versions.json, got %s", r.URL.Path)
		}

		response := VersionsResponse{
			Versions: []Version{
				{ID: 1, Name: "v1.0", Status: "open"},
				{ID: 2, Name: "v2.0", Status: "locked"},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	result, err := client.ListVersions(context.Background(), "test-project")
	if err != nil {
		t.Fatalf("ListVersions failed: %v", err)
	}

	if len(result.Versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(result.Versions))
	}
}

func TestCreateVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var req VersionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Version.Name != "v3.0" {
			t.Errorf("Expected name v3.0, got %s", req.Version.Name)
		}

		w.WriteHeader(http.StatusCreated)
		response := VersionResponse{
			Version: Version{
				ID:   3,
				Name: req.Version.Name,
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	version := Version{Name: "v3.0"}
	result, err := client.CreateVersion(context.Background(), "test-project", version)
	if err != nil {
		t.Fatalf("CreateVersion failed: %v", err)
	}

	if result.Version.ID != 3 {
		t.Errorf("Expected ID 3, got %d", result.Version.ID)
	}
}

func TestDeleteVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	err := client.DeleteVersion(context.Background(), 123)
	if err != nil {
		t.Fatalf("DeleteVersion failed: %v", err)
	}
}
