package redmine

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListRoles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/roles.json" {
			t.Errorf("Expected path /roles.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-API-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		response := RolesResponse{
			Roles: []Role{
				{ID: 1, Name: "Manager"},
				{ID: 2, Name: "Developer"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ListRoles()
	if err != nil {
		t.Fatalf("ListRoles failed: %v", err)
	}

	if len(result.Roles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(result.Roles))
	}
	if result.Roles[0].Name != "Manager" {
		t.Errorf("Expected name 'Manager', got %s", result.Roles[0].Name)
	}
}

func TestShowRole(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/roles/1.json" {
			t.Errorf("Expected path /roles/1.json, got %s", r.URL.Path)
		}

		response := RoleResponse{
			Role: Role{
				ID:   1,
				Name: "Manager",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ShowRole(1)
	if err != nil {
		t.Fatalf("ShowRole failed: %v", err)
	}

	if result.Role.Name != "Manager" {
		t.Errorf("Expected name 'Manager', got %s", result.Role.Name)
	}
}
