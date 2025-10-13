package redmine

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMyAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/my/account.json" {
			t.Errorf("Expected path /my/account.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-API-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		response := UserResponse{
			User: User{
				ID:        1,
				Login:     "testuser",
				Firstname: "Test",
				Lastname:  "User",
				Mail:      "test@example.com",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.GetMyAccount()
	if err != nil {
		t.Fatalf("GetMyAccount failed: %v", err)
	}

	if result.User.Login != "testuser" {
		t.Errorf("Expected login 'testuser', got %s", result.User.Login)
	}
	if result.User.Mail != "test@example.com" {
		t.Errorf("Expected mail 'test@example.com', got %s", result.User.Mail)
	}
}
