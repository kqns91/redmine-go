package redmine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := UsersResponse{
			Users: []User{
				{ID: 1, Login: "user1", Firstname: "John", Lastname: "Doe"},
				{ID: 2, Login: "user2", Firstname: "Jane", Lastname: "Smith"},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	result, err := client.ListUsers(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}

	if len(result.Users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(result.Users))
	}
}

func TestGetCurrentUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/current.json" {
			t.Errorf("Expected path /users/current.json, got %s", r.URL.Path)
		}

		response := UserResponse{
			User: User{
				ID:    1,
				Login: "current_user",
				Admin: true,
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	result, err := client.GetCurrentUser(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetCurrentUser failed: %v", err)
	}

	if result.User.Login != "current_user" {
		t.Errorf("Expected login 'current_user', got '%s'", result.User.Login)
	}
	if !result.User.Admin {
		t.Error("Expected admin to be true")
	}
}
