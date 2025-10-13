package redmine

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListMemberships(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test-project/memberships.json" {
			t.Errorf("Expected path /projects/test-project/memberships.json, got %s", r.URL.Path)
		}

		response := MembershipsResponse{
			Memberships: []Membership{
				{ID: 1},
				{ID: 2},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	result, err := client.ListMemberships("test-project")
	if err != nil {
		t.Fatalf("ListMemberships failed: %v", err)
	}

	if len(result.Memberships) != 2 {
		t.Errorf("Expected 2 memberships, got %d", len(result.Memberships))
	}
}

func TestCreateMembership(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var req MembershipRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Membership.UserID != 5 {
			t.Errorf("Expected user_id 5, got %d", req.Membership.UserID)
		}

		w.WriteHeader(http.StatusCreated)
		response := MembershipResponse{
			Membership: Membership{ID: 999},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	membership := MembershipCreateUpdate{
		UserID:  5,
		RoleIDs: []int{1, 2},
	}
	result, err := client.CreateMembership("test-project", membership)
	if err != nil {
		t.Fatalf("CreateMembership failed: %v", err)
	}

	if result.Membership.ID != 999 {
		t.Errorf("Expected ID 999, got %d", result.Membership.ID)
	}
}

func TestUpdateMembership(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	err := client.UpdateMembership(123, []int{3, 4})
	if err != nil {
		t.Fatalf("UpdateMembership failed: %v", err)
	}
}
