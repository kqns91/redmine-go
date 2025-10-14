package redmine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListGroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/groups.json" {
			t.Errorf("Expected path /groups.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-Api-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		response := GroupsResponse{
			Groups: []Group{
				{ID: 1, Name: "Developers"},
				{ID: 2, Name: "Managers"},
			},
			TotalCount: 2,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ListGroups(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListGroups failed: %v", err)
	}

	if len(result.Groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(result.Groups))
	}
	if result.Groups[0].Name != "Developers" {
		t.Errorf("Expected name 'Developers', got %s", result.Groups[0].Name)
	}
}

func TestShowGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/groups/1.json" {
			t.Errorf("Expected path /groups/1.json, got %s", r.URL.Path)
		}

		response := GroupResponse{
			Group: Group{
				ID:   1,
				Name: "Developers",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ShowGroup(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("ShowGroup failed: %v", err)
	}

	if result.Group.Name != "Developers" {
		t.Errorf("Expected name 'Developers', got %s", result.Group.Name)
	}
}

func TestCreateGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/groups.json" {
			t.Errorf("Expected path /groups.json, got %s", r.URL.Path)
		}

		var req GroupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.Group.Name != "New Group" {
			t.Errorf("Expected name 'New Group', got %s", req.Group.Name)
		}

		w.WriteHeader(http.StatusCreated)
		response := GroupResponse{
			Group: Group{
				ID:   3,
				Name: req.Group.Name,
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	group := Group{
		Name: "New Group",
	}
	result, err := client.CreateGroup(context.Background(), group)
	if err != nil {
		t.Fatalf("CreateGroup failed: %v", err)
	}

	if result.Group.Name != "New Group" {
		t.Errorf("Expected name 'New Group', got %s", result.Group.Name)
	}
}

func TestDeleteGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/groups/1.json" {
			t.Errorf("Expected path /groups/1.json, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	err := client.DeleteGroup(context.Background(), 1)
	if err != nil {
		t.Fatalf("DeleteGroup failed: %v", err)
	}
}
